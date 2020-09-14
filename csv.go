package main

// Import required packages
import (
	"encoding/csv"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Parameters struct
// Used for parsing parameters from input file
type Parameters struct {
	Species                  string
	RhabdomLength            float64
	RhabdomWidth             float64
	EyeDiameter              float64
	FacetWidth               float64
	ApertureDiameter         float64
	CytoplasmRefractiveIndex float64
	RhabdomRefractiveIndex   float64
	BlurCircleExtent         float64
	ProximalRhabdomAngle     float64
}

// ParseCSV function
// Takes a pointer to an os.File object as a parameter
// Returns a pointer to an array of Parameters or an error object
func ParseCSV(csvFile *os.File) (*[]Parameters, error) {
	// Load CSV file into two dimensional string array
	rows, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}

	// Check species names are unique
	speciesNames := make([]string, len(rows))
	speciesNameCount := make(map[string]int)
	for i := range rows {
		speciesNames[i] = string(rows[i][0])
		speciesNameCount[speciesNames[i]]++
		if speciesNameCount[speciesNames[i]] > 1 {
			return nil, errors.New("duplicate species name found")
		}
	}

	// Define inputParameters array of type Parameters
	inputParameters := make([]Parameters, len(rows))

	// Define regex for species name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return nil, err
	}

	// Iterate over rows of array and parse values
	for i := range rows {
		if len(rows[i]) != 10 {
			return nil, errors.New("incorrect number of parameters, 10 expected")
		}
		spVal := reg.ReplaceAllString(strings.ToLower(string(rows[i][0])), "")
		rlVal, err := strconv.ParseFloat(string(rows[i][1]), 64)
		if err != nil {
			return nil, err
		}
		rwVal, err := strconv.ParseFloat(string(rows[i][2]), 64)
		if err != nil {
			return nil, err
		}
		edVal, err := strconv.ParseFloat(string(rows[i][3]), 64)
		if err != nil {
			return nil, err
		}
		fwVal, err := strconv.ParseFloat(string(rows[i][4]), 64)
		if err != nil {
			return nil, err
		}
		adVal, err := strconv.ParseFloat(string(rows[i][5]), 64)
		if err != nil {
			return nil, err
		}
		criVal, err := strconv.ParseFloat(string(rows[i][6]), 64)
		if err != nil {
			return nil, err
		}
		rriVal, err := strconv.ParseFloat(string(rows[i][7]), 64)
		if err != nil {
			return nil, err
		}
		bceVal, err := strconv.ParseFloat(string(rows[i][8]), 64)
		if err != nil {
			return nil, err
		}
		praVal, err := strconv.ParseFloat(string(rows[i][9]), 64)
		if err != nil {
			return nil, err
		}

		// Assign Paramaters object to index position in inputParameters array
		inputParameters[i] = Parameters{
			Species:                  spVal,
			RhabdomLength:            rlVal,
			RhabdomWidth:             rwVal,
			EyeDiameter:              edVal,
			FacetWidth:               fwVal,
			ApertureDiameter:         adVal,
			CytoplasmRefractiveIndex: criVal,
			RhabdomRefractiveIndex:   rriVal,
			BlurCircleExtent:         bceVal,
			ProximalRhabdomAngle:     praVal,
		}
	}
	return &inputParameters, nil
}

// OutputCSV function
// Takes a pointer to an os.File object and a 2 dimensional float64 array as a parameter
// Returns a boolean or an error object
func OutputCSV(csvfile *os.File, outputVals *[][]float64) (bool, error) {

	return true, nil
}
