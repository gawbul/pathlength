// FILE: csv.go
// This file contains functions for reading and writing CSV data.

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// absorptionCoefficient is the rhabdom absorption coefficient in um^-1, used in the
// Beer-Lambert absorbance 1 - exp(-k*L). Reported values for crustacean rhabdoms
// span roughly 0.0067 to 0.01 um^-1.
const absorptionCoefficient = 0.01

// ringArea is the area, in squared facet widths, of the annulus of rhabdoms lying at
// the given whole-rhabdom offset from the optic axis. It is the same measure used to
// weight the contributing facets, so dividing an area-weighted total by it yields a
// radial area density.
func ringArea(offset int) float64 {
	if offset == 0 {
		return math.Pi * 0.25
	}
	outer := math.Pi * math.Pow(float64(offset)+0.5, 2)
	inner := math.Pi * math.Pow(float64(offset)-0.5, 2)
	return outer - inner
}

// deposit adds an amount at the given offset, growing the slice as required. Earlier
// versions used a fixed 21-element array and silently discarded everything beyond it,
// which lost up to a quarter of the absorbed light for widely blurred eyes.
func deposit(dst []float64, offset int, amount float64) []float64 {
	if amount == 0 {
		return dst
	}
	for len(dst) <= offset {
		dst = append(dst, 0)
	}
	dst[offset] += amount
	return dst
}

// blockSummary holds the resolution and sensitivity derived from one pigment block.
type blockSummary struct {
	// FWHMDegrees is the acceptance angle: the full width at half maximum of the
	// angular sensitivity function, in degrees. NaN when the profile carries no light
	// or is annular, in which case there is no acceptance angle to report.
	FWHMDegrees float64
	// Percentage of incident light absorbed, averaged over the eyeshine patch (0-100).
	SensitivityPercent float64
	// PeakOffset is the rhabdom offset carrying the most light.
	PeakOffset int
	// Annular is set when the profile dips below half its maximum on the optic axis,
	// so the light forms a ring rather than a central spot.
	Annular bool
}

// summariseBlock converts one block's area-weighted absorption profile into
// resolution and sensitivity.
func (m *Model) summariseBlock(rhabdoms []float64) blockSummary {
	out := blockSummary{FWHMDegrees: math.NaN()}

	// Sensitivity: the area-weighted mean of the absorbed percentage over the
	// eyeshine patch. The facet weights telescope to exactly pi*(N-0.5)^2, so
	// dividing by that area makes this a true weighted mean in the range 0-100.
	total := 0.0
	for _, v := range rhabdoms {
		total += v
	}
	patchArea := math.Pi * math.Pow(float64(m.NumberOfFacets)-0.5, 2)
	if patchArea > 0 {
		out.SensitivityPercent = total / patchArea
	}

	if len(rhabdoms) == 0 {
		return out
	}

	// Point spread function: light per unit area at each rhabdom offset. The
	// contributing facets are weighted by their source annulus, so the light arriving
	// at an offset must be divided by the annulus it is spread over to recover an
	// intensity. Without this the profile rises monotonically with offset simply
	// because outer annuli contain more ommatidia.
	// The profile ends at the outermost rhabdom that receives any light, so the next
	// offset out is genuinely dark. Including that zero captures the falling edge of
	// the blur circle, which is where a top-hat profile crosses its half maximum.
	psf := make([]float64, len(rhabdoms)+1)
	for j := range rhabdoms {
		psf[j] = rhabdoms[j] / ringArea(j)
	}

	peak := 0
	for j, v := range psf {
		if v > psf[peak] {
			peak = j
		}
	}
	out.PeakOffset = peak
	if psf[peak] <= 0 {
		return out
	}
	half := psf[peak] / 2.0

	// A profile that is already below half its maximum on the optic axis is annular:
	// the light forms a ring, and the region above half maximum is a band that does
	// not contain the axis. There is no acceptance angle to report, so the width is
	// left undefined rather than substituting the ring's thickness for it.
	if psf[0] < half {
		out.Annular = true
		return out
	}

	// The angular sensitivity function is even about the optic axis - offset j stands
	// for both +j and -j - so its full width at half maximum is twice the radius at
	// which it first falls below half. That radius is measured from the axis, not from
	// the peak: on a flat-topped profile whose maximum sits slightly off-axis,
	// measuring from the peak would understate the width by the peak's own offset.
	for i := 0; i < len(psf)-1; i++ {
		if psf[i] >= half && psf[i+1] < half {
			frac := (psf[i] - half) / (psf[i] - psf[i+1])
			out.FWHMDegrees = 2.0 * (float64(i) + frac) * m.OmmatidialAngle
			break
		}
	}
	return out
}

