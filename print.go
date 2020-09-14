package main

// Import required packages
import (
	"flag"
	"fmt"
	"os"
)

// PrintLicense function
// Outputs license details to the command line
func PrintLicense() {
	licenseText := `pathlength - calculates resolution and sensitivity in reflective superposition compound eyes.

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
along with this program.  If not, see <https://www.gnu.org/licenses/>`

	fmt.Println(licenseText)
	os.Exit(0)
}

// PrintCitation function
// Outputs citation information to the command line
func PrintCitation() {
	citationText := `Gaten, E., Moss, S., Johnson, M. 2013. The Reniform Reflecting Superposition Compound Eyes of Nephrops Norvegicus:
Optics, Susceptibility to Light-Induced Damage, Electrophysiology and a Ray Tracing Model. In: M. L. Johnson and M. P. Johnson, ed(s).
Advances in Marine Biology: The Ecology and Biology of Nephrops norvegicus. Oxford: Academic Press, 107:148.`

	fmt.Println(citationText)
	os.Exit(0)
}

// PrintUsage function
// Outputs usage information to the command line
func PrintUsage(exitCode int) {
	fmt.Println("Usage:")
	flag.PrintDefaults()
	os.Exit(exitCode)
}

// PrintVersion function
// Outputs the version of the program to the command line
func PrintVersion() {
	fmt.Println("pathlength version", pathlengthVersion)
	os.Exit(0)
}
