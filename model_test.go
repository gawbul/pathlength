// FILE: model_test.go
// This file contains tests for the model initialization and simulation logic.

package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestInitialCalculations(t *testing.T) {
	params := Parameters{
		SpeciesName:              "test_eye",
		RhabdomLength:            180,
		RhabdomWidth:             25,
		EyeDiameter:              7800,
		FacetWidth:               50,
		ApertureDiameter:         3200,
		CytoplasmRefractiveIndex: 1.34,
		RhabdomRefractiveIndex:   1.37,
		BlurCircleExtent:         18,
		ProximalRhabdomAngle:     12.5,
	}

	model := NewModel(params)

	expectedCircumference := math.Pi * 7800.0
	if math.Abs(model.CircumferenceOfEye-expectedCircumference) > 1e-6 {
		t.Errorf("Expected CircumferenceOfEye %f, got %f", expectedCircumference, model.CircumferenceOfEye)
	}

	expectedOmmatidialAngle := (50.0 / expectedCircumference) * 360.0
	if math.Abs(model.OmmatidialAngle-expectedOmmatidialAngle) > 1e-6 {
		t.Errorf("Expected OmmatidialAngle %f, got %f", expectedOmmatidialAngle, model.OmmatidialAngle)
	}

	// Snell's law critical angle
	expectedSnell := math.Asin(1.34/1.37) / (math.Pi / 180.0)
	expectedCritical := 90.0 - expectedSnell
	if math.Abs(model.CriticalAngle-expectedCritical) > 1e-6 {
		t.Errorf("Expected CriticalAngle %f, got %f", expectedCritical, model.CriticalAngle)
	}

	if model.NumberOfFacets <= 0 {
		t.Errorf("Expected NumberOfFacets > 0, got %d", model.NumberOfFacets)
	}
}

func TestRunModelBlurCircleAndPointyRhabdom(t *testing.T) {
	// Compare flat lateral vs pointy lateral simulations
	paramsFlat := Parameters{
		SpeciesName:              "test_flat",
		RhabdomLength:            180,
		RhabdomWidth:             25,
		EyeDiameter:              7800,
		FacetWidth:               50,
		ApertureDiameter:         3200,
		CytoplasmRefractiveIndex: 1.34,
		RhabdomRefractiveIndex:   1.37,
		BlurCircleExtent:         18,
		ProximalRhabdomAngle:     0,
	}

	paramsPointy := Parameters{
		SpeciesName:              "test_pointy",
		RhabdomLength:            180,
		RhabdomWidth:             25,
		EyeDiameter:              7800,
		FacetWidth:               50,
		ApertureDiameter:         3200,
		CytoplasmRefractiveIndex: 1.34,
		RhabdomRefractiveIndex:   1.37,
		BlurCircleExtent:         18,
		ProximalRhabdomAngle:     12.5,
	}

	modelFlat := NewModel(paramsFlat)
	modelPointy := NewModel(paramsPointy)

	modelFlat.runModel()
	modelPointy.runModel()

	defer os.Remove("test_flat_pathlengths.csv")
	defer os.Remove("test_pointy_pathlengths.csv")

	// Read flat pathlengths file
	fileFlat, err := os.Open("test_flat_pathlengths.csv")
	if err != nil {
		t.Fatalf("Failed to open flat pathlengths file: %v", err)
	}
	defer fileFlat.Close()

	var linesFlat []string
	scanner := bufio.NewScanner(fileFlat)
	for scanner.Scan() {
		linesFlat = append(linesFlat, scanner.Text())
	}

	// Verify leading zeros exist for off-axis blur circle facets
	foundLeadingZero := false
	blockCount := 0
	for _, line := range linesFlat {
		if line == "999" {
			blockCount++
		}
		if strings.HasPrefix(line, "0,") {
			foundLeadingZero = true
		}
	}

	if !foundLeadingZero {
		t.Errorf("Expected leading zeros in blur circle facets, but none found")
	}

	if blockCount != 121 {
		t.Errorf("Expected 121 pigment blocks (11x11), got %d", blockCount)
	}

	// Read pointy pathlengths file and verify difference due to ProximalRhabdomAngle
	filePointy, err := os.Open("test_pointy_pathlengths.csv")
	if err != nil {
		t.Fatalf("Failed to open pointy pathlengths file: %v", err)
	}
	defer filePointy.Close()

	var linesPointy []string
	scannerPointy := bufio.NewScanner(filePointy)
	for scannerPointy.Scan() {
		linesPointy = append(linesPointy, scannerPointy.Text())
	}

	hasDifference := false
	for i := 0; i < len(linesFlat) && i < len(linesPointy); i++ {
		if linesFlat[i] != linesPointy[i] {
			hasDifference = true
			break
		}
	}

	if !hasDifference {
		t.Errorf("Expected differences between flat and pointy rhabdom simulations, but got identical results")
	}
}

