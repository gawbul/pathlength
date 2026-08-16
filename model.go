// FILE: model.go
// This file contains the core data structures and simulation logic for the model.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

// Parameters holds all the eye-specific configuration.
type Parameters struct {
	SpeciesName              string
	RhabdomLength            float64
	RhabdomWidth             float64
	EyeDiameter              float64
	FacetWidth               float64
	ApertureDiameter         float64
	CytoplasmRefractiveIndex float64
	RhabdomRefractiveIndex   float64
	BlurCircleExtent         float64
	ProximalRhabdomAngle     float64
}

// Model holds the calculated parameters and state of the simulation.
type Model struct {
	Params             Parameters
	CircumferenceOfEye float64
	ApertureRadius     float64
	EyeRadius          float64
	DistanceToAperture float64
	AngleAtCenter      float64
	ApertureArc        float64
	OmmatidialAngle    float64
	NumberOfFacets     int
	RhabdomRadius      float64
	CriticalAngle      float64
	DebugMode          bool
}

const (
	degToRadConv = math.Pi / 180.0
	radToDegConv = 180.0 / math.Pi

	// pigmentSteps is the number of migration positions sampled for each pigment,
	// from fully retracted (0) to fully covering the rhabdom (RhabdomLength).
	pigmentSteps = 11

	// maxPropagationAngle is the largest angle to the rhabdom axis at which a ray
	// can still advance towards the proximal end. At or beyond 90 degrees the ray
	// travels perpendicular to (or back along) the axis and is treated as lost.
	maxPropagationAngle = 90.0
)

// NewModel initialises a model with a given set of parameters. It returns an
// error if the parameters do not describe a physically realisable eye, rather
// than allowing NaNs to propagate silently into the results.
func NewModel(params Parameters) (*Model, error) {
	if err := validateParameters(params); err != nil {
		return nil, err
	}
	m := &Model{Params: params}
	m.initialCalculations()
	if err := m.validateGeometry(); err != nil {
		return nil, err
	}
	return m, nil
}

// validateParameters rejects inputs that would produce NaNs or nonsensical optics.
func validateParameters(p Parameters) error {
	if strings.TrimSpace(p.SpeciesName) == "" {
		return fmt.Errorf("species name is required")
	}
	// Reject non-finite values before any range check. strconv.ParseFloat accepts
	// "NaN" and "Inf", and every ordered comparison against NaN is false, so a NaN
	// would slip past the constraints below and reappear as an undefined critical
	// angle or blur offset - reproducing the very silent failures the range checks
	// exist to prevent.
	for _, c := range []struct {
		name string
		v    float64
	}{
		{"rhabdom length", p.RhabdomLength},
		{"rhabdom width", p.RhabdomWidth},
		{"eye diameter", p.EyeDiameter},
		{"facet width", p.FacetWidth},
		{"aperture diameter", p.ApertureDiameter},
		{"cytoplasm refractive index", p.CytoplasmRefractiveIndex},
		{"rhabdom refractive index", p.RhabdomRefractiveIndex},
		{"blur circle extent", p.BlurCircleExtent},
		{"proximal rhabdom angle", p.ProximalRhabdomAngle},
	} {
		if math.IsNaN(c.v) || math.IsInf(c.v, 0) {
			return fmt.Errorf("%s must be a finite number, got %g", c.name, c.v)
		}
	}

	for _, c := range []struct {
		name string
		v    float64
	}{
		{"rhabdom length", p.RhabdomLength},
		{"rhabdom width", p.RhabdomWidth},
		{"eye diameter", p.EyeDiameter},
		{"facet width", p.FacetWidth},
		{"aperture diameter", p.ApertureDiameter},
	} {
		if c.v <= 0 {
			return fmt.Errorf("%s must be greater than 0 um, got %g", c.name, c.v)
		}
	}
	if p.ApertureDiameter >= p.EyeDiameter {
		return fmt.Errorf("aperture diameter (%g um) must be smaller than eye diameter (%g um)",
			p.ApertureDiameter, p.EyeDiameter)
	}
	if p.CytoplasmRefractiveIndex <= 1.0 {
		return fmt.Errorf("cytoplasm refractive index must be greater than 1.0, got %g",
			p.CytoplasmRefractiveIndex)
	}
	// Without n_rhabdom > n_cytoplasm there is no waveguiding and no critical angle:
	// asin(n_cyt/n_rhab) would be NaN and every total-internal-reflection test would
	// silently evaluate to false.
	if p.RhabdomRefractiveIndex <= p.CytoplasmRefractiveIndex {
		return fmt.Errorf("rhabdom refractive index (%g) must exceed cytoplasm refractive index (%g) for total internal reflection",
			p.RhabdomRefractiveIndex, p.CytoplasmRefractiveIndex)
	}
	if p.BlurCircleExtent < 1 {
		return fmt.Errorf("blur circle extent must be at least 1 rhabdom, got %g", p.BlurCircleExtent)
	}
	if p.ProximalRhabdomAngle < 0 {
		return fmt.Errorf("proximal rhabdom angle must not be negative, got %g", p.ProximalRhabdomAngle)
	}
	return nil
}

