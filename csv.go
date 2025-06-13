// FILE: csv.go
// This file contains functions for reading and writing CSV data.

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// calculateRessens calculates resolution and sensitivity from the generated pathlengths file.
// This is a port of the `summarise_data` function from the Python script.
func (m *Model) calculateRessens() {
	p := m.Params
	fmt.Println("INFO: Calculating resolution and sensitivity...")

	// Open output files
	resFile, err := os.Create(fmt.Sprintf("%s_summary_res.csv", p.SpeciesName))
	if err != nil {
		log.Fatalf("Error creating resolution file: %v", err)
	}
	defer resFile.Close()
	resWriter := bufio.NewWriter(resFile)
	defer resWriter.Flush()

	sensFile, err := os.Create(fmt.Sprintf("%s_summary_sen.csv", p.SpeciesName))
	if err != nil {
		log.Fatalf("Error creating sensitivity file: %v", err)
	}
	defer sensFile.Close()
	sensWriter := bufio.NewWriter(sensFile)
	defer sensWriter.Flush()

	// Open the pathlengths file for reading
	pathlengthsFile, err := os.Open(fmt.Sprintf("%s_pathlengths.csv", p.SpeciesName))
	if err != nil {
		log.Fatalf("Error opening pathlengths file for reading: %v", err)
	}
	defer pathlengthsFile.Close()

	// --- State variables for summary calculation ---
	rhabdoms := make([]float64, 21)
	matrixSens := []string{}
	matrixRes := []string{}
	facet := 0.0
	arem := 0.0
	cc, dd := 0, 0

	scanner := bufio.NewScanner(pathlengthsFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.Contains(line, "998") {
			parts := strings.Split(line, ",")
			rhabdom := 0
			tot := 0.0
			area := math.Pi * math.Pow(facet+0.5, 2)
			inci := math.Pi * math.Pow(facet-0.5, 2)
			if facet == 0 {
				inci = 0
			}
			torus := area - inci
			if area > arem {
				arem = area
			}

			for _, part := range parts {
				if part == "998" {
					break
				}
				pathlength, _ := strconv.ParseFloat(part, 64)
				var absorbance, bx float64
				if pathlength > 0 {
					absorbance = 1 - math.Exp(-0.01*pathlength)
				} else {
					absorbance = 0
				}
				if rhabdom == 0 && absorbance > 0 {
					bx = 100 * absorbance
				} else if rhabdom > 0 && absorbance > 0 {
					bx = 100 * ((1 - tot) * absorbance)
				}
				if absorbance == 0 {
					bx = 0
				}
				tot += bx / 100.0
				bx *= torus
				if rhabdom < len(rhabdoms) {
					rhabdoms[rhabdom] += bx
				}
				rhabdom++
			}
			facet++
		} else if line == "999" {
			// End of block, summarize
			sens := 0.0
			for _, r := range rhabdoms {
				sens += r
			}
			halfwayPoint := rhabdoms[0] / 2.0
			opticAxis := 0.0
			xz, yy := rhabdoms[0], rhabdoms[1]

			for i := 1; i < 12; i++ {
				if halfwayPoint < rhabdoms[i] {
					xz = rhabdoms[i]
					if i+1 < len(rhabdoms) {
						yy = rhabdoms[i+1]
					}
					opticAxis = m.OmmatidialAngle * float64(i)
					break
				}
			}
			diff := xz - yy
			hwp := xz - halfwayPoint
			frac := hwp / (diff + 0.1) // prevent div by zero
			oab := frac * m.OmmatidialAngle
			res := oab + opticAxis

			if cc == 0 && dd > 0 {
				fmt.Fprintln(sensWriter, strings.Join(matrixSens, ","))
				fmt.Fprintln(resWriter, strings.Join(matrixRes, ","))
				matrixSens = []string{}
				matrixRes = []string{}
			}

			if arem > 0 {
				matrixSens = append(matrixSens, fmt.Sprintf("%d", int(sens/arem)))
			} else {
				matrixSens = append(matrixSens, "0")
			}
			matrixRes = append(matrixRes, fmt.Sprintf("%d", int(res*200)))

			cc++
			if cc == 11 {
				dd++
				cc = 0
			}
			// Reset for next block
			for i := range rhabdoms {
				rhabdoms[i] = 0
			}
			facet = 0
		}
	}
	// Write the final line of data
	if len(matrixSens) > 0 {
		fmt.Fprintln(sensWriter, strings.Join(matrixSens, ","))
		fmt.Fprintln(resWriter, strings.Join(matrixRes, ","))
	}
}

// parseInputParameters reads a parameter file using Go's standard CSV reader.
// It returns a slice of parsed parameters and an error if parsing fails.
func parseInputParameters(filename string) ([]Parameters, error) {
	var paramsList []Parameters // A slice to hold multiple parameter sets

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open parameter file %s: %w", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		// Check for a specific parsing error (like wrong number of fields)
		if parseErr, ok := err.(*csv.ParseError); ok && parseErr.Err == csv.ErrFieldCount {
			log.Printf("Skipping malformed record on line %d: %v", parseErr.Line, err)
			continue // Skip to the next record
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		if len(record) != 10 {
			log.Printf("Skipping malformed record (expected 10 fields, got %d): %v", len(record), record)
			continue
		}

		var params Parameters
		// (sn, rl, rw, ed, fw, ad, cri, rri, bce, pra)
		params.SpeciesName = strings.TrimSpace(record[0])
		params.RhabdomLength, _ = strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		params.RhabdomWidth, _ = strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		params.EyeDiameter, _ = strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
		params.FacetWidth, _ = strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
		params.ApertureDiameter, _ = strconv.ParseFloat(strings.TrimSpace(record[5]), 64)
		params.CytoplasmRefractiveIndex, _ = strconv.ParseFloat(strings.TrimSpace(record[6]), 64)
		params.RhabdomRefractiveIndex, _ = strconv.ParseFloat(strings.TrimSpace(record[7]), 64)
		bce, _ := strconv.ParseFloat(strings.TrimSpace(record[8]), 64)
		params.BlurCircleExtent = math.Max(1.0, bce) // Ensure blur circle is at least 1
		params.ProximalRhabdomAngle, _ = strconv.ParseFloat(strings.TrimSpace(record[9]), 64)

		paramsList = append(paramsList, params)
	}

	if len(paramsList) == 0 {
		return nil, fmt.Errorf("no valid parameter data found in %s", filename)
	}

	return paramsList, nil
}
