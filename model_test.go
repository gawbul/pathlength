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

// nephropsFlatLateral is the reference parameter set used throughout the tests.
func nephropsFlatLateral(name string) Parameters {
	return Parameters{
		SpeciesName:              name,
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
}

func mustModel(t *testing.T, p Parameters) *Model {
	t.Helper()
	m, err := NewModel(p)
	if err != nil {
		t.Fatalf("NewModel(%s) returned an unexpected error: %v", p.SpeciesName, err)
	}
	return m
}

func TestInitialCalculations(t *testing.T) {
	params := nephropsFlatLateral("test_eye")
	params.ProximalRhabdomAngle = 12.5
	model := mustModel(t, params)

	expectedCircumference := math.Pi * 7800.0
	if math.Abs(model.CircumferenceOfEye-expectedCircumference) > 1e-6 {
		t.Errorf("Expected CircumferenceOfEye %f, got %f", expectedCircumference, model.CircumferenceOfEye)
	}

	expectedOmmatidialAngle := (50.0 / expectedCircumference) * 360.0
	if math.Abs(model.OmmatidialAngle-expectedOmmatidialAngle) > 1e-6 {
		t.Errorf("Expected OmmatidialAngle %f, got %f", expectedOmmatidialAngle, model.OmmatidialAngle)
	}

	// boa is measured from the rhabdom axis, so the angle at the wall normal is
	// (90 - boa) and light is guided while boa < CriticalAngle.
	expectedCritical := 90.0 - math.Asin(1.34/1.37)*radToDegConv
	if math.Abs(model.CriticalAngle-expectedCritical) > 1e-6 {
		t.Errorf("Expected CriticalAngle %f, got %f", expectedCritical, model.CriticalAngle)
	}

	if model.NumberOfFacets != 33 {
		t.Errorf("Expected 33 facets across the eyeshine patch, got %d", model.NumberOfFacets)
	}
}

// TestNewModelRejectsUnphysicalParameters covers inputs that previously produced
// NaNs which then silently disabled the total-internal-reflection test or emptied
// the simulation, in both cases without any warning to the user.
func TestNewModelRejectsUnphysicalParameters(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Parameters)
		want   string
	}{
		{"CytoplasmIndexExceedsRhabdom", func(p *Parameters) {
			p.CytoplasmRefractiveIndex = 1.40
			p.RhabdomRefractiveIndex = 1.37
		}, "total internal reflection"},
		{"EqualRefractiveIndices", func(p *Parameters) {
			p.RhabdomRefractiveIndex = p.CytoplasmRefractiveIndex
		}, "total internal reflection"},
		{"ApertureExceedsEye", func(p *Parameters) {
			p.ApertureDiameter = 9000
		}, "aperture diameter"},
		{"ZeroRhabdomLength", func(p *Parameters) { p.RhabdomLength = 0 }, "rhabdom length"},
		{"NegativeFacetWidth", func(p *Parameters) { p.FacetWidth = -50 }, "facet width"},
		{"BlurCircleBelowOne", func(p *Parameters) { p.BlurCircleExtent = 0 }, "blur circle extent"},
		// astacodes ships with an 18-rhabdom blur circle but only 7 facets across the
		// eyeshine patch, which leaves 11 rhabdom offsets receiving no light at all.
		{"BlurCircleExceedsFacets", func(p *Parameters) {
			p.RhabdomLength, p.RhabdomWidth = 84, 16
			p.EyeDiameter, p.FacetWidth, p.ApertureDiameter = 890, 32, 445
			p.BlurCircleExtent = 18
		}, "exceeds the 7 facets"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			params := nephropsFlatLateral("test_invalid")
			tc.mutate(&params)
			m, err := NewModel(params)
			if err == nil {
				t.Fatalf("Expected an error, got a model with %d facets", m.NumberOfFacets)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("Expected error mentioning %q, got: %v", tc.want, err)
			}
		})
	}
}

