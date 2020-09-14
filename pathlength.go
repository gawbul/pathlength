package main

// Import required packages
import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Program input variables
var (
	inputFilename string
	showLicense   bool
	showCitation  bool
	showUsage     bool
	showVersion   bool
)

// Set program version
const (
	pathlengthVersion = "0.0.1"
)

// Init function
func init() {
	// Setup command line parameters
	flag.StringVar(&inputFilename, "filename", "", "Input filename in CSV format")
	flag.BoolVar(&showLicense, "license", false, "Display the license for the program")
	flag.BoolVar(&showCitation, "citation", false, "Display the citation for the program")
	flag.BoolVar(&showUsage, "usage", false, "Display program usage")
	flag.BoolVar(&showVersion, "version", false, "Display program version")
}

// Main function
func main() {
	// Parse input arguments
	flag.Parse()

	// Process input arguments
	if showUsage {
		PrintUsage(0)
	} else if showVersion {
		PrintVersion()
	} else if showLicense {
		PrintLicense()
	} else if showCitation {
		PrintCitation()
	} else if inputFilename == "" {
		PrintUsage(1)
	}

	// Open input file
	fmt.Println("Opening input file", inputFilename)
	csvFile, err := os.Open(inputFilename)
	if err != nil {
		log.Fatalln("Error opening input file", err)
	}

	// Parse input file
	fmt.Println("Parsing input parameters...")
	inputParameters, err := ParseCSV(csvFile)
	if err != nil {
		log.Fatalln("Error parsing input parameters", err)
	}
	fmt.Printf("Processed %d parameter sets.\n", len(*inputParameters))

	// Run the model
	fmt.Println("Running model...")
	for _, params := range *inputParameters {
		RunModel(&params)
	}
}