// validateGeometry checks the values derived from the parameters.
func (m *Model) validateGeometry() error {
	if m.NumberOfFacets < 1 {
		return fmt.Errorf("eyeshine patch spans no facets (aperture arc %.2f um / facet width %.2f um); "+
			"check aperture and facet dimensions", m.ApertureArc, m.Params.FacetWidth)
	}
	// The blur circle spreads the light gathered across the eyeshine patch over
	// BlurCircleExtent rhabdoms. With more blur steps than facets, the mapping from
	// facet to rhabdom offset leaves offsets with no contributing facet at all,
	// producing a comb-shaped profile whose half-maximum is a binning artefact.
	if m.Params.BlurCircleExtent > float64(m.NumberOfFacets) {
		return fmt.Errorf("blur circle extent (%g rhabdoms) exceeds the %d facets across the eyeshine patch; "+
			"the resulting profile would contain rhabdom offsets that receive no light",
			m.Params.BlurCircleExtent, m.NumberOfFacets)
	}
	return nil
}

// initialCalculations sets up the derived optical values of the model.
func (m *Model) initialCalculations() {
	p := m.Params
	m.CircumferenceOfEye = math.Pi * p.EyeDiameter
	m.ApertureRadius = p.ApertureDiameter / 2.0
	m.EyeRadius = p.EyeDiameter / 2.0
	m.DistanceToAperture = math.Sqrt(math.Pow(m.EyeRadius, 2) - math.Pow(m.ApertureRadius, 2))
	m.AngleAtCenter = math.Atan(m.ApertureRadius/m.DistanceToAperture) * radToDegConv
	m.ApertureArc = (m.AngleAtCenter / 360.0) * m.CircumferenceOfEye
	m.OmmatidialAngle = (p.FacetWidth / m.CircumferenceOfEye) * 360.0
	m.RhabdomRadius = p.RhabdomWidth / 2.0

	// Number of facets along the radius of the eyeshine patch.
	m.NumberOfFacets = int(math.Round(m.ApertureArc / p.FacetWidth))

	// Angle of total internal reflection. boa is measured from the rhabdom axis, so
	// the angle at the wall normal is (90 - boa) and light is guided when
	// (90 - boa) exceeds the Snell critical angle, i.e. when boa < CriticalAngle.
	snellsLaw := math.Asin(p.CytoplasmRefractiveIndex/p.RhabdomRefractiveIndex) * radToDegConv
	m.CriticalAngle = 90.0 - snellsLaw
}

// refractedAngle applies the empirical corneal refraction regression to an angle
// of incidence, in degrees.
func refractedAngle(incidence float64) float64 {
	switch {
	case incidence <= 0:
		return 0.0
	case incidence <= 15:
		return (incidence * 0.9494) + 0.004667
	case incidence <= 35:
		return (incidence * 0.9407) + 0.1648
	case incidence <= 50:
		return (incidence * 0.9196) + 0.8676
	case incidence <= 60:
		return (incidence * 0.8677) + 3.38
	default:
		return math.NaN()
	}
}

// facetTransmission is the fraction of light a facet admits at the given angle of
// incidence, relative to a facet viewed normally. It is a flux factor in [0, 1]
// and is applied to the absorbed intensity, not to the geometric path length.
func (m *Model) facetTransmission(facetIndex int) float64 {
	incidence := float64(facetIndex) * m.OmmatidialAngle
	refracted := refractedAngle(incidence)
	if math.IsNaN(refracted) {
		// Beyond the range the corneal regression covers, no light gets through.
		return 0.0
	}
	if refracted == 0 {
		return 1.0
	}
	cc := m.Params.FacetWidth / math.Abs(math.Tan(refracted*degToRadConv))
	var fw float64
	if cc > m.Params.FacetWidth*2.0 {
		fw = math.Cos(incidence*degToRadConv) * m.Params.FacetWidth
	} else {
		ll := (2.0 * cc) - (2.0 * m.Params.FacetWidth)
		fw = math.Sin(incidence*degToRadConv) * ll
	}
	return math.Min(1.0, math.Max(0.0, fw/m.Params.FacetWidth))
}