// TestBlurOffsetSpansExtentEvenly guards the blur-circle mapping. The previous
// `facet > fd*i` formulation aliased facets unevenly onto whole rhabdom offsets and
// skipped an offset entirely wherever fd*i landed on an exact integer, cutting a
// notch into the profile that the half-maximum search then locked onto.
func TestBlurOffsetSpansExtentEvenly(t *testing.T) {
	model := mustModel(t, nephropsFlatLateral("test_blur"))

	if got := model.blurOffset(0); got != 0 {
		t.Errorf("Expected the central facet to be undisplaced, got %f", got)
	}
	want := model.Params.BlurCircleExtent - 1
	if got := model.blurOffset(model.NumberOfFacets - 1); math.Abs(got-want) > 1e-9 {
		t.Errorf("Expected the outermost facet at offset %f, got %f", want, got)
	}

	// The mapping must be strictly increasing with a constant step, so no offset is
	// starved of contributing facets and no facet index sits on a tie boundary.
	step := model.blurOffset(1) - model.blurOffset(0)
	if step <= 0 {
		t.Fatalf("Expected a positive blur step, got %f", step)
	}
	for facet := 1; facet < model.NumberOfFacets; facet++ {
		delta := model.blurOffset(facet) - model.blurOffset(facet-1)
		if math.Abs(delta-step) > 1e-9 {
			t.Errorf("Facet %d: expected a uniform blur step of %f, got %f", facet, step, delta)
		}
	}

	// A single-rhabdom blur circle means no displacement at all.
	unblurred := nephropsFlatLateral("test_noblur")
	unblurred.BlurCircleExtent = 1
	plain := mustModel(t, unblurred)
	for facet := 0; facet < plain.NumberOfFacets; facet++ {
		if got := plain.blurOffset(facet); got != 0 {
			t.Errorf("Facet %d: expected no displacement with a 1-rhabdom blur circle, got %f", facet, got)
		}
	}
}

// TestRaysStayWithinPhysicalGeometry guards against the ray tracer folding rays back
// on themselves. Taking |tan| and |cos| of an angle past 90 degrees used to yield
// path lengths many times the rhabdom length for a ray that cannot in fact advance
// towards the proximal end at all.
func TestRaysStayWithinPhysicalGeometry(t *testing.T) {
	for _, params := range []Parameters{
		nephropsFlatLateral("test_geom_flat"),
		func() Parameters {
			p := nephropsFlatLateral("test_geom_pointy")
			p.ProximalRhabdomAngle = 12.5
			return p
		}(),
	} {
		model := mustModel(t, params)
		increment := params.RhabdomLength / 10.0

		for pStep := 0; pStep < pigmentSteps; pStep++ {
			for tStep := 0; tStep < pigmentSteps; tStep++ {
				for facet := 0; facet < model.NumberOfFacets; facet++ {
					trace := model.traceRay(facet, float64(pStep)*increment, float64(tStep)*increment)
					if trace.Lost {
						t.Fatalf("%s: facet %d lost the ray at pigment step (%d,%d)",
							params.SpeciesName, facet, pStep, tStep)
					}
					if trace.MaxAngle < 0 || trace.MaxAngle >= maxPropagationAngle {
						t.Fatalf("%s: facet %d reached %.2f deg to the rhabdom axis",
							params.SpeciesName, facet, trace.MaxAngle)
					}

					// Every segment covers some axial depth at an angle no greater
					// than MaxAngle, and the ray can traverse the rhabdom at most
					// twice (down, then back off the tapetum). So the total path is
					// bounded by 2*RhabdomLength/cos(MaxAngle). Taking |cos| of an
					// angle past 90 degrees used to break this by a factor of ten.
					limit := 2.0*params.RhabdomLength/math.Cos(trace.MaxAngle*degToRadConv) + 1e-9
					total := 0.0
					for i, v := range trace.Pathlengths {
						if math.IsNaN(v) || v < 0 {
							t.Fatalf("%s: facet %d rhabdom %d has an invalid pathlength %f",
								params.SpeciesName, facet, i, v)
						}
						total += v
					}
					if total > limit {
						t.Errorf("%s: facet %d at pigment step (%d,%d) traced %.1f um, "+
							"exceeding the %.1f um bound for a maximum angle of %.2f deg",
							params.SpeciesName, facet, pStep, tStep, total, limit, trace.MaxAngle)
					}
				}
			}
		}
	}
}

