// FILE: model.go
// This file contains the core data structures and simulation logic for the model.

package main

import (
	"bufio"
	"fmt"
	"log"
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
	Params                   Parameters
	TapetalPigment           float64
	ShieldingPigment         float64
	CircumferenceOfEye       float64
	ApertureRadius           float64
	EyeRadius                float64
	DistanceToAperture       float64
	AngleAtCenter            float64
	ApertureArc              float64
	OmmatidialAngle          float64
	NumberOfFacets           int
	RhabdomRadius            float64
	IncidenceOmmatidialAngle float64
	CriticalAngle            float64
	RefractedOmmatidialAngle float64
	RhabdomAlphaAngle        float64
	OldRhabdomLength         float64 // Store original length
	DebugMode                bool
}

var degToRadConv = math.Pi / 180.0
var radToDegConv = 180.0 / math.Pi

// NewModel initialises a model with a given set of parameters.
func NewModel(params Parameters) *Model {
	m := &Model{Params: params}
	m.initialCalculations()
	return m
}

// initialCalculations sets up the initial state of the model based on parameters.
func (m *Model) initialCalculations() {
	p := m.Params
	m.OldRhabdomLength = p.RhabdomLength // Store for reset
	m.CircumferenceOfEye = math.Pi * p.EyeDiameter
	m.ApertureRadius = p.ApertureDiameter / 2.0
	m.EyeRadius = p.EyeDiameter / 2.0
	m.DistanceToAperture = math.Sqrt(math.Pow(m.EyeRadius, 2) - math.Pow(m.ApertureRadius, 2))
	m.AngleAtCenter = math.Atan(m.ApertureRadius/m.DistanceToAperture) / degToRadConv
	m.ApertureArc = (m.AngleAtCenter / 360.0) * m.CircumferenceOfEye
	m.OmmatidialAngle = (p.FacetWidth / m.CircumferenceOfEye) * 360
	m.RhabdomRadius = p.RhabdomWidth / 2.0

	// In the python version, this is calculated, but doesn't seem to be used.
	// rhabdomSlope := math.Sqrt(math.Pow(p.RhabdomLength, 2) + math.Pow(m.RhabdomRadius, 2))
	// m.RhabdomAlphaAngle = math.Acos(p.RhabdomLength/rhabdomSlope) / degToRadConv

	// Number of facets in eyeshine patch axis
	m.NumberOfFacets = int(math.Round(m.ApertureArc / p.FacetWidth))

	// Angle of total internal reflection
	snellsLaw := math.Asin(p.CytoplasmRefractiveIndex / p.RhabdomRefractiveIndex)
	m.CriticalAngle = 90.0 - (snellsLaw / degToRadConv)

	// Reset angles
	m.IncidenceOmmatidialAngle = 0.0
	m.RefractedOmmatidialAngle = 0.0
}