func TestDebugFlagOutput(t *testing.T) {
	paramsNoDebug := Parameters{
		SpeciesName:              "test_nodebug",
		RhabdomLength:            100,
		RhabdomWidth:             20,
		EyeDiameter:              1000,
		FacetWidth:               20,
		ApertureDiameter:         500,
		CytoplasmRefractiveIndex: 1.34,
		RhabdomRefractiveIndex:   1.37,
		BlurCircleExtent:         1,
		ProximalRhabdomAngle:     0,
	}
	modelNoDebug := NewModel(paramsNoDebug)
	modelNoDebug.DebugMode = false
	modelNoDebug.runModel()

	defer os.Remove("test_nodebug_pathlengths.csv")
	if _, err := os.Stat("test_nodebug_debug.csv"); !os.IsNotExist(err) {
		os.Remove("test_nodebug_debug.csv")
		t.Errorf("Expected test_nodebug_debug.csv to NOT exist when DebugMode is false")
	}

	paramsDebug := Parameters{
		SpeciesName:              "test_debug",
		RhabdomLength:            100,
		RhabdomWidth:             20,
		EyeDiameter:              1000,
		FacetWidth:               20,
		ApertureDiameter:         500,
		CytoplasmRefractiveIndex: 1.34,
		RhabdomRefractiveIndex:   1.37,
		BlurCircleExtent:         1,
		ProximalRhabdomAngle:     0,
	}
	modelDebug := NewModel(paramsDebug)
	modelDebug.DebugMode = true
	modelDebug.runModel()

	defer os.Remove("test_debug_pathlengths.csv")
	defer os.Remove("test_debug_debug.csv")
	if _, err := os.Stat("test_debug_debug.csv"); os.IsNotExist(err) {
		t.Errorf("Expected test_debug_debug.csv to exist when DebugMode is true")
	}
}

func TestCalculateRessensFullMatrix(t *testing.T) {
	params := Parameters{
		SpeciesName:              "test_matrix",
		RhabdomLength:            127,
		RhabdomWidth:             15.8,
		EyeDiameter:              2480,
		FacetWidth:               22.5,
		ApertureDiameter:         870,
		CytoplasmRefractiveIndex: 1.34,
		RhabdomRefractiveIndex:   1.37,
		BlurCircleExtent:         1,
		ProximalRhabdomAngle:     0,
	}

	model := NewModel(params)
	model.runModel()
	model.calculateRessens()

	defer os.Remove("test_matrix_pathlengths.csv")
	defer os.Remove("test_matrix_summary_res.csv")
	defer os.Remove("test_matrix_summary_sen.csv")

	// Verify resolution file has 11 lines of 11 comma-separated values
	resBytes, err := os.ReadFile("test_matrix_summary_res.csv")
	if err != nil {
		t.Fatalf("Failed to read resolution summary file: %v", err)
	}

	resLines := strings.Split(strings.TrimSpace(string(resBytes)), "\n")
	if len(resLines) != 11 {
		t.Errorf("Expected 11 rows in summary resolution file, got %d", len(resLines))
	}

	for rowIdx, line := range resLines {
		parts := strings.Split(line, ",")
		if len(parts) != 11 {
			t.Errorf("Row %d: expected 11 columns, got %d", rowIdx, len(parts))
		}
		for colIdx, valStr := range parts {
			val, err := strconv.Atoi(strings.TrimSpace(valStr))
			if err != nil {
				t.Errorf("Row %d, Col %d: invalid integer value '%s'", rowIdx, colIdx, valStr)
			}
			if val <= 0 {
				t.Errorf("Row %d, Col %d: expected positive resolution value, got %d", rowIdx, colIdx, val)
			}
		}
	}
}

func TestCalculateRessensNephropsWideBlurCircle(t *testing.T) {
	params := Parameters{
		SpeciesName:              "test_nephrops_bce18",
		RhabdomLength:            180,
		RhabdomWidth:             25,
		EyeDiameter:              7800,
		FacetWidth:               50,
		ApertureDiameter:         3200,
		CytoplasmRefractiveIndex: 1.34,
		RhabdomRefractiveIndex:   1.37,
		BlurCircleExtent:         18,
		ProximalRhabdomAngle:     0,
	}

	model := NewModel(params)
	model.runModel()
	model.calculateRessens()

	defer os.Remove("test_nephrops_bce18_pathlengths.csv")
	defer os.Remove("test_nephrops_bce18_summary_res.csv")
	defer os.Remove("test_nephrops_bce18_summary_sen.csv")

	resBytes, err := os.ReadFile("test_nephrops_bce18_summary_res.csv")
	if err != nil {
		t.Fatalf("Failed to read resolution summary file: %v", err)
	}

	resLines := strings.Split(strings.TrimSpace(string(resBytes)), "\n")
	if len(resLines) != 11 {
		t.Fatalf("Expected 11 rows in summary resolution file, got %d", len(resLines))
	}

	for rowIdx, line := range resLines {
		parts := strings.Split(line, ",")
		if len(parts) != 11 {
			t.Errorf("Row %d: expected 11 columns, got %d", rowIdx, len(parts))
		}
		for colIdx, valStr := range parts {
			val, err := strconv.Atoi(strings.TrimSpace(valStr))
			if err != nil {
				t.Errorf("Row %d, Col %d: invalid integer value '%s'", rowIdx, colIdx, valStr)
			}
			if val <= 0 {
				t.Errorf("Row %d, Col %d: expected positive resolution value, got %d", rowIdx, colIdx, val)
			}
			// Verify backward extrapolation fix (diff < 0 caps at 1616)
			if (rowIdx == 0 && colIdx == 1) || (rowIdx == 1 && colIdx == 4) || (rowIdx == 2 && colIdx == 4) {
				if val != 1616 {
					t.Errorf("Row %d, Col %d: expected backward extrapolation to cap at 1616, got %d", rowIdx, colIdx, val)
				}
			}
		}
	}
}