// blurOffset is the displacement, in rhabdoms, of the image formed by the facet at
// the given radial position in the eyeshine patch.
//
// The blur circle spans BlurCircleExtent rhabdoms, so the outermost facet is
// displaced by (BlurCircleExtent - 1) and the central facet by zero. The offset is
// continuous in the facet index: quantising it to whole rhabdoms (as earlier
// versions did, via a chain of `facet > fd*i` tests) aliased the facets unevenly
// across the available offsets and cut notches into the profile at offsets where
// fd*i landed on an exact integer.
func (m *Model) blurOffset(facetIndex int) float64 {
	if m.NumberOfFacets <= 1 {
		return 0.0
	}
	return float64(facetIndex) * (m.Params.BlurCircleExtent - 1.0) / float64(m.NumberOfFacets-1)
}

// traceResult is the outcome of tracing one ray through the rhabdom array.
type traceResult struct {
	// Pathlengths through each successive rhabdom the ray enters, in micrometres of
	// raw geometry. Facet transmission is NOT folded in here; it is a flux factor
	// applied when the absorbed intensity is computed.
	Pathlengths []float64
	// TerminalCase records which branch ended the trace, for the debug output.
	TerminalCase string
	// MaxAngle is the largest angle to the rhabdom axis reached during the trace.
	// Every segment covers some axial depth d at an angle no greater than this, so
	// the total path is bounded by 2*RhabdomLength/cos(MaxAngle).
	MaxAngle float64
	// Lost is set when the ray stopped propagating towards the proximal end.
	Lost bool
}

// traceRay follows a single ray from the given facet through the rhabdom array for
// one combination of pigment positions.
func (m *Model) traceRay(facetIndex int, shielding, tapetal float64) traceResult {
	p := m.Params
	res := traceResult{}

	// Angle to the rhabdom axis on entry: corneal refraction plus the blur-circle
	// displacement, which tilts the ray by one ommatidial angle per rhabdom offset.
	boa := refractedAngle(float64(facetIndex)*m.OmmatidialAngle) + m.blurOffset(facetIndex)*m.OmmatidialAngle

	rhabdomLength := p.RhabdomLength
	cz := 0

	for {
		// Tapered ("pointy") rhabdom tip widens the acceptance angle on first entry.
		if boa > m.CriticalAngle && cz == 0 {
			boa -= p.ProximalRhabdomAngle
			if boa < 0 {
				boa = 0
			}
		}

		// A ray at 90 degrees or more to the axis cannot advance towards the
		// proximal end. Earlier versions took the absolute value of tan and cos,
		// which silently folded such rays back and produced path lengths many
		// times the rhabdom length.
		if boa >= maxPropagationAngle || math.IsNaN(boa) {
			res.TerminalCase = "lost"
			res.Lost = true
			return res
		}
		if boa > res.MaxAngle {
			res.MaxAngle = boa
		}

		if facetIndex == 0 {
			// CASE 4: axial ray. Equivalent to case 3 at boa = 0, kept explicit.
			val := rhabdomLength
			if tapetal > 0 && shielding == 0 {
				val = rhabdomLength * 2.0
			}
			res.Pathlengths = append(res.Pathlengths, val)
			res.TerminalCase = "C4"
			return res
		}

		sin := math.Sin(boa * degToRadConv)
		cos := math.Cos(boa * degToRadConv)
		tan := math.Tan(boa * degToRadConv)

		// Axial distance travelled before the ray meets the rhabdom wall.
		y := m.RhabdomRadius / tan

		switch {
		case y >= rhabdomLength:
			// CASE 3: the ray reaches the base without meeting the wall and is
			// reflected by the tapetum at the base.
			x := rhabdomLength / cos
			v := p.RhabdomLength / cos
			val := x
			if tapetal > 0 && shielding == 0 {
				val = x + v
			}
			res.Pathlengths = append(res.Pathlengths, val)
			res.TerminalCase = "C3"
			return res

		case y > (rhabdomLength-shielding) || y > (rhabdomLength-tapetal) || boa < m.CriticalAngle:
			// CASE 2: the ray is reflected at the wall, either by total internal
			// reflection or by the tapetal mirror.
			x := m.RhabdomRadius / sin

			// A guided ray never leaves the rhabdom, so the proximal screening
			// pigment - which lies in the cytoplasm outside the rhabdom - cannot
			// absorb it. An unguided ray exits through the wall and is absorbed
			// where the pigment starts.
			guided := boa < m.CriticalAngle
			axial := rhabdomLength - y
			if shielding > 0 && !guided {
				axial = math.Max(0, rhabdomLength-shielding-y)
			}
			z := axial / cos
			v := p.RhabdomLength / cos

			val := x + z
			if tapetal > 0 && shielding == 0 {
				val = x + z + v
			}
			res.Pathlengths = append(res.Pathlengths, val)
			res.TerminalCase = "C2"
			return res

		default:
			// CASE 1: no reflection. The ray crosses the wall into the adjacent
			// rhabdom, and the inter-rhabdom angle steps by one ommatidial angle.
			res.Pathlengths = append(res.Pathlengths, m.RhabdomRadius/sin)
			rhabdomLength -= y
			boa += m.OmmatidialAngle
			cz = 1
			if rhabdomLength <= tapetal || rhabdomLength <= shielding {
				res.TerminalCase = "C1"
				return res
			}
		}
	}
}

