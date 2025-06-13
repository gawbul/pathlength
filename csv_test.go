// FILE: csv_test.go
// This file contains tests for the functions in csv.go

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestParseInputParameters checks the main scenarios for parsing the parameter file.
func TestParseInputParameters(t *testing.T) {
	t.Run("ValidFile", func(t *testing.T) {
		// Create a temporary valid CSV file with multiple entries
		content := `test_species1,100,10,1000,20,500,1.3,1.4,10,0
test_species2,200,20,2000,40,1000,1.5,1.6,20,1`
		tmpfile, err := os.CreateTemp("", "params-*.csv")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}

		paramsList, err := parseInputParameters(tmpfile.Name())
		if err != nil {
			t.Errorf("parseInputParameters() returned an unexpected error: %v", err)
		}

		if len(paramsList) != 2 {
			t.Fatalf("Expected to parse 2 parameter sets, but got %d", len(paramsList))
		}

		// Check first record
		if paramsList[0].SpeciesName != "test_species1" {
			t.Errorf("Expected SpeciesName to be 'test_species1', got '%s'", paramsList[0].SpeciesName)
		}
		if paramsList[0].RhabdomLength != 100 {
			t.Errorf("Expected RhabdomLength to be 100, got '%f'", paramsList[0].RhabdomLength)
		}

		// Check second record
		if paramsList[1].SpeciesName != "test_species2" {
			t.Errorf("Expected SpeciesName to be 'test_species2', got '%s'", paramsList[1].SpeciesName)
		}
		if paramsList[1].RhabdomLength != 200 {
			t.Errorf("Expected RhabdomLength to be 200, got '%f'", paramsList[1].RhabdomLength)
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := parseInputParameters("non_existent_file.csv")
		if err == nil {
			t.Error("parseInputParameters() was expected to return an error for a non-existent file, but it didn't")
		}
	})

	t.Run("MalformedFile", func(t *testing.T) {
		// Create a temporary malformed CSV file with one valid and one invalid line
		content := `good_species,1,1,1,1,1,1,1,1,1
bad_species,1,1`
		tmpfile, err := os.CreateTemp("", "params-*.csv")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}

		// It should process the good line and log an error for the bad one, but not return a fatal error.
		paramsList, err := parseInputParameters(tmpfile.Name())
		if err != nil {
			t.Errorf("parseInputParameters() returned an unexpected error for a recoverable issue: %v", err)
		}
		if len(paramsList) != 1 {
			t.Errorf("Expected to parse 1 valid line, but got %d", len(paramsList))
		}
		if paramsList[0].SpeciesName != "good_species" {
			t.Errorf("Expected to parse 'good_species', but got '%s'", paramsList[0].SpeciesName)
		}
	})
}

// TestCalculateRessens checks if the summary data is calculated correctly.
func TestCalculateRessens(t *testing.T) {
	// Setup a model with some basic parameters
	params := Parameters{
		SpeciesName:   "test_summary",
		RhabdomLength: 100,
		EyeDiameter:   1000,
		FacetWidth:    20,
	}
	model := NewModel(params)

	// Create a mock pathlengths file
	pathlengthsContent := `0.000000
0.000000
84.000000,998
999
`
	pathlengthsFileName := fmt.Sprintf("%s_pathlengths.csv", params.SpeciesName)
	if err := os.WriteFile(pathlengthsFileName, []byte(pathlengthsContent), 0644); err != nil {
		t.Fatalf("Failed to write mock pathlengths file: %v", err)
	}
	defer os.Remove(pathlengthsFileName)

	// Run the function to be tested
	model.calculateRessens()

	// Check if the output files were created
	resFileName := fmt.Sprintf("%s_summary_res.csv", params.SpeciesName)
	sensFileName := fmt.Sprintf("%s_summary_sen.csv", params.SpeciesName)
	defer os.Remove(resFileName)
	defer os.Remove(sensFileName)

	if _, err := os.Stat(resFileName); os.IsNotExist(err) {
		t.Fatalf("Expected resolution file '%s' was not created", resFileName)
	}
	if _, err := os.Stat(sensFileName); os.IsNotExist(err) {
		t.Fatalf("Expected sensitivity file '%s' was not created", sensFileName)
	}

	// Check the content of the output files
	resContent, err := os.ReadFile(resFileName)
	if err != nil {
		t.Fatalf("Failed to read resolution file: %v", err)
	}
	expectedRes := "228"
	if !strings.Contains(string(resContent), expectedRes) {
		t.Errorf("Resolution file content is wrong. Got '%s', expected to contain '%s'", string(resContent), expectedRes)
	}

	sensContent, err := os.ReadFile(sensFileName)
	if err != nil {
		t.Fatalf("Failed to read sensitivity file: %v", err)
	}
	expectedSens := "56"
	if !strings.Contains(string(sensContent), expectedSens) {
		t.Errorf("Sensitivity file content is wrong. Got '%s', expected to contain '%s'", string(sensContent), expectedSens)
	}
}
