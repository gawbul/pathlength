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
   - `Model` struct: Contains calculated parameters and simulation state
   - `NewModel()`: Initializes model from parameters
   - `initialCalculations()`: Computes derived values (eye radius, critical angle, facet counts, etc.)
   - `runModel()`: Main simulation loop implementing ray tracing logic
     - Iterates over pigment combinations (tapetal and shielding)
     - For each combination, iterates over facets in the eyeshine patch
     - Calculates light paths through rhabdoms based on angle of incidence
     - Four main cases: perpendicular ray, no reflection, edge reflection, base bounce
     - Outputs pathlengths to CSV file (and optional debug file when `-d` is passed)

3. **csv.go** - File I/O and data processing
   - `parseInputParameters()`: Reads CSV parameter file, returns slice of Parameters
   - `calculateRessens()`: Post-processes pathlengths to calculate resolution and sensitivity
     - Reads `{species}_pathlengths.csv`
     - Calculates absorbance, accumulates values across rhabdoms
     - Outputs `{species}_summary_res.csv` and `{species}_summary_sen.csv`

4. **csv_test.go / model_test.go** - Test suite
   - `TestParseInputParameters`: Tests CSV parsing with valid, nonexistent, and malformed files
   - `TestCalculateRessens`: Tests resolution/sensitivity calculation
   - `TestInitialCalculations`: Tests derived optical calculations
   - `TestRunModelBlurCircleAndPointyRhabdom`: Tests blur circle shifts and pointy rhabdom behavior
   - `TestDebugFlagOutput`: Tests optional debug CSV file generation
   - `TestCalculateRessensFullMatrix`: Tests resolution and sensitivity matrix generation

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
- `X_pathlengths.csv` - Raw pathlength data for each facet/pigment combination
- `X_summary_res.csv` - Calculated resolution values
- `X_summary_sen.csv` - Calculated sensitivity values
- `X_debug.csv` - (Optional, generated with `-d`) Debug information from simulation

## Key Implementation Notes

- The ray tracing simulation handles refraction at the cornea using empirical regression equations for different angle ranges
- The critical angle for total internal reflection is calculated using Snell's law
- Pigment migration is simulated by incrementing tapetal and shielding pigment values through the rhabdom length
- The blur circle feature prepends zero shifts to account for optical aberrations
- Pointy-ended rhabdoms are simulated with the proximal rhabdom acceptance angle offset

## Dependencies

- Go 1.26+ (as configured in go.mod / devcontainer)
- No external dependencies required

## Citation

When using this program, cite:
Gaten, E., Moss, S., Johnson, M. 2013. The Reniform Reflecting Superposition Compound Eyes of Nephrops Norvegicus: Optics, Susceptibility to Light-Induced Damage, Electrophysiology and a Ray Tracing Model. In: M. L. Johnson and M. P. Johnson, ed(s). Advances in Marine Biology: The Ecology and Biology of Nephrops norvegicus. Oxford: Academic Press, 107:148.