// TestPathlengthsAreRawGeometry confirms that facet transmission is no longer folded
// into the geometry. The axial ray must traverse exactly the rhabdom length, and
// twice that when the tapetum reflects it back with no screening pigment in the way.
func TestPathlengthsAreRawGeometry(t *testing.T) {
	params := nephropsFlatLateral("test_geometry")
	model := mustModel(t, params)

	single := model.traceRay(0, 0, 0)
	if len(single.Pathlengths) != 1 || math.Abs(single.Pathlengths[0]-params.RhabdomLength) > 1e-9 {
		t.Errorf("Expected the axial ray to traverse %.1f um once, got %v",
			params.RhabdomLength, single.Pathlengths)
	}

	reflected := model.traceRay(0, 0, params.RhabdomLength)
	if len(reflected.Pathlengths) != 1 || math.Abs(reflected.Pathlengths[0]-2*params.RhabdomLength) > 1e-9 {
		t.Errorf("Expected the tapetum to double the axial path to %.1f um, got %v",
			2*params.RhabdomLength, reflected.Pathlengths)
	}
}

// TestGuidedRayIgnoresScreeningPigment covers the proximal screening pigment rule.
// The pigment lies in the cytoplasm outside the rhabdom, so it can only absorb light
// that has already crossed the wall; a totally internally reflected ray stays inside
// and must not be truncated by it.
func TestGuidedRayIgnoresScreeningPigment(t *testing.T) {
	params := nephropsFlatLateral("test_guided")
	// A narrow blur circle keeps the entry angle below the critical angle so the ray
	// is guided rather than leaking through the wall.
	params.BlurCircleExtent = 1
	model := mustModel(t, params)

	facet := 1
	boa := refractedAngle(float64(facet)*model.OmmatidialAngle) + model.blurOffset(facet)*model.OmmatidialAngle
	if boa >= model.CriticalAngle {
		t.Fatalf("Test needs a guided ray: boa %.4f is not below the critical angle %.4f",
			boa, model.CriticalAngle)
	}

	unscreened := model.traceRay(facet, 0, 0)
	screened := model.traceRay(facet, params.RhabdomLength, 0)
	if len(unscreened.Pathlengths) != 1 || len(screened.Pathlengths) != 1 {
		t.Fatalf("Expected single-segment traces, got %v and %v",
			unscreened.Pathlengths, screened.Pathlengths)
	}
	if math.Abs(unscreened.Pathlengths[0]-screened.Pathlengths[0]) > 1e-9 {
		t.Errorf("Screening pigment truncated a guided ray: %.6f um without, %.6f um with",
			unscreened.Pathlengths[0], screened.Pathlengths[0])
	}
}

