# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

## Project Overview

PathLength is a ray tracing model that calculates the resolution and sensitivity of reflective superposition compound eyes (found in crustaceans like Nephrops norvegicus). It's a Go rewrite of an original QBASIC program from 1995, ported from a Python version.

## Development Commands

### Build
```bash
go build
```
This creates an executable named `pathlength` in the current directory.

### Run from source
```bash
go run .
```
Note: The program requires a `-f` flag with a CSV parameter file path.

### Run with parameter file
```bash
go run . -f example_data/nephrops_parameters.txt
# or with built binary:
./pathlength -f example_data/nephrops_parameters.txt
```

### Test
```bash
# Run all tests
go test -v

# Run a specific test
go test -v -run TestParseInputParameters
go test -v -run TestCalculateRessens
```

### Display help
```bash
./pathlength -h
```

## Development Workflow

**IMPORTANT: Always run tests before committing changes.**

This project uses pre-commit hooks to automatically enforce code quality. The hooks will:
- Format Go code with `go fmt`
- Run static analysis with `go vet`
- Run all tests with `go test`
- Ensure `go.mod` and `go.sum` are tidy
- Check for common issues (merge conflicts, large files, etc.)

### Setup pre-commit (first time only)

```bash
# Install pre-commit (if not already installed)
pip install pre-commit
# or: brew install pre-commit

# Install the git hooks
pre-commit install
```

### Manual test run

You can manually run tests before committing:

```bash
go test -v
```

If tests fail, fix the issues before committing. This ensures the codebase remains stable and the scientific calculations remain accurate.

### Commit Message Format

Use conventional commit format for all commit messages:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring (no functional changes)
- `test`: Adding or updating tests
- `docs`: Documentation changes
- `chore`: Maintenance tasks (dependencies, build config, etc.)
- `perf`: Performance improvements

**Examples:**
```
feat(model): add support for pointy-ended rhabdoms
fix(csv): handle malformed records gracefully
test(parser): add edge case tests for parameter validation
docs: update README with installation instructions
refactor(model): simplify ray tracing calculation logic
```

## Code Architecture

### Core Components

The codebase is organized into Go files in a single package:

1. **pathlength.go** - Entry point with CLI argument parsing
   - Handles flags: `-f` (file), `-d` (debug), `-h` (help), `-v` (version), `-c` (citation), `-l` (license)
   - Parses parameter file and orchestrates simulation runs
   - Main loop iterates over each parameter set from CSV input

2. **model.go** - Simulation engine and data structures
   - `Parameters` struct: Holds 10 eye-specific configuration values from CSV
   - `Model` struct: Derived optics (ommatidial angle, facet count, critical angle)
   - `NewModel()`: Validates the parameters and computes the derived optics. Returns
     an error for any eye that is not physically realisable, rather than letting NaNs
     propagate into the results
   - `refractedAngle()`: Empirical corneal refraction regression
   - `facetTransmission()`: Light a facet admits at a given incidence, as a flux
     factor in [0, 1]. Applied to absorbed intensity, never to path length
   - `blurOffset()`: Continuous displacement in rhabdoms of the image formed by a
     facet, spanning 0 to (BlurCircleExtent - 1) across the eyeshine patch
   - `traceRay()`: Follows one ray through the rhabdom array. Four cases: axial ray,
     no reflection, wall reflection, base bounce. Returns raw geometry in µm, the
     terminating case, and the largest angle reached
   - `runModel()`: Iterates the 121 pigment states, writes the pathlengths CSV, and
     returns the per-state summaries. The absorption profile is accumulated as the
     rays are traced rather than by reading the file back, so the summary does not
     depend on the output format at all

3. **csv.go** - File I/O and summary
   - `parseInputParameters()`: Reads CSV parameter file, returns slice of Parameters
   - `ringArea()`: Area of the annulus of rhabdoms at a given offset. Used both to
     weight contributing facets and to normalise the light arriving at an offset
   - `deposit()`: Accumulates into a growable profile, so no rhabdom is ever silently
     dropped
   - `summariseBlock()`: Turns one block's area-weighted profile into an acceptance
     angle (FWHM, degrees) and a sensitivity (percent absorbed)
   - `accumulate()`: Adds one facet's traced ray into the area-weighted absorption
     profile
   - `calculateRessens()`: Writes the two summary matrices from the per-state
     summaries accumulated during the simulation

