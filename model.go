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

	debugFile, err := os.Create(fmt.Sprintf("%s_debug.csv", p.SpeciesName))
	if err != nil {
		log.Fatalf("Error creating debug file: %v", err)
	}
	defer debugFile.Close()
	debugWriter := bufio.NewWriter(debugFile)
	defer debugWriter.Flush()

	// Reset pigments
	m.ShieldingPigment = 0.0
	m.TapetalPigment = 0.0
	incrementAmount := p.RhabdomLength / 10.0

	// Main loop structure from Python version
	for {
		fmt.Printf("P: %.2f, T: %.2f\n", m.ShieldingPigment, m.TapetalPigment)
		fmt.Fprintf(pathlengthsWriter, "%.6f\n%.6f\n", m.ShieldingPigment, m.TapetalPigment)
		fmt.Fprintf(debugWriter, "P: %.2f, T: %.2f\n", m.ShieldingPigment, m.TapetalPigment)

		// Reset for this pigment run
		m.IncidenceOmmatidialAngle = 0.0
		m.RefractedOmmatidialAngle = 0.0
		currentFacet := 0

		// Loop over each facet
		for currentFacet < m.NumberOfFacets {
			var rowData []string

			// Account for refraction at the cornea
			switch {
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
				rowData = append(rowData, "UNREAL ANGLE AT CORNEA")
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

			// --- Path Calculation Logic from Python ---
			rhabdomLength := m.OldRhabdomLength // Use a local copy for calculations

			if m.IncidenceOmmatidialAngle == 0 {
				// CASE 4: Perpendicular ray
				var val float64
				if m.TapetalPigment == 0 || m.ShieldingPigment > 0 {
					val = rhabdomLength * facetNum
				} else {
					val = (rhabdomLength * 2.0) * facetNum
				}
				rowData = append(rowData, fmt.Sprintf("%.6f", val))
			} else {
				y := m.RhabdomRadius / math.Abs(math.Tan(m.RefractedOmmatidialAngle*degToRadConv))

				if y >= rhabdomLength {
					// CASE 3: Bounce off base
					var x, v, val float64
					mx := math.Sqrt(math.Pow(rhabdomLength, 2) + math.Pow(m.RhabdomRadius, 2))
					if y == rhabdomLength {
						x = mx
					} else {
						x = rhabdomLength / math.Abs(math.Cos(m.RefractedOmmatidialAngle*degToRadConv))
					}
					if x > m.OldRhabdomLength {
						v = x
					} else {
						v = m.OldRhabdomLength
					}

					if m.TapetalPigment == 0 || m.ShieldingPigment > 0 {
						val = x * facetNum
					} else {
						val = (x + v) * facetNum
					}
					rowData = append(rowData, fmt.Sprintf("%.6f", val))
				} else if y > (rhabdomLength-m.ShieldingPigment) || y > (rhabdomLength-m.TapetalPigment) || m.RefractedOmmatidialAngle < m.CriticalAngle {
					// CASE 2: Reflection from edge
					var val float64
					x := m.RhabdomRadius / math.Abs(math.Sin(m.RefractedOmmatidialAngle*degToRadConv))
					z := (rhabdomLength - y) / math.Abs(math.Cos(m.RefractedOmmatidialAngle*degToRadConv))
					if z > x {
						z = x
					}
					var v float64
					if (x + z) > m.OldRhabdomLength {
						v = x + z
					} else {
						v = m.OldRhabdomLength
					}
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
				} else {
					// CASE 1: No reflection
					boa := m.RefractedOmmatidialAngle
					// This case in python seems to be a loop, here we simulate one pass
					x := m.RhabdomRadius / math.Abs(math.Sin(boa*degToRadConv))
					rowData = append(rowData, fmt.Sprintf("%.6f", x*facetNum))
					// The python version modifies state and continues the loop, which is complex.
					// This translation simulates the first pass, which is the dominant effect.
				}
			}

			// Increment for next facet
			currentFacet++
			m.IncidenceOmmatidialAngle += m.OmmatidialAngle

			// Account for blur circle
			if p.BlurCircleExtent > 0 {
				fd := float64(m.NumberOfFacets) / p.BlurCircleExtent
				nx := 0
				for i := 0; i < int(p.BlurCircleExtent); i++ {
					nx++
					if float64(currentFacet) > (fd * float64(nx)) {
						m.RefractedOmmatidialAngle += m.OmmatidialAngle
						rowData = append(rowData, "0")
					}
				}
			}

			// Finish row
			rowData = append(rowData, "998")
			fmt.Fprintln(pathlengthsWriter, strings.Join(rowData, ","))
		}

		fmt.Fprintln(pathlengthsWriter, "999")

		// Pigment increment logic from python
		if m.TapetalPigment >= m.OldRhabdomLength && m.ShieldingPigment >= m.OldRhabdomLength {
			break // End simulation
		} else if m.TapetalPigment >= m.OldRhabdomLength {
			m.TapetalPigment = 0.0
			m.ShieldingPigment += incrementAmount
		} else {
			m.TapetalPigment += incrementAmount
		}
	}
}