// runModel executes the main simulation loop.
func (m *Model) runModel() {
	p := m.Params
	// Open files for output
	pathlengthsFile, err := os.Create(fmt.Sprintf("%s_pathlengths.csv", p.SpeciesName))
	if err != nil {
		log.Fatalf("Error creating pathlengths file: %v", err)
	}
	defer pathlengthsFile.Close()
	pathlengthsWriter := bufio.NewWriter(pathlengthsFile)
	defer pathlengthsWriter.Flush()

	var debugWriter *bufio.Writer
	if m.DebugMode {
		debugFile, err := os.Create(fmt.Sprintf("%s_debug.csv", p.SpeciesName))
		if err != nil {
			log.Fatalf("Error creating debug file: %v", err)
		}
		defer debugFile.Close()
		debugWriter = bufio.NewWriter(debugFile)
		defer debugWriter.Flush()
	}

	incrementAmount := p.RhabdomLength / 10.0

	// Main pigment loops (11 steps for Shielding pigment, 11 steps for Tapetal pigment)
	for pStep := 0; pStep <= 10; pStep++ {
		m.ShieldingPigment = float64(pStep) * incrementAmount
		for tStep := 0; tStep <= 10; tStep++ {
			m.TapetalPigment = float64(tStep) * incrementAmount

			fmt.Printf("P: %.2f, T: %.2f\n", m.ShieldingPigment, m.TapetalPigment)
			fmt.Fprintf(pathlengthsWriter, "%.6f\n%.6f\n", m.ShieldingPigment, m.TapetalPigment)
			if debugWriter != nil {
				fmt.Fprintf(debugWriter, "P: %.2f, T: %.2f\n", m.ShieldingPigment, m.TapetalPigment)
			}

			// Loop over each facet across the eyeshine patch axis
			for currentFacet := 0; currentFacet < m.NumberOfFacets; currentFacet++ {
				m.IncidenceOmmatidialAngle = float64(currentFacet) * m.OmmatidialAngle

				// Account for refraction at the cornea
				switch {
				case m.IncidenceOmmatidialAngle == 0:
					m.RefractedOmmatidialAngle = 0.0
				case m.IncidenceOmmatidialAngle > 0 && m.IncidenceOmmatidialAngle <= 15:
					m.RefractedOmmatidialAngle = (m.IncidenceOmmatidialAngle * 0.9494) + 0.004667
				case m.IncidenceOmmatidialAngle > 15 && m.IncidenceOmmatidialAngle <= 35:
					m.RefractedOmmatidialAngle = (m.IncidenceOmmatidialAngle * 0.9407) + 0.1648
				case m.IncidenceOmmatidialAngle > 35 && m.IncidenceOmmatidialAngle <= 50:
					m.RefractedOmmatidialAngle = (m.IncidenceOmmatidialAngle * 0.9196) + 0.8676
				case m.IncidenceOmmatidialAngle > 50 && m.IncidenceOmmatidialAngle <= 60:
					m.RefractedOmmatidialAngle = (m.IncidenceOmmatidialAngle * 0.8677) + 3.38
				case m.IncidenceOmmatidialAngle > 60:
					fmt.Println("UNREAL ANGLE AT CORNEA")
					m.RefractedOmmatidialAngle = 60.0
				}

				// Light loss at cone due to angle of incidence
				var facetNum float64
				if m.RefractedOmmatidialAngle == 0 {
					facetNum = 1.0
				} else {
					cc := p.FacetWidth / math.Abs(math.Tan(m.RefractedOmmatidialAngle*degToRadConv))
					var fw float64
					if cc > p.FacetWidth*2.0 {
						fw = math.Cos(m.IncidenceOmmatidialAngle*degToRadConv) * p.FacetWidth
					} else {
						ll := (2.0 * cc) - (2.0 * p.FacetWidth)
						fw = math.Sin(m.IncidenceOmmatidialAngle*degToRadConv) * ll
					}
					facetNum = fw / p.FacetWidth
				}
				if facetNum > 1.0 {
					facetNum = 1.0
				}
				if facetNum < 0.0 {
					facetNum = 0.0
				}

				// Leading zeros for blur circle off-axis shift and angle adjustment
				var rowData []string
				boa := m.RefractedOmmatidialAngle
				if p.BlurCircleExtent > 0 {
					fd := float64(m.NumberOfFacets) / p.BlurCircleExtent
					for i := 1; i <= int(p.BlurCircleExtent); i++ {
						if float64(currentFacet) > (fd * float64(i)) {
							boa += m.OmmatidialAngle
							rowData = append(rowData, "0")
						}
					}
				}

				// Ray tracing through rhabdom array
				rhabdomLength := m.OldRhabdomLength
				cz := 0
				for {
					// Account for tapered/pointy rhabdom shape at proximal entrance
					if boa > m.CriticalAngle && cz == 0 {
						boa -= p.ProximalRhabdomAngle
					}

					if m.IncidenceOmmatidialAngle == 0 {
						// CASE 4: Perpendicular ray
						var val float64
						if m.TapetalPigment == 0 || m.ShieldingPigment > 0 {
							val = rhabdomLength * facetNum
						} else {
							val = (rhabdomLength * 2.0) * facetNum
						}
						rowData = append(rowData, fmt.Sprintf("%.6f", val))
						break
					}

					y := m.RhabdomRadius / math.Abs(math.Tan(boa*degToRadConv))

					if y >= rhabdomLength {
						// CASE 3: Bounce off base
						var x, v, val float64
						mx := math.Sqrt(math.Pow(rhabdomLength, 2) + math.Pow(m.RhabdomRadius, 2))
						if y == rhabdomLength {
							x = mx
						} else {
							x = rhabdomLength / math.Abs(math.Cos(boa*degToRadConv))
						}
						v = m.OldRhabdomLength / math.Abs(math.Cos(boa*degToRadConv))

						if m.TapetalPigment == 0 || m.ShieldingPigment > 0 {
							val = x * facetNum
						} else {
							val = (x + v) * facetNum
						}
						rowData = append(rowData, fmt.Sprintf("%.6f", val))
						break
					} else if y > (rhabdomLength-m.ShieldingPigment) || y > (rhabdomLength-m.TapetalPigment) || boa < m.CriticalAngle {
						// CASE 2: Reflection from edge
						var val float64
						x := m.RhabdomRadius / math.Abs(math.Sin(boa*degToRadConv))
						z := (rhabdomLength - y) / math.Abs(math.Cos(boa*degToRadConv))
						if z > x {
							z = x
						}
						v := m.OldRhabdomLength / math.Abs(math.Cos(boa*degToRadConv))
						if m.TapetalPigment == 0 {
							val = (x + z) * facetNum
						} else {
							val = (x + z + v) * facetNum
						}
						if m.ShieldingPigment > 0 {
							val = (x + z) * facetNum
						}
						if m.ShieldingPigment > (rhabdomLength - y) {
							val = x * facetNum
						}
						rowData = append(rowData, fmt.Sprintf("%.6f", val))
						break
					} else {
						// CASE 1: No reflection (ray passes through rhabdom wall into adjacent rhabdom)
						x := m.RhabdomRadius / math.Abs(math.Sin(boa*degToRadConv))
						rowData = append(rowData, fmt.Sprintf("%.6f", x*facetNum))
						rhabdomLength -= y
						boa += m.OmmatidialAngle
						cz = 1
						if rhabdomLength <= m.TapetalPigment || rhabdomLength <= m.ShieldingPigment {
							break
						}
					}
				}

				// Write row
				fmt.Fprintln(pathlengthsWriter, strings.Join(rowData, ","))
			}

			fmt.Fprintln(pathlengthsWriter, "999")
		}
	}
}
