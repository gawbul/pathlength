# Superposition Eye Pathlength Program

PathLength implements a ray tracing model to calculate the resolution and sensitivity of reflective superposition compound eyes.

Original QBASIC version by Dr Magnus L Johnson and Genevre Parker, 1995

Golang rewrite by Dr Stephen P Moss, 2020

Author: Dr Stephen P Moss

Website: [https://www.gawbul.io](https://www.gawbul.io)

Email: gawbul@gmail.com

## Install Go compiler

### macOS

```bash
# Install brew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"

# Install Golang
brew install go@1.24
```

### Linux (Debian/Ubuntu)

```bash
# Download tarball
wget -P /tmp https://dl.google.com/go/go1.24.4.linux-amd64.tar.gz

# Extract tarball
sudo tar -zxvf /tmp/go1.24.4.linux-amd64.tar.gz -C /usr/local

# Setup environment
echo "export GOROOT=/usr/local/go" >> ~/.profile
echo "export GOPATH=$HOME/go" >> ~/.profile
echo "export PATH=$GOPATH/bin:$GOROOT/bin:$PATH" >> ~/.profile

# Source environment
source ~/.profile
```

## Install dependencies

### macOS

```bash
# Install git
brew install git

# Install Go packages
go get github.com/stretchr/testify
```

### Linux (Debian/Ubuntu)

```bash
# Install git
sudo apt-get update
sudo apt-get install git

# Install Go packages
go get github.com/stretchr/testify
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
Usage: pathlength -f filename [-h] [-v]
  -c    Show the program citation.
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
2025/06/13 15:05:08 Skipping malformed record (expected 10 fields, got 3): [test_species 100 10]
--- PASS: TestParseInputParameters (0.00s)
    --- PASS: TestParseInputParameters/ValidFile (0.00s)
    --- PASS: TestParseInputParameters/NonExistentFile (0.00s)
    --- PASS: TestParseInputParameters/MalformedFile (0.00s)
=== RUN   TestCalculateRessens
INFO: Calculating resolution and sensitivity...
--- PASS: TestCalculateRessens (0.00s)
PASS
ok      pathlength      0.271s
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
Usage: pathlength -f filename [-h] [-v]
  -c    Show the program citation.
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
Initialising model parameters for acanthephyra...
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
P: 12.70, T: 50.80
P: 12.70, T: 63.50
P: 12.70, T: 76.20
P: 12.70, T: 88.90
P: 12.70, T: 101.60
P: 12.70, T: 114.30
P: 12.70, T: 127.00
P: 25.40, T: 0.00
P: 25.40, T: 12.70
P: 25.40, T: 25.40
P: 25.40, T: 38.10
P: 25.40, T: 50.80
P: 25.40, T: 63.50
P: 25.40, T: 76.20
P: 25.40, T: 88.90
P: 25.40, T: 101.60
P: 25.40, T: 114.30
P: 25.40, T: 127.00
P: 38.10, T: 0.00
P: 38.10, T: 12.70
P: 38.10, T: 25.40
P: 38.10, T: 38.10
P: 38.10, T: 50.80
P: 38.10, T: 63.50
P: 38.10, T: 76.20
P: 38.10, T: 88.90
P: 38.10, T: 101.60
P: 38.10, T: 114.30
P: 38.10, T: 127.00
P: 50.80, T: 0.00
P: 50.80, T: 12.70
P: 50.80, T: 25.40
P: 50.80, T: 38.10
P: 50.80, T: 50.80
P: 50.80, T: 63.50
P: 50.80, T: 76.20
P: 50.80, T: 88.90
P: 50.80, T: 101.60
P: 50.80, T: 114.30
P: 50.80, T: 127.00
P: 63.50, T: 0.00
P: 63.50, T: 12.70
P: 63.50, T: 25.40
P: 63.50, T: 38.10
P: 63.50, T: 50.80
P: 63.50, T: 63.50
P: 63.50, T: 76.20
P: 63.50, T: 88.90
P: 63.50, T: 101.60
P: 63.50, T: 114.30
P: 63.50, T: 127.00
P: 76.20, T: 0.00
P: 76.20, T: 12.70
P: 76.20, T: 25.40
P: 76.20, T: 38.10
P: 76.20, T: 50.80
P: 76.20, T: 63.50
P: 76.20, T: 76.20
P: 76.20, T: 88.90
P: 76.20, T: 101.60
P: 76.20, T: 114.30
P: 76.20, T: 127.00
P: 88.90, T: 0.00
P: 88.90, T: 12.70
P: 88.90, T: 25.40
P: 88.90, T: 38.10
P: 88.90, T: 50.80
P: 88.90, T: 63.50
P: 88.90, T: 76.20
P: 88.90, T: 88.90
P: 88.90, T: 101.60
P: 88.90, T: 114.30
P: 88.90, T: 127.00
P: 101.60, T: 0.00
P: 101.60, T: 12.70
P: 101.60, T: 25.40
P: 101.60, T: 38.10
P: 101.60, T: 50.80
P: 101.60, T: 63.50
P: 101.60, T: 76.20
P: 101.60, T: 88.90
P: 101.60, T: 101.60
P: 101.60, T: 114.30
P: 101.60, T: 127.00
P: 114.30, T: 0.00
P: 114.30, T: 12.70
P: 114.30, T: 25.40
P: 114.30, T: 38.10
P: 114.30, T: 50.80
P: 114.30, T: 63.50
P: 114.30, T: 76.20
P: 114.30, T: 88.90
P: 114.30, T: 101.60
P: 114.30, T: 114.30
P: 114.30, T: 127.00
P: 127.00, T: 0.00
P: 127.00, T: 12.70
P: 127.00, T: 25.40
P: 127.00, T: 38.10
P: 127.00, T: 50.80
P: 127.00, T: 63.50
P: 127.00, T: 76.20
P: 127.00, T: 88.90
P: 127.00, T: 101.60
P: 127.00, T: 114.30
P: 127.00, T: 127.00
INFO: Calculating resolution and sensitivity...
Done.
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

Four output files are created:

* genus_debug.csv
* genus_pathlengths.csv
* genus_resolution.csv
* genus_sensitivity.csv

The pathlengths file contains multiple rows for each facet with the various combinations of tapetal and shielding pigment lengths in the adjacent columns and then multiple columns with the pathlengths.

The resolution file contains the calculated resolution values.

The sensitivity file contains the calculated sensitivity values.

## Citation

If you use this program, please cite:

> Gaten, E., Moss, S., Johnson, M. 2013. The Reniform Reflecting Superposition Compound Eyes of Nephrops Norvegicus: Optics, Susceptibility to Light-Induced Damage, Electrophysiology and a Ray Tracing Model. In: M. L. Johnson and M. P. Johnson, ed(s). Advances in Marine Biology: The Ecology and Biology of Nephrops norvegicus. Oxford: Academic Press, 107:148.
