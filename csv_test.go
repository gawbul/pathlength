// FILE: csv_test.go
// This file contains tests for the functions in csv.go

package main

import (
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

// TestParseInputParameters checks the main scenarios for parsing the parameter file.
func TestParseInputParameters(t *testing.T) {
	t.Run("ValidFile", func(t *testing.T) {
		content := `test_species1,100,10,1000,20,500,1.3,1.4,10,0
test_species2,200,20,2000,40,1000,1.5,1.6,20,1`
		paramsList, err := parseInputParameters(writeTempFile(t, content))
		if err != nil {
			t.Errorf("parseInputParameters() returned an unexpected error: %v", err)
		}
		if len(paramsList) != 2 {
			t.Fatalf("Expected to parse 2 parameter sets, but got %d", len(paramsList))
		}
		if paramsList[0].SpeciesName != "test_species1" {
			t.Errorf("Expected SpeciesName to be 'test_species1', got '%s'", paramsList[0].SpeciesName)
		}
		if paramsList[0].RhabdomLength != 100 {
			t.Errorf("Expected RhabdomLength to be 100, got '%f'", paramsList[0].RhabdomLength)
		}
		if paramsList[1].RhabdomLength != 200 {
			t.Errorf("Expected RhabdomLength to be 200, got '%f'", paramsList[1].RhabdomLength)
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		if _, err := parseInputParameters("non_existent_file.csv"); err == nil {
			t.Error("parseInputParameters() was expected to return an error for a non-existent file")
		}
	})

	t.Run("MalformedFile", func(t *testing.T) {
		content := `good_species,1,1,1,1,1,1,1,1,1
bad_species,1,1`
		paramsList, err := parseInputParameters(writeTempFile(t, content))
		if err != nil {
			t.Errorf("parseInputParameters() returned an unexpected error: %v", err)
		}
		if len(paramsList) != 1 {
			t.Errorf("Expected to parse 1 valid line, but got %d", len(paramsList))
		}
	})

	// A non-numeric optical value used to be swallowed by a discarded ParseFloat
	// error and silently become zero.
	t.Run("NonNumericField", func(t *testing.T) {
		content := `good_species,100,10,1000,20,500,1.3,1.4,10,0
bad_species,100,10,1000,20,500,1.3,1.4,not_a_number,0`
		paramsList, err := parseInputParameters(writeTempFile(t, content))
		if err != nil {
			t.Errorf("parseInputParameters() returned an unexpected error: %v", err)
		}
		if len(paramsList) != 1 {
			t.Fatalf("Expected the non-numeric record to be rejected, got %d records", len(paramsList))
		}
		if paramsList[0].SpeciesName != "good_species" {
			t.Errorf("Expected 'good_species', got '%s'", paramsList[0].SpeciesName)
		}
	})
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "params-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(tmpfile.Name()) })
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	return tmpfile.Name()
}

// TestDepositGrowsBeyondFixedArray covers the accumulator that used to be a fixed
// 21-element array, silently discarding every rhabdom past the twenty-first.
func TestDepositGrowsBeyondFixedArray(t *testing.T) {
	var acc []float64
	acc = deposit(acc, 40, 3.5)
	if len(acc) != 41 {
		t.Fatalf("Expected the accumulator to grow to 41 entries, got %d", len(acc))
	}
	if acc[40] != 3.5 {
		t.Errorf("Expected 3.5 at offset 40, got %f", acc[40])
	}
	acc = deposit(acc, 40, 1.5)
	if acc[40] != 5.0 {
		t.Errorf("Expected deposits to accumulate to 5.0, got %f", acc[40])
	}
	// A zero deposit must not extend the profile with meaningless trailing offsets.
	before := len(acc)
	if acc = deposit(acc, 99, 0); len(acc) != before {
		t.Errorf("Expected a zero deposit to leave the length at %d, got %d", before, len(acc))
	}
}

// TestSummariseBlockResolution checks the half-maximum interpolation against
// profiles whose full width at half maximum is known exactly. The stored profile is
// area-weighted, so each offset is scaled by its ring area and the summary divides
// it back out to recover the intensity.
func TestSummariseBlockResolution(t *testing.T) {
	model := &Model{OmmatidialAngle: 2.0, NumberOfFacets: 4}

	cases := []struct {
		name     string
		psf      []float64
		wantFWHM float64
	}{
		// Falls from the peak straight to zero: the half-maximum lies midway across
		// the first step, so the half width is 0.5 ommatidial angles.
		{"SingleStep", []float64{1, 0}, 2.0 * 0.5 * 2},
		// Sits exactly on the half-maximum at offset 1, giving a half width of 1.
		{"ExactHalfAtUnitOffset", []float64{1, 0.5, 0}, 2.0 * 1.0 * 2},
		// A wider profile crossing half-maximum midway between offsets 2 and 3.
		{"WiderProfile", []float64{1, 1, 1, 0}, 2.0 * 2.5 * 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			weighted := make([]float64, len(tc.psf))
			for j, v := range tc.psf {
				weighted[j] = v * ringArea(j)
			}
			got := model.summariseBlock(weighted)
			if math.Abs(got.FWHMDegrees-tc.wantFWHM) > 1e-9 {
				t.Errorf("Expected FWHM %.6f deg, got %.6f", tc.wantFWHM, got.FWHMDegrees)
			}
			if got.PeakOffset != 0 {
				t.Errorf("Expected an on-axis peak, got offset %d", got.PeakOffset)
			}
		})
	}

	// A flat top-hat profile is measured at its own edge: rhabdoms beyond the
	// outermost illuminated one are genuinely dark, so the half maximum falls midway
	// across the final step.
	t.Run("TopHatMeasuredAtItsEdge", func(t *testing.T) {
		weighted := []float64{ringArea(0), ringArea(1), ringArea(2)}
		got := model.summariseBlock(weighted)
		want := 2.0 * 2.5 * model.OmmatidialAngle
		if math.Abs(got.FWHMDegrees-want) > 1e-9 {
			t.Errorf("Expected FWHM %.6f deg, got %.6f", want, got.FWHMDegrees)
		}
	})

	t.Run("EmptyProfile", func(t *testing.T) {
		got := model.summariseBlock(nil)
		if !math.IsNaN(got.FWHMDegrees) {
			t.Errorf("Expected an undefined FWHM for an empty profile, got %f", got.FWHMDegrees)
		}
		if got.SensitivityPercent != 0 {
			t.Errorf("Expected zero sensitivity for an empty profile, got %f", got.SensitivityPercent)
		}
	})
}

// singleFacetModel is an eye whose eyeshine patch spans exactly one facet, so the
// whole aperture is the axial ray and the summary can be checked by hand.
func singleFacetModel(t *testing.T, name string) *Model {
	t.Helper()
	model := mustModel(t, Parameters{
		SpeciesName:              name,
		RhabdomLength:            100,
		RhabdomWidth:             20,
		EyeDiameter:              1000,
		FacetWidth:               20,
		ApertureDiameter:         40,
		CytoplasmRefractiveIndex: 1.34,
		RhabdomRefractiveIndex:   1.37,
		BlurCircleExtent:         1,
		ProximalRhabdomAngle:     0,
	})
	if model.NumberOfFacets != 1 {
		t.Fatalf("Test needs a single-facet patch, got %d facets", model.NumberOfFacets)
	}
	return model
}

// TestAccumulateAndSummarise checks the summary against a hand-computed value.
func TestAccumulateAndSummarise(t *testing.T) {
	model := singleFacetModel(t, "test_summary")

	profile := model.accumulate(nil, 0, []float64{84.0})
	got := model.summariseBlock(profile)

	// The axial facet transmits fully, so the absorbed percentage is simply the
	// Beer-Lambert absorbance over 84 um. The area weight cancels against the patch
	// area for a single facet.
	wantSens := 100.0 * (1.0 - math.Exp(-absorptionCoefficient*84.0))
	if math.Abs(got.SensitivityPercent-wantSens) > 1e-9 {
		t.Errorf("Expected sensitivity %.6f%%, got %.6f%%", wantSens, got.SensitivityPercent)
	}

	// All the light lands on the axial rhabdom and the next offset out is dark, so
	// the half maximum falls midway across that single step: an acceptance angle of
	// exactly one ommatidial angle.
	if math.Abs(got.FWHMDegrees-model.OmmatidialAngle) > 1e-9 {
		t.Errorf("Expected an acceptance angle of %.6f deg, got %.6f",
			model.OmmatidialAngle, got.FWHMDegrees)
	}

	// A ray that absorbs nothing must leave the profile empty rather than depositing
	// a zero that would extend it with a meaningless offset.
	if empty := model.accumulate(nil, 0, nil); len(empty) != 0 {
		t.Errorf("Expected an empty profile for a ray with no path, got %v", empty)
	}
}

// TestCalculateRessensWritesMatrices covers the file output separately from the
// physics, which is now computed as the rays are traced rather than by reading the
// pathlengths file back.
func TestCalculateRessensWritesMatrices(t *testing.T) {
	model := singleFacetModel(t, "test_write")

	summaries := make([]blockSummary, pigmentSteps*pigmentSteps)
	for i := range summaries {
		summaries[i] = blockSummary{FWHMDegrees: float64(i), SensitivityPercent: float64(i) / 2}
	}
	if err := model.calculateRessens(summaries); err != nil {
		t.Fatalf("calculateRessens returned an unexpected error: %v", err)
	}
	defer os.Remove("test_write_summary_res.csv")
	defer os.Remove("test_write_summary_sen.csv")

	res := readMatrix(t, "test_write_summary_res.csv")
	sens := readMatrix(t, "test_write_summary_sen.csv")
	// Rows vary the shielding pigment, columns the tapetal pigment, in the order the
	// simulation produced them.
	for row := 0; row < pigmentSteps; row++ {
		for col := 0; col < pigmentSteps; col++ {
			want := float64(row*pigmentSteps + col)
			if res[row][col] != want {
				t.Errorf("res[%d][%d] = %f, expected %f", row, col, res[row][col], want)
			}
			if sens[row][col] != want/2 {
				t.Errorf("sens[%d][%d] = %f, expected %f", row, col, sens[row][col], want/2)
			}
		}
	}

	// A short or long slice means the caller and the writer disagree about the matrix
	// shape, which must be an error rather than a partially written file.
	if err := model.calculateRessens(summaries[:5]); err == nil {
		t.Error("Expected an error for the wrong number of pigment states")
	}
}

// TestSummaryMatricesAreUsable runs the reference eye end to end and checks that the
// resolution matrix is fully defined. Earlier versions passed a `value > 0` check
// even when every cell was the fabricated no-crossing fallback.
func TestSummaryMatricesAreUsable(t *testing.T) {
	params := nephropsFlatLateral("test_matrix")
	model := mustModel(t, params)

	summaries, err := model.runModel()
	if err != nil {
		t.Fatalf("runModel failed: %v", err)
	}
	if err := model.calculateRessens(summaries); err != nil {
		t.Fatalf("calculateRessens failed: %v", err)
	}
	defer os.Remove("test_matrix_pathlengths.csv")
	defer os.Remove("test_matrix_summary_res.csv")
	defer os.Remove("test_matrix_summary_sen.csv")

	res := readMatrix(t, "test_matrix_summary_res.csv")
	sens := readMatrix(t, "test_matrix_summary_sen.csv")

	for row := range res {
		for col := range res[row] {
			if math.IsNaN(res[row][col]) {
				t.Errorf("Resolution [%d][%d] is undefined", row, col)
			} else if res[row][col] <= 0 {
				t.Errorf("Resolution [%d][%d] = %f, expected a positive acceptance angle",
					row, col, res[row][col])
			}
			if s := sens[row][col]; s < 0 || s > 100 {
				t.Errorf("Sensitivity [%d][%d] = %f, expected a percentage in 0-100", row, col, s)
			}
		}
	}

	// Migrating the screening pigment across the whole rhabdom must not raise
	// sensitivity: a light-adapted eye absorbs less than a dark-adapted one.
	if sens[pigmentSteps-1][0] > sens[0][0] {
		t.Errorf("Full screening pigment raised sensitivity from %.4f%% to %.4f%%",
			sens[0][0], sens[pigmentSteps-1][0])
	}
	// Extending the tapetum must not lower sensitivity: reflected light is absorbed twice.
	if sens[0][pigmentSteps-1] < sens[0][0] {
		t.Errorf("Full tapetum lowered sensitivity from %.4f%% to %.4f%%",
			sens[0][0], sens[0][pigmentSteps-1])
	}
}

func readMatrix(t *testing.T, filename string) [][]float64 {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", filename, err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != pigmentSteps {
		t.Fatalf("%s: expected %d rows, got %d", filename, pigmentSteps, len(lines))
	}
	matrix := make([][]float64, len(lines))
	for i, line := range lines {
		fields := strings.Split(line, ",")
		if len(fields) != pigmentSteps {
			t.Fatalf("%s row %d: expected %d columns, got %d", filename, i, pigmentSteps, len(fields))
		}
		matrix[i] = make([]float64, len(fields))
		for j, field := range fields {
			v, err := strconv.ParseFloat(strings.TrimSpace(field), 64)
			if err != nil {
				t.Fatalf("%s row %d column %d: invalid value %q", filename, i, j, field)
			}
			matrix[i][j] = v
		}
	}
	return matrix
}
