# Superposition Eye Pathlength Program

PathLength implements a ray tracing model to calculate the resolution and sensitivity of reflective superposition compound eyes.

Original QBASIC version by Dr Magnus L Johnson and Genevre Parker, 1995

Golang rewrite by Dr Stephen P Moss, 2020

Author: Dr Stephen P Moss

Website: [https://www.gawbul.io](https://www.gawbul.io)

Email: gawbul@gmail.com

# Install Go compiler

## macOS
```bash
# Install brew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"

# Install Golang
brew install go@1.15
```

## Linux (Debian/Ubuntu)
```bash
# Download tarball
wget -P /tmp https://dl.google.com/go/go1.15.1.linux-amd64.tar.gz

# Extract tarball
sudo tar -zxvf /tmp/go1.15.1.linux-amd64.tar.gz -C /usr/local

# Setup environment
echo "export GOROOT=/usr/local/go" >> ~/.profile
echo "export GOPATH=$HOME/go" >> ~/.profile
echo "export PATH=$GOPATH/bin:$GOROOT/bin:$PATH" >> ~/.profile

# Source environment
source ~/.profile
```

# Install dependencies

## macOS
```bash
# Install git
brew install git

# Install Go packages
go get github.com/stretchr/testify
```

## Linux (Debian/Ubuntu)
```bash
# Install git
sudo apt-get update
sudo apt-get install git

# Install Go packages
go get github.com/stretchr/testify
```

# Checkout source code

```bash
# Create and change to projects directory
mkdir -p ~/projects
cd ~/projects

# Clone the GitHub repository
git clone git@github.com:gawbul/pathlength.git
```
 
# Usage

## Run the program from source
```bash
cd ~/projects/pathlength
go run .
```
Outputs:
```
Usage:
  -citation
        Display the citation for the program
  -filename string
        Input filename in CSV format
  -license
        Display the license for the program
  -usage
        Display program usage
exit status 1
```

## Run the test suite
```bash
cd ~/projects/pathlength
go test -v
```
Outputs:
```
=== RUN   TestParseParameters
--- PASS: TestParseParameters (0.00s)
PASS
ok      _/Users/stephenmoss/Dropbox/Code/pathlength     0.320s
```

## Build the program
```bash
cd ~/projects/pathlength
go build
chmod +x pathlength
```

## Display program usage
```bash
./pathlength -usage
```
Outputs:
```
Usage:
  -citation
        Display the citation for the program
  -filename string
        Input filename in CSV format
  -license
        Display the license for the program
  -usage
        Display program usage
  -version
        Display program version
```
*Also displays if you don't pass in any arguments, as it expects a filename as input.*

## Display citation information
```bash
./pathlength -citation
```
Outputs:
```
Gaten, E., Moss, S., Johnson, M. 2013. The Reniform Reflecting Superposition Compound Eyes of Nephrops Norvegicus:
Optics, Susceptibility to Light-Induced Damage, Electrophysiology and a Ray Tracing Model. In: M. L. Johnson and M. P. Johnson, ed(s).
Advances in Marine Biology: The Ecology and Biology of Nephrops norvegicus. Oxford: Academic Press, 107:148.
```

## Display license information
```bash
./pathlength -license
```
Outputs:
```
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

## Display program version
```bash
./pathlength -version
```
Outputs:
```
pathlength version 0.0.1
```

# Run the program
```bash
./pathlength -filename <input_filename.csv>
```
Outputs:
```

```

# Required parameters
A CSV format file is required as input to the program. You can provide multiple lines for separate runs of the model. The format should be as follows:
```
nephropsfl,180,25,7800,50,3200,1.34,1.37,18,0
nephropspl,180,25,7800,50,3200,1.34,1.37,18,12.5
nephropsfa,180,25,6760,50,3060,1.34,1.37,10,0
nephropspa,180,25,6760,50,3060,1.34,1.37,10,12.5
```

Each row is comprised of the following fields:
```
genus	=	A prefix for the output filenames e.g. organism genus name (lowercase alphanumeric only)
180 	=	Rhabdom Length
25 		=	Rhabdom Width
7800 	=	Eye Diameter
50 		=	Facet Width
3200	=	Aperture Diameter
1.34	=	Cytoplasm Refractive Index
1.37	=	Rhabdom Refractive Index
18		=	Blur Circle Extent
0		=	Proximal Rhabdom Angle (used to create pointy-ended rhabdoms)
```
*NB: The genus name is NOT case sensitive. It is always converted to lowercase and should be unique to avoid filename conflicts.*

# Output files
Three output files are created:

* genus_pathlengths.csv
* genus_resolution.csv
* genus_sensitivity.csv

The first file contains multiple rows for each facet with the various combinations of tapetal and shielding pigment lengths and then multiple rows with the pathlengths 

This outputs two files (where genus is the name you give when setting up the object):

	genus_output_one.txt	=	Each record is separated by 999 in the text and contains the length of the reflective tapetum and shielding pigment initially, followed by the path length values for each rhabdom the light passes through, starting at the axial rhabdom.
	genus_output_two.txt	=	**Description needed**


This outputs three files (where genus is the name you give when setting up the object):

	genus_summary_one.txt	= **Description needed**
	genus_summary_res.txt	= Resolution output
	genus_summary_sen.txt	= Sensitivity output



# Citation

If you use this program, please cite:

Gaten, E., Moss, S., Johnson, M. 2013. The Reniform Reflecting Superposition Compound Eyes of Nephrops Norvegicus: Optics, Susceptibility to Light-Induced Damage, Electrophysiology and a Ray Tracing Model. In: M. L. Johnson and M. P. Johnson, ed(s). Advances in Marine Biology: The Ecology and Biology of Nephrops norvegicus. Oxford: Academic Press, 107:148.