// pathlengthsHeader labels the columns of the raw geometry output. Every row carries
// its own keys, so the file is a plain rectangular CSV with no positional state and
// no block terminator.
const pathlengthsHeader = "block,shielding_um,tapetal_um,facet,rhabdom,pathlength_um"

// runModel executes the main simulation loop, writes the raw pathlength geometry, and
// returns the resolution and sensitivity of each pigment state.
//
// The absorption profile is accumulated as the rays are traced rather than by reading
// the file back, so the summary does not depend on the output format at all.
func (m *Model) runModel() ([]blockSummary, error) {
	p := m.Params

	pathlengthsFile, err := os.Create(fmt.Sprintf("%s_pathlengths.csv", p.SpeciesName))
	if err != nil {
		return nil, fmt.Errorf("creating pathlengths file: %w", err)
	}
	defer pathlengthsFile.Close()
	pathlengthsWriter := bufio.NewWriter(pathlengthsFile)
	defer pathlengthsWriter.Flush()

	fmt.Fprintln(pathlengthsWriter, pathlengthsHeader)

	var debugWriter *bufio.Writer
	if m.DebugMode {
		debugFile, err := os.Create(fmt.Sprintf("%s_debug.csv", p.SpeciesName))
		if err != nil {
			return nil, fmt.Errorf("creating debug file: %w", err)
		}
		defer debugFile.Close()
		debugWriter = bufio.NewWriter(debugFile)
		defer debugWriter.Flush()
		fmt.Fprintln(debugWriter,
			"block,shielding_um,tapetal_um,facet,incidence_deg,refracted_deg,blur_offset_rhabdoms,entry_boa_deg,facet_transmission,terminal_case,rhabdoms_entered,pathlengths_um")
	}

	incrementAmount := p.RhabdomLength / 10.0
	lostRays := 0
	block := 0
	summaries := make([]blockSummary, 0, pigmentSteps*pigmentSteps)

	for pStep := 0; pStep < pigmentSteps; pStep++ {
		shielding := float64(pStep) * incrementAmount
		for tStep := 0; tStep < pigmentSteps; tStep++ {
			tapetal := float64(tStep) * incrementAmount

			// Area-weighted absorbed light at each rhabdom offset from the optic axis.
			var profile []float64

			for facet := 0; facet < m.NumberOfFacets; facet++ {
				trace := m.traceRay(facet, shielding, tapetal)
				if trace.Lost {
					lostRays++
				}
				profile = m.accumulate(profile, facet, trace.Pathlengths)

				if len(trace.Pathlengths) == 0 {
					// A lost ray absorbs nothing, but the facet still belongs in the
					// record, so emit an explicit zero for it.
					fmt.Fprintf(pathlengthsWriter, "%d,%.6f,%.6f,%d,0,0.000000\n",
						block, shielding, tapetal, facet)
				}
				for rhabdom, v := range trace.Pathlengths {
					fmt.Fprintf(pathlengthsWriter, "%d,%.6f,%.6f,%d,%d,%.6f\n",
						block, shielding, tapetal, facet, rhabdom, v)
				}

				if debugWriter != nil {
					parts := make([]string, len(trace.Pathlengths))
					for i, v := range trace.Pathlengths {
						parts[i] = fmt.Sprintf("%.6f", v)
					}
					incidence := float64(facet) * m.OmmatidialAngle
					fmt.Fprintf(debugWriter, "%d,%.4f,%.4f,%d,%.4f,%.4f,%.4f,%.4f,%.6f,%s,%d,%s\n",
						block, shielding, tapetal, facet, incidence, refractedAngle(incidence),
						m.blurOffset(facet),
						refractedAngle(incidence)+m.blurOffset(facet)*m.OmmatidialAngle,
						m.facetTransmission(facet), trace.TerminalCase,
						len(trace.Pathlengths), strings.Join(parts, " "))
				}
			}

			summaries = append(summaries, m.summariseBlock(profile))
			block++
		}
	}

	if lostRays > 0 {
		fmt.Printf("WARNING: %d of %d rays exceeded 90 degrees to the rhabdom axis and were discarded.\n",
			lostRays, pigmentSteps*pigmentSteps*m.NumberOfFacets)
	}
	return summaries, nil
}
