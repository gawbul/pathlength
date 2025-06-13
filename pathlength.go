// FILE: pathlength.go
// This file is the main entry point for the application.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const version = "0.6.0"

func main() {
	// --- Command Line Argument Parsing ---
	paramFile := flag.String("f", "", "Path to a parameter file (CSV format). (Required)")
	showCitation := flag.Bool("c", false, "Show the program citation.")
	showHelp := flag.Bool("h", false, "Show this help message.")
	showLicense := flag.Bool("l", false, "Show the program license.")
	showVersion := flag.Bool("v", false, "Show program version.")
	flag.Parse()

	if *showLicense {
		fmt.Println(`pathlength - calculates resolution and sensitivity in reflective superposition compound eyes.

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
along with this program.  If not, see <https://www.gnu.org/licenses/>`)
		os.Exit(0)
	}

	if *showCitation {
		fmt.Println(`Gaten, E., Moss, S., Johnson, M. 2013. The Reniform Reflecting Superposition Compound Eyes of Nephrops Norvegicus:
Optics, Susceptibility to Light-Induced Damage, Electrophysiology and a Ray Tracing Model. In: M. L. Johnson and M. P. Johnson, ed(s).
Advances in Marine Biology: The Ecology and Biology of Nephrops norvegicus. Oxford: Academic Press, 107:148.`)
		os.Exit(0)
	}

	if *showHelp {
		fmt.Printf("Usage: %s -f filename [-h] [-v]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("%s version %s\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}

	// --- Initialisation ---
	// Exit if no parameter file is provided.
	if *paramFile == "" {
		fmt.Printf("Usage: %s -f filename [-h] [-v]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		log.Fatal("Error: No parameter file supplied. Use the -f flag to specify a file.")
	}

	fmt.Printf("Parsing input parameters from %s...\n", *paramFile)
	paramsList, err := parseInputParameters(*paramFile)
	if err != nil {
		log.Fatalf("Error parsing parameter file: %v", err)
	}

	// --- Loop over each parameter set and run the model ---
	for _, params := range paramsList {
		model := NewModel(params)

		fmt.Printf("--- Running simulation for %s ---\n", model.Params.SpeciesName)

		// --- Run Simulation & Calculate Results ---
		fmt.Printf("Calculating pathlengths for %s...\n", model.Params.SpeciesName)
		model.runModel()

		model.calculateRessens()

		fmt.Printf("--- Finished simulation for %s ---\n\n", model.Params.SpeciesName)
	}

	fmt.Println("All simulations complete.")
}
