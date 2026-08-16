# Superposition Eye Pathlength Program

PathLength implements a ray tracing model to calculate the resolution and sensitivity of reflective superposition compound eyes.

Original QBASIC version by Dr Magnus L Johnson and Genevre Parker, 1995

Golang rewrite by Dr Stephen P Moss, 2025

Author: Dr Stephen P Moss

Website: [https://www.gawbul.io](https://www.gawbul.io)

Email: gawbul@gmail.com

## Install Go compiler

### macOS

```bash
# Install brew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"

# Install Golang
brew install go@1.26
```

### Linux (Debian/Ubuntu)

```bash
# Download tarball
wget -P /tmp https://dl.google.com/go/go1.26.6.linux-amd64.tar.gz

# Extract tarball
sudo tar -zxvf /tmp/go1.26.6.linux-amd64.tar.gz -C /usr/local

# Setup environment
echo "export GOROOT=/usr/local/go" >> ~/.profile
echo "export GOPATH=$HOME/go" >> ~/.profile
echo "export PATH=$GOPATH/bin:$GOROOT/bin:$PATH" >> ~/.profile

# Source environment
source ~/.profile
```

## Install dependencies

The application uses standard Go library packages with no external runtime dependencies.

### macOS

```bash
# Install git
brew install git
```

### Linux (Debian/Ubuntu)

```bash
# Install git
sudo apt-get update
sudo apt-get install git
```

## Checkout source code

```bash
# Create and change to projects directory
mkdir -p ~/projects
cd ~/projects

# Clone the GitHub repository
git clone git@github.com:gawbul/pathlength.git
```

## Usage

### Run the program from source

```bash
cd ~/projects/pathlength
go run .
```

Outputs:

```bash
Usage: pathlength -f filename [-c] [-d] [-h] [-l] [-v]
  -c    Show the program citation.
  -d    Generate debug CSV output file.
  -f string
        Path to a parameter file (CSV format). (Required)
  -h    Show this help message.
  -l    Show the program license.
  -v    Show program version.
2025/06/13 14:58:20 Error: No parameter file supplied. Use the -f flag to specify a file.
exit status 1
```

### Run the test suite

```bash
cd ~/projects/pathlength
go test -v
```

Outputs:

```bash
=== RUN   TestParseInputParameters
=== RUN   TestParseInputParameters/ValidFile
=== RUN   TestParseInputParameters/NonExistentFile
=== RUN   TestParseInputParameters/MalformedFile
2026/08/15 13:06:29 Skipping malformed record on line 2: record on line 2: wrong number of fields
--- PASS: TestParseInputParameters (0.00s)
    --- PASS: TestParseInputParameters/ValidFile (0.00s)
    --- PASS: TestParseInputParameters/NonExistentFile (0.00s)
    --- PASS: TestParseInputParameters/MalformedFile (0.00s)
=== RUN   TestCalculateRessens
INFO: Calculating resolution and sensitivity...
--- PASS: TestCalculateRessens (0.00s)
=== RUN   TestInitialCalculations
--- PASS: TestInitialCalculations (0.00s)
=== RUN   TestRunModelBlurCircleAndPointyRhabdom
--- PASS: TestRunModelBlurCircleAndPointyRhabdom (0.01s)
=== RUN   TestDebugFlagOutput
--- PASS: TestDebugFlagOutput (0.00s)
=== RUN   TestCalculateRessensFullMatrix
INFO: Calculating resolution and sensitivity...
--- PASS: TestCalculateRessensFullMatrix (0.00s)
PASS
ok  	pathlength	0.215s
```

### Build the program

```bash
cd ~/projects/pathlength
go build
chmod +x pathlength
```

### Display program usage

```bash
./pathlength -h
```

Outputs:

```bash
Usage: pathlength -f filename [-c] [-d] [-h] [-l] [-v]
  -c    Show the program citation.
  -d    Generate debug CSV output file.
  -f string
        Path to a parameter file (CSV format). (Required)
  -h    Show this help message.
  -l    Show the program license.
  -v    Show program version.
```

*Also displays if you don't pass in any arguments, as it expects a filename as input.*

### Display citation information

```bash
./pathlength -c
```

Outputs:

```bash
Gaten, E., Moss, S., Johnson, M. 2013. The Reniform Reflecting Superposition Compound Eyes of Nephrops Norvegicus:
Optics, Susceptibility to Light-Induced Damage, Electrophysiology and a Ray Tracing Model. In: M. L. Johnson and M. P. Johnson, ed(s).
Advances in Marine Biology: The Ecology and Biology of Nephrops norvegicus. Oxford: Academic Press, 107:148.
```

### Display license information

```bash
./pathlength -l
```

Outputs:

```bash
pathlength - calculates resolution and sensitivity in reflective superposition compound eyes.

Copyright (C) 2020 Dr Stephen P Moss

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
```

### Display program version

```bash
./pathlength -v
```

Outputs:

```bash
pathlength version 0.6.0
```

## Run the program

```bash
./pathlength -f example_data/acanthephyra_parameters.txt
```

Outputs:

```bash
Parsing input parameters from example_data/acanthephyra_parameters.txt...
--- Running simulation for acanthephyra ---
Calculating pathlengths for acanthephyra...
P: 0.00, T: 0.00
P: 0.00, T: 12.70
P: 0.00, T: 25.40
P: 0.00, T: 38.10
P: 0.00, T: 50.80
P: 0.00, T: 63.50
P: 0.00, T: 76.20
P: 0.00, T: 88.90
P: 0.00, T: 101.60
P: 0.00, T: 114.30
P: 0.00, T: 127.00
P: 12.70, T: 0.00
P: 12.70, T: 12.70
P: 12.70, T: 25.40
P: 12.70, T: 38.10
...
P: 127.00, T: 101.60
P: 127.00, T: 114.30
P: 127.00, T: 127.00
INFO: Calculating resolution and sensitivity...
--- Finished simulation for acanthephyra ---

--- Running simulation for acanthephyra_bce3 ---
Calculating pathlengths for acanthephyra_bce3...
...
INFO: Calculating resolution and sensitivity...
--- Finished simulation for acanthephyra_bce3 ---

--- Running simulation for acanthephyra_bce6 ---
Calculating pathlengths for acanthephyra_bce6...
...
INFO: Calculating resolution and sensitivity...
--- Finished simulation for acanthephyra_bce6 ---

All simulations complete.
```

### Run with debug output

To generate optional `{species}_debug.csv` files containing simulation debug logs (e.g. pigment migration steps):

```bash
./pathlength -f example_data/acanthephyra_parameters.txt -d
```

## Required parameters

A CSV format file is required as input to the program. You can provide multiple lines for separate runs of the model. The format should be as follows:

```text
nephropsfl,180,25,7800,50,3200,1.34,1.37,18,0
nephropspl,180,25,7800,50,3200,1.34,1.37,18,12.5
nephropsfa,180,25,6760,50,3060,1.34,1.37,10,0
nephropspa,180,25,6760,50,3060,1.34,1.37,10,12.5
```

Each row is comprised of the following fields:

```text
genus	= A prefix for the output filenames e.g. organism genus name (lowercase alphanumeric only)
180 	= Rhabdom Length
25 	= Rhabdom Width
7800 	= Eye Diameter
50 	= Facet Width
3200	= Aperture Diameter
1.34	= Cytoplasm Refractive Index
1.37	= Rhabdom Refractive Index
18	= Blur Circle Extent
0	= Proximal Rhabdom Angle (used to create pointy-ended rhabdoms)
```

*NB: The genus name is NOT case sensitive. It is always converted to lowercase and should be unique to avoid filename conflicts.*

## Output files

The following output files are created:

* `genus_pathlengths.csv` - Raw ray geometry for each facet and pigment combination
* `genus_summary_res.csv` - Acceptance angle matrix
* `genus_summary_sen.csv` - Sensitivity matrix
* `genus_debug.csv` - (Optional) Per-ray trace, enabled with `-d`

### `genus_pathlengths.csv`

A plain rectangular CSV with a header row and one row per rhabdom entered. Every row
carries its own keys, so there is no positional state and no block terminator:

```csv
block,shielding_um,tapetal_um,facet,rhabdom,pathlength_um
0,0.000000,0.000000,0,0,180.000000
0,0.000000,0.000000,1,0,180.032715
...
0,0.000000,0.000000,12,0,55.332562
0,0.000000,0.000000,12,1,52.437957
```

| Column | Meaning |
| --- | --- |
| `block` | Pigment state, 0–120 |
| `shielding_um` | Shielding (proximal screening) pigment position, µm |
| `tapetal_um` | Tapetal (reflecting) pigment position, µm |
| `facet` | Facet index across the eyeshine patch, 0 at the optic axis |
| `rhabdom` | Which rhabdom along that ray, 0 being the one it enters first |
| `pathlength_um` | Path length through that rhabdom, µm |

The path lengths are **raw geometry**: facet transmission is a flux factor and is
applied when the absorbed intensity is computed, not folded into the path length. A
ray that stopped propagating still contributes one row, with a path length of zero,
so every facet is accounted for.

The summary is accumulated as the rays are traced rather than by reading this file
back, so it is purely an output artefact and can be loaded directly:

```python
df = pd.read_csv("nephropsfl_pathlengths.csv")
df.groupby(["block", "facet"]).pathlength_um.sum()
```

### `genus_summary_res.csv` and `genus_summary_sen.csv`

Both are 11×11 matrices. **Rows vary the shielding pigment** from fully retracted
(row 0) to fully covering the rhabdom (row 10); **columns vary the tapetal pigment**
over the same range.

| File | Quantity | Units |
| --- | --- | --- |
| `summary_res` | Acceptance angle: FWHM of the point spread function | degrees |
| `summary_sen` | Incident light absorbed, area-weighted over the eyeshine patch | percent (0–100) |

A resolution cell reading `NaN` means the profile never falls to half its maximum,
so the acceptance angle is undefined for that pigment state.

The point spread function is the light arriving at each whole-rhabdom offset from the
optic axis, divided by the area of the annulus it is spread over. Facets are weighted
by their own source annulus, since the number of ommatidia at a given radius in the
eyeshine patch grows with that radius.

## Model notes

* **Blur circle extent** is the width of the blur circle in rhabdoms. 1 is a perfect
  point focus; the outermost facet is displaced by (extent − 1) rhabdoms. It may not
  exceed the number of facets across the eyeshine patch, since the light would
  otherwise have to fill rhabdom offsets that no facet reaches.
* **Critical angle.** Ray angles (`boa`) are measured from the rhabdom axis, so the
  angle at the wall normal is (90° − boa) and light is guided while
  `boa < 90° − asin(n_cytoplasm / n_rhabdom)`.
* **Absorption coefficient** is fixed at 0.01 µm⁻¹ in the Beer-Lambert absorbance
  `1 − exp(−kL)`. Reported values for crustacean rhabdoms span roughly
  0.0067–0.01 µm⁻¹.
* **Tapetal reflectance** is implicitly 1.0: the return path is added at full length
  with no loss term.
* **Parameter validation.** Parameter sets that cannot describe a physically
  realisable eye are rejected with a diagnostic and skipped, rather than being
  allowed to produce NaNs that silently disable the total-internal-reflection test.

### Known limitation

A trace terminates at the first reflection rather than following the reflected ray
onward. Because a ray that is *not* reflected continues leaking into adjacent
rhabdoms and absorbing there, extending the tapetum can occasionally shorten the
total absorbing path and so lower the reported sensitivity, by up to about 6
percentage points. A tapetal mirror can only add path length in reality, so this is a
limitation of the 1995 case structure rather than a property of the eye.

## Citation

If you use this program, please cite:

> Gaten, E., Moss, S., Johnson, M. 2013. The Reniform Reflecting Superposition Compound Eyes of Nephrops Norvegicus: Optics, Susceptibility to Light-Induced Damage, Electrophysiology and a Ray Tracing Model. In: M. L. Johnson and M. P. Johnson, ed(s). Advances in Marine Biology: The Ecology and Biology of Nephrops norvegicus. Oxford: Academic Press, 107:148.