func TestRunModelProducesWellFormedBlocks(t *testing.T) {
	flat := nephropsFlatLateral("test_flat")
	pointy := nephropsFlatLateral("test_pointy")
	pointy.ProximalRhabdomAngle = 12.5

	modelFlat := mustModel(t, flat)
	modelPointy := mustModel(t, pointy)

	if _, err := modelFlat.runModel(); err != nil {
		t.Fatalf("runModel(flat) failed: %v", err)
	}
	if _, err := modelPointy.runModel(); err != nil {
		t.Fatalf("runModel(pointy) failed: %v", err)
	}
	defer os.Remove("test_flat_pathlengths.csv")
	defer os.Remove("test_pointy_pathlengths.csv")

	linesFlat := readLines(t, "test_flat_pathlengths.csv")
	linesPointy := readLines(t, "test_pointy_pathlengths.csv")

	// The file is a plain rectangular CSV: a header, then one row per rhabdom, each
	// carrying its own keys. There is no block terminator and no positional state.
	if linesFlat[0] != pathlengthsHeader {
		t.Fatalf("Expected the header %q, got %q", pathlengthsHeader, linesFlat[0])
	}
	if strings.Contains(strings.Join(linesFlat, "\n"), "\n999\n") {
		t.Error("Expected no 999 block terminator in the output")
	}

	// Every (block, facet) pair must appear exactly once as a group, with rhabdom
	// indices running from zero.
	type key struct{ block, facet int }
	seen := map[key]int{}
	for i, line := range linesFlat[1:] {
		fields := strings.Split(line, ",")
		if len(fields) != 6 {
			t.Fatalf("Row %d: expected 6 fields, got %d: %q", i, len(fields), line)
		}
		for j, field := range fields {
			if _, err := strconv.ParseFloat(field, 64); err != nil {
				t.Fatalf("Row %d field %d is not numeric: %q", i, j, field)
			}
		}
		block, _ := strconv.Atoi(fields[0])
		facet, _ := strconv.Atoi(fields[3])
		rhabdom, _ := strconv.Atoi(fields[4])
		k := key{block, facet}
		if rhabdom != seen[k] {
			t.Fatalf("Block %d facet %d: expected rhabdom %d, got %d", block, facet, seen[k], rhabdom)
		}
		seen[k]++
	}
	if want := pigmentSteps * pigmentSteps * modelFlat.NumberOfFacets; len(seen) != want {
		t.Errorf("Expected %d block/facet groups, got %d", want, len(seen))
	}

	same := len(linesFlat) == len(linesPointy)
	if same {
		for i := range linesFlat {
			if linesFlat[i] != linesPointy[i] {
				same = false
				break
			}
		}
	}
	if same {
		t.Error("Expected the proximal rhabdom angle to change the traced pathlengths")
	}
}

func TestDebugFlagOutput(t *testing.T) {
	params := nephropsFlatLateral("test_nodebug")
	modelNoDebug := mustModel(t, params)
	modelNoDebug.DebugMode = false
	if _, err := modelNoDebug.runModel(); err != nil {
		t.Fatalf("runModel failed: %v", err)
	}
	defer os.Remove("test_nodebug_pathlengths.csv")
	if _, err := os.Stat("test_nodebug_debug.csv"); !os.IsNotExist(err) {
		os.Remove("test_nodebug_debug.csv")
		t.Errorf("Expected test_nodebug_debug.csv to NOT exist when DebugMode is false")
	}

	params.SpeciesName = "test_debug"
	modelDebug := mustModel(t, params)
	modelDebug.DebugMode = true
	if _, err := modelDebug.runModel(); err != nil {
		t.Fatalf("runModel failed: %v", err)
	}
	defer os.Remove("test_debug_pathlengths.csv")
	defer os.Remove("test_debug_debug.csv")

	lines := readLines(t, "test_debug_debug.csv")
	// One header plus one row per traced ray, rather than the block headings alone
	// that earlier versions emitted.
	wantRows := 1 + pigmentSteps*pigmentSteps*modelDebug.NumberOfFacets
	if len(lines) != wantRows {
		t.Errorf("Expected %d debug rows, got %d", wantRows, len(lines))
	}
	if !strings.HasPrefix(lines[0], "block,shielding_um,tapetal_um,facet,") {
		t.Errorf("Expected a descriptive debug header, got %q", lines[0])
	}
	if fields := strings.Split(lines[1], ","); len(fields) != 12 {
		t.Errorf("Expected 12 debug fields per ray, got %d: %q", len(fields), lines[1])
	}
}

func readLines(t *testing.T, filename string) []string {
	t.Helper()
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", filename, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Failed to read %s: %v", filename, err)
	}
	return lines
}