// accumulate adds one facet's traced ray into the area-weighted absorption profile,
// which records how much light reaches each whole-rhabdom offset from the optic axis.
func (m *Model) accumulate(profile []float64, facetIndex int, pathlengths []float64) []float64 {
	// Light gathered by this facet, and the rhabdom offset its image lands on.
	transmission := m.facetTransmission(facetIndex)
	sourceArea := ringArea(facetIndex)
	offset := m.blurOffset(facetIndex)
	base := int(math.Floor(offset))
	frac := offset - float64(base)

	tot := 0.0
	for rhabdom, pathlength := range pathlengths {
		if pathlength <= 0 {
			continue
		}
		// Fraction of the light still travelling that this rhabdom absorbs.
		absorbed := (1.0 - tot) * (1.0 - math.Exp(-absorptionCoefficient*pathlength))
		tot += absorbed
		// Facet transmission attenuates the flux entering the eye; it does not shorten
		// the geometric path, so it multiplies the absorbed intensity rather than the
		// exponent.
		weighted := 100.0 * transmission * absorbed * sourceArea

		// The blur displacement is continuous, so split the light between the two
		// rhabdom offsets that bracket it.
		profile = deposit(profile, base+rhabdom, weighted*(1.0-frac))
		if frac > 0 {
			profile = deposit(profile, base+rhabdom+1, weighted*frac)
		}
	}
	return profile
}

// calculateRessens writes the resolution and sensitivity matrices for the pigment
// states accumulated during the simulation.
func (m *Model) calculateRessens(summaries []blockSummary) error {
	p := m.Params
	fmt.Println("INFO: Calculating resolution and sensitivity...")

	if len(summaries) != pigmentSteps*pigmentSteps {
		return fmt.Errorf("expected %d pigment states, got %d", pigmentSteps*pigmentSteps, len(summaries))
	}

	if err := writeSummaryMatrix(fmt.Sprintf("%s_summary_res.csv", p.SpeciesName), summaries,
		func(b blockSummary) float64 { return b.FWHMDegrees }); err != nil {
		return err
	}
	if err := writeSummaryMatrix(fmt.Sprintf("%s_summary_sen.csv", p.SpeciesName), summaries,
		func(b blockSummary) float64 { return b.SensitivityPercent }); err != nil {
		return err
	}

	dark, annular := 0, 0
	for _, s := range summaries {
		switch {
		case s.Annular:
			annular++
		case math.IsNaN(s.FWHMDegrees):
			dark++
		}
	}
	if annular > 0 {
		fmt.Printf("WARNING: %d of %d pigment states have an annular profile, with the light "+
			"forming a ring rather than a central spot; they have no acceptance angle and are "+
			"reported as NaN.\n", annular, len(summaries))
	}
	if dark > 0 {
		fmt.Printf("WARNING: %d of %d pigment states absorb no light; their resolution is "+
			"reported as NaN.\n", dark, len(summaries))
	}
	return nil
}

// writeSummaryMatrix writes an 11x11 matrix with shielding pigment position varying
// down the rows and tapetal pigment position across the columns.
func writeSummaryMatrix(filename string, summaries []blockSummary, value func(blockSummary) float64) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating %s: %w", filename, err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for row := 0; row < pigmentSteps; row++ {
		cells := make([]string, pigmentSteps)
		for col := 0; col < pigmentSteps; col++ {
			cells[col] = strconv.FormatFloat(value(summaries[row*pigmentSteps+col]), 'f', 4, 64)
		}
		if _, err := fmt.Fprintln(writer, strings.Join(cells, ",")); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}
	return writer.Flush()
}

// parseInputParameters reads a parameter file using Go's standard CSV reader.
// It returns a slice of parsed parameters and an error if parsing fails.
func parseInputParameters(filename string) ([]Parameters, error) {
	var paramsList []Parameters // A slice to hold multiple parameter sets

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open parameter file %s: %w", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		// Check for a specific parsing error (like wrong number of fields)
		if parseErr, ok := err.(*csv.ParseError); ok && parseErr.Err == csv.ErrFieldCount {
			log.Printf("Skipping malformed record on line %d: %v", parseErr.Line, err)
			continue // Skip to the next record
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		if len(record) != 10 {
			log.Printf("Skipping malformed record (expected 10 fields, got %d): %v", len(record), record)
			continue
		}

		var params Parameters
		// (sn, rl, rw, ed, fw, ad, cri, rri, bce, pra)
		numbers := make([]float64, 9)
		bad := false
		for i := 1; i < 10; i++ {
			v, err := strconv.ParseFloat(strings.TrimSpace(record[i]), 64)
			if err != nil {
				log.Printf("Skipping record %q: field %d (%q) is not a number", record[0], i+1, record[i])
				bad = true
				break
			}
			numbers[i-1] = v
		}
		if bad {
			continue
		}

		params.SpeciesName = strings.TrimSpace(record[0])
		params.RhabdomLength = numbers[0]
		params.RhabdomWidth = numbers[1]
		params.EyeDiameter = numbers[2]
		params.FacetWidth = numbers[3]
		params.ApertureDiameter = numbers[4]
		params.CytoplasmRefractiveIndex = numbers[5]
		params.RhabdomRefractiveIndex = numbers[6]
		params.BlurCircleExtent = numbers[7]
		params.ProximalRhabdomAngle = numbers[8]

		paramsList = append(paramsList, params)
	}

	if len(paramsList) == 0 {
		return nil, fmt.Errorf("no valid parameter data found in %s", filename)
	}

	return paramsList, nil
}