4. **csv_test.go / model_test.go** - Test suite
   - `TestInitialCalculations`: Derived optical calculations
   - `TestNewModelRejectsUnphysicalParameters`: Refractive index ordering, aperture
     against eye diameter, blur circle against facet count
   - `TestBlurOffsetSpansExtentEvenly`: Blur mapping is uniform and tie-free
   - `TestRaysStayWithinPhysicalGeometry`: Total path is bounded by
     2 * RhabdomLength / cos(maxAngle), and no ray exceeds 90 degrees to the axis
   - `TestPathlengthsAreRawGeometry`: Facet transmission is not folded into geometry
   - `TestGuidedRayIgnoresScreeningPigment`: Total internal reflection versus the
     proximal screening pigment
   - `TestDepositGrowsBeyondFixedArray`: Profile accumulation has no fixed ceiling
   - `TestSummariseBlockResolution`: Half-maximum interpolation against known profiles
   - `TestCalculateRessens`: Beer-Lambert absorbance for a single-facet patch
   - `TestSummaryMatricesAreUsable`: Full matrices are defined, in range, and respond
     correctly to pigment migration
   - `TestDebugFlagOutput`: Per-ray debug CSV

### Data Flow

1. CSV parameter file → `parseInputParameters()` → slice of `Parameters`
2. For each Parameters → `NewModel()` → `Model` with calculated values
3. `Model.runModel()` → writes `{species}_pathlengths.csv` (and `{species}_debug.csv` if `-d` is set)
4. `Model.calculateRessens()` → reads pathlengths → writes summary CSV files

### Input File Format

CSV files in `example_data/` directory with 10 comma-separated fields per row:
```
species_name, rhabdom_length, rhabdom_width, eye_diameter, facet_width,
aperture_diameter, cytoplasm_refractive_index, rhabdom_refractive_index,
blur_circle_extent, proximal_rhabdom_angle
```

Each row represents a separate simulation run. Multiple rows can be processed in one execution.

### Output Files

For each parameter set with species name "X", generates:
- `X_pathlengths.csv` - Raw ray geometry, in µm: a rectangular CSV with the header
  `block,shielding_um,tapetal_um,facet,rhabdom,pathlength_um` and one row per rhabdom
  entered. No block terminator and no positional state
- `X_summary_res.csv` - Acceptance angle (FWHM of the point spread function), degrees
- `X_summary_sen.csv` - Incident light absorbed, percent (0-100)
- `X_debug.csv` - (Optional, generated with `-d`) One row per traced ray

Both summary files are 11x11: **rows vary the shielding (proximal screening) pigment,
columns vary the tapetal (reflecting) pigment**, each from fully retracted to fully
covering the rhabdom. A resolution cell reading `NaN` means the profile never falls to
half its maximum, so the acceptance angle is undefined for that state.

## Key Implementation Notes

- Corneal refraction uses empirical regression equations over four angle ranges
- Ray angles (`boa`) are measured from the rhabdom axis, so the angle at the wall
  normal is (90 - boa) and light is guided while `boa < CriticalAngle`, where
  `CriticalAngle = 90 - asin(n_cytoplasm / n_rhabdom)`
- Pigment migration is simulated over 11 positions per pigment, giving 121 states
- The blur circle displaces each facet's image by a **continuous** offset, spanning 0
  to (BlurCircleExtent - 1) across the eyeshine patch. Light is split between the two
  whole-rhabdom offsets that bracket it. Quantising this to whole rhabdoms aliases the
  facets unevenly and cuts notches into the profile that the half-maximum search then
  locks onto
- The blur circle extent may not exceed the facet count across the patch
- Pointy-ended rhabdoms are simulated with the proximal rhabdom acceptance angle offset
- Rays reaching 90 degrees or more to the rhabdom axis cannot advance towards the
  proximal end and are discarded with a warning, rather than being folded back by
  taking absolute values of `tan` and `cos`
- The absorption coefficient is fixed at 0.01 µm⁻¹ and tapetal reflectance at 1.0

### Known limitation

A trace terminates at the first reflection rather than following the reflected ray
onward. Because an unreflected ray continues leaking into adjacent rhabdoms and
absorbing there, extending the tapetum can occasionally shorten the total absorbing
path and lower the reported sensitivity, by up to about 6 percentage points. A tapetal
mirror can only add path length in reality, so this is a limitation of the 1995 case
structure rather than a property of the eye.

## Dependencies

- Go 1.26+ (as configured in go.mod / devcontainer)
- No external dependencies required

## Citation

When using this program, cite:
Gaten, E., Moss, S., Johnson, M. 2013. The Reniform Reflecting Superposition Compound Eyes of Nephrops Norvegicus: Optics, Susceptibility to Light-Induced Damage, Electrophysiology and a Ray Tracing Model. In: M. L. Johnson and M. P. Johnson, ed(s). Advances in Marine Biology: The Ecology and Biology of Nephrops norvegicus. Oxford: Academic Press, 107:148.
