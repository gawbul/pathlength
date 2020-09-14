package main

// Import required packages
import (
	"fmt"
	"math"
)

// Eye struct
// Defines the parameters of an eye object
type Eye struct {
	Circumference      float64
	Diameter           float64
	Radius             float64
	DistanceToAperture float64
	RefractiveIndex    float64
	OmmatidialAngle    float64
	AngleAtCenter      float64
	Aperture           Aperture
	Rhabdoms           []Rhabdom
}

// Aperture struct
// Defines the parameters of an aperture object
type Aperture struct {
	Diameter float64
	Radius   float64
	Curve    float64
}

// Rhabdom struct
// Defines the parameters of a rhabdom object
type Rhabdom struct {
	Length          float64
	Diameter        float64
	Radius          float64
	RefractiveIndex float64
}

// Define variables
var (
	tapetalPigment          float64
	shieldingPigment        float64
	eyeCircumference        float64
	eyeRadius               float64
	apertureRadius          float64
	distanceToAperture      float64
	angleAtCenter           float64
	apertureCurve           float64
	ommatidialAngle         float64
	facetAdjustment         int
	numberOfFacets          int
	rhabdomRadius           float64
	interOmmatidialAngle    float64
	exitOmmatidialAngle     float64
	increaseAcceptanceAngle bool
	snellsLaw               float64
	criticalAngle           float64
	rhabdomSlope            float64
)

// DegreesToRadians function
// Takes a pointer to a float64 and returns a pointer to a float64
func DegreesToRadians(val float64) float64 {
	return val * (math.Pi / 180)
}

// RadiansToDegrees function
// Takes a pointer to a float64 and returns a pointer to a float64
func RadiansToDegrees(val float64) float64 {
	return val * (180 / math.Pi)
}

// RunModel function
// Takes pointer to a Parameters array as a parameter
// Returns a pointer to a 2 dimensional float64 array or error object
func RunModel(parameters *Parameters) (*[][]float64, error) {
	fmt.Println(parameters)
	// Initialise variables
	tapetalPigment = 0.0
	shieldingPigment = 0.0
	eyeCircumference = math.Pi * parameters.EyeDiameter
	eyeRadius = parameters.EyeDiameter / 2
	apertureRadius = parameters.ApertureDiameter / 2
	distanceToAperture = math.Sqrt((eyeRadius - eyeRadius) - (apertureRadius - apertureRadius))
	angleAtCenter = RadiansToDegrees(math.Atan(apertureRadius / distanceToAperture))
	apertureCurve = eyeCircumference * (angleAtCenter / 360)
	numberOfFacets = int(math.Round(parameters.ApertureDiameter / parameters.FacetWidth))
	interOmmatidialAngle = 0.0
	exitOmmatidialAngle = 0.0
	snellsLaw = math.Asin(parameters.CytoplasmRefractiveIndex / parameters.RhabdomRefractiveIndex)
	criticalAngle = RadiansToDegrees(snellsLaw)
	rhabdomSlope = math.Sqrt(math.Pow(parameters.RhabdomLength, 2) + math.Pow(rhabdomRadius, 2))
	increaseAcceptanceAngle = false

	// Initialise rhabdoms
	rhabdoms := make([]Rhabdom, numberOfFacets)
	for i := range rhabdoms {
		rhabdoms[i] = Rhabdom{
			Length:          parameters.RhabdomLength,
			Diameter:        parameters.RhabdomWidth,
			Radius:          parameters.RhabdomWidth / 2,
			RefractiveIndex: parameters.RhabdomRefractiveIndex,
		}
	}

	// Create new Eye object
	eye := Eye{
		Circumference:      eyeCircumference,
		Diameter:           parameters.EyeDiameter,
		Radius:             eyeRadius,
		DistanceToAperture: distanceToAperture,
		RefractiveIndex:    parameters.CytoplasmRefractiveIndex,
		AngleAtCenter:      angleAtCenter,
		OmmatidialAngle:    ommatidialAngle,
		Aperture: Aperture{
			Diameter: parameters.ApertureDiameter,
			Radius:   apertureRadius,
			Curve:    apertureCurve,
		},
		Rhabdoms: rhabdoms,
	}

	// Iterate over each facet
	for i := 0; i < numberOfFacets; i++ {
		// Iterate through tapetal pigment lengths
		for tapetalPigment = 0.0; tapetalPigment <= eye.Rhabdoms[i].Length+1; tapetalPigment += eye.Rhabdoms[i].Length / 10 {
			// Iterate through shielding pigment lengths
			for shieldingPigment = 0.0; shieldingPigment <= eye.Rhabdoms[i].Length+1; shieldingPigment += eye.Rhabdoms[i].Length / 10 {
				fmt.Printf("%d %0.2f %0.2f\n", i, tapetalPigment, shieldingPigment)
			}
		}
	}

	fmt.Println(eye)
	return nil, nil
}
