/*
Filename: pathlength.c
Author: Steve Moss (gawbul@gmail.com)
Copyright: (C) Steve Moss, 2014

Models pencils of light travelling through a reflective superposition compound eye.
Their pathlength through the rhabdoms is used to determine resolution and sensitivity.

1. First we assume the light starts at a single point outside the eye
2. Light travels as a single pencil of light into the eye
3. Movement of the light is implicit in the flow of the program
4. The light travels vertically through the facet of the eye
5. The light is reflected 
6. 

We move to the next facet and repeat steps 4 onwards.
*/

/* include required header files */
#include <stdlib.h>
#include <stdio.h>
#include <math.h>
#include <string.h>

/********************************/
/* ONLY CHANGE THESE PARAMETERS */
/********************************/

/* declare global constant eye parameters */
const char* SPECIES_NAME = "astacodes";				/* species name */
const double RHABDOM_LENGTH = 84.0;					/* length of the rhabdom */
const double RHABDOM_WIDTH = 16.0;					/* width of the rhabdom */
const double EYE_DIAMETER = 890.0;					/* diameter of the eye */
const double FACET_WIDTH = 32.0;					/* width of the ommatidia's facet */
const double APERTURE_DIAMETER = 445.0;				/* diameter of the aperture of the eye */
const double CYTOPLASM_REFRACTIVE_INDEX = 1.34;		/* refractive index of the eye's cytoplasm */
const double RHABDOM_REFRACTIVE_INDEX = 1.37;		/* refractive index within the eye */
const double BLUR_CIRCLE_EXTENT = 18.0;				/* diameter of the blur circle */
const double PROXIMAL_RHABDOM_ANGLE = 0;			/* allow for a different angle at the tips of the rhabdoms */

/*************************/
/* NO CHANGES BELOW HERE */
/*************************/

/* define constant values */
#ifndef M_PI
	#define M_PI 3.14159265358979323846			/* store value of pi */
#endif
const double RAD_TO_DEG_CONV = M_PI / 180.0;	/* store conversion from radians to degrees */

/* declare model functions */
void initialise(void);
void run_model(void);
void calculate_ressens(void);

/* declare global model parameters */
double tapetal_pigment;				/* store the extent of the tapetal pigment */
double shielding_pigment;			/* store the extent of the shielding pigment */
double circumference_of_eye;		/* store the circumference of the eye */
double aperture_radius;				/* store the aperture radius */
double eye_radius;					/* store the eye radius */
double distance_to_aperture;		/* store distance from center of the eye to the aperture */
double angle_at_center;				/* store angle at center of the eye */
double aperture_diameter;			/* store the aperture diameter */
double ommatidial_angle;			/* store the ommatidial_angle */
/* ####### vvvvvvv ####### */
double facet_adjustment;			/* ### WHAT IS THIS? FA = 1 IN ORIGINAL CODE ### */
/* ####### ^^^^^^^ ####### */
int number_of_facets;				/* store the number of facets in the aperture */
double rhabdom_radius;				/* store the rhabdom radius */
/* ####### vvvvvvv ####### */
double old_rhabdom_length;			/* store the old rhabdom length - DEPRECATE? */
double max_rhabdom_length;			/* store the max rhabdom length - DEPRECATE? */
/* ####### ^^^^^^^ ####### */
double inter_ommatidial_angle;		/* store the interommatidial angle */
int current_facet;					/* store the current facet number */
double snells_law;					/* store Snell's Law formula */
double critical_angle;				/* store critical angle */
double rhabdom_slope;				/* store rhabdom slope */
int increase_acceptance_angle;		/* store boolean value for increasing acceptance angle */
double pathlength;					/* store pathlength value */
double exit_ommatidial_angle;		/* store exit ommatidial angle */

/* store values in 2d array for writing */
double output_data[(int) (APERTURE_DIAMETER / FACET_WIDTH)][60];

/* main subroutine */
int main() {
	/* initialise parameters */
	printf("Initialising model parameters for %s...\n", SPECIES_NAME);
	initialise();
	
	/* run model */
	printf("Calculating pathlengths for %s...\n", SPECIES_NAME);
	run_model();
	
	/* calculate resolution and sensitivity */
	printf("Calculating resolution and sensitivity for %s...\n", SPECIES_NAME);
	calculate_ressens();

	/* let user know we've finished */
	printf("Done.\n");
}

/* initialise the model parameters */
void initialise(void) {
	/* initial calculations */
	tapetal_pigment = 0.0;		/* Initialise tapetal pigment to zero */
	shielding_pigment = 0.0;	/* Initialise shielding pigment to zero */
	circumference_of_eye = M_PI * EYE_DIAMETER;														/* calculate the circumference of the eye */
	aperture_radius =  APERTURE_DIAMETER / 2.0;														/* calculate the aperture radius */
	eye_radius = EYE_DIAMETER / 2.0;																/* calculate the eye radius */
	distance_to_aperture = sqrt((eye_radius * eye_radius) - (aperture_radius * aperture_radius));	/* use pythagorus theorum to work out distance to aperture */
	angle_at_center = atan(aperture_radius / distance_to_aperture) / RAD_TO_DEG_CONV;				/* and determine angle at center */
	aperture_diameter = circumference_of_eye * (angle_at_center / 360.0);							/* calculate the curve of half the aperture */
	/* can also do this by without converting radians to degrees first */
	/* angle_at_center = atan(aperture_radius / distance_to_aperture);
	aperture_diameter = circumference_of_eye * (angle_at_center / (2 * M_PI)); */
	ommatidial_angle = (FACET_WIDTH / circumference_of_eye) * 360.0;								/* calculate ommatidial_angle */

	/* number of facets in eyeshine patch axis */
	facet_adjustment = 1.0;									/* set facet adjustment */
	number_of_facets = aperture_diameter / FACET_WIDTH;		/* calculate the number of facets */
	rhabdom_radius = RHABDOM_WIDTH / 2.0;					/* calculate rhabdom radius */
	old_rhabdom_length = RHABDOM_LENGTH;					/* store old rhabdom length */
	max_rhabdom_length = RHABDOM_LENGTH;					/* store max rhabdom length */
	inter_ommatidial_angle = 0.0;	/* set the inter ommatidial angle */
	current_facet = 1.0;			/* set the current facet */

	/* calculate the angle of total internal reflection */
	snells_law = asin(CYTOPLASM_REFRACTIVE_INDEX / RHABDOM_REFRACTIVE_INDEX);						/* calculate value using Snell's Law formula */
	critical_angle = 90 - snells_law;																/* calculate critical angle below which light is totally internally reflected */
	rhabdom_slope = sqrt((RHABDOM_LENGTH * RHABDOM_LENGTH) + (rhabdom_radius * rhabdom_radius));	/* calculate the slope/hypotenuse of the rhambdom */
	/* can also do this with pythagorus theorum */
	/* adjacent_side = sqrt((RHABDOM_REFRACTIVE_INDEX * RHABDOM_REFRACTIVE_INDEX) - (CYTOPLASM_REFRACTIVE_INDEX * CYTOPLASM_REFRACTIVE_INDEX));
	critical_below = atan(CYTOPLASM_REFRACTIVE_INDEX / adjacent_side);
	critical_angle = 90 - (critical_below / RAD_TO_DEG_CONV); */
	increase_acceptance_angle = 0; /* set increase acceptance angle as false */
}

/* run the model */
void run_model(void) {
	/* open files for output */
	char* filename;
	asprintf(&filename, "%s_pathlengths.csv", SPECIES_NAME);
	FILE *output = fopen(filename, "w");
	if (output == NULL)
	{
	    printf("Error opening file!\n");
	    exit(1);
	}

	/* iterate through tapetal pigment lengths until they are the full length of the rhabdom */
	/* need to add 1 to rhabdom length due to problems with floating point precision */
	while (tapetal_pigment <= RHABDOM_LENGTH + 1) {
		/* iterate through shielding pigment lengths until they are the full length of the rhabdom */
		/* need to add 1 to rhabdom length due to problems with floating point precision */
		while (shielding_pigment <= RHABDOM_LENGTH + 1) {
			/* display shielding and tapetal pigment values */
			printf("T: %0.2f\n", tapetal_pigment);
			printf("P: %0.2f\n", shielding_pigment);
			/* save values to output file */
			fprintf(output, "%0.2f\n", tapetal_pigment);
			fprintf(output, "%0.2f\n", shielding_pigment);

			/* iterate over each facet in turn */
			while (current_facet <= number_of_facets) {
				/* calculate pathlength through rhabdoms for each facet */
				/* calculate prox-dist length of first pass */
				if (exit_ommatidial_angle > critical_angle && increase_acceptance_angle == 0) {
					exit_ommatidial_angle = exit_ommatidial_angle - PROXIMAL_RHABDOM_ANGLE;
				}
				/* if inter ommatidial angle is zero then we have a perpendicular ray */
				if (inter_ommatidial_angle == 0) {	
					/* CASE 4 - PERPENDICULAR RAY */
					if (tapetal_pigment == 0 || shielding_pigment > 0) {
						fprintf(output, "%0.6f,", (RHABDOM_LENGTH * facet_adjustment));
					}
					else if (tapetal_pigment > 0) {
						fprintf(output, "%0.6f,", ((RHABDOM_LENGTH * 2.0) * facet_adjustment));
					}					
				}
				else {
					/* use tangent of the angle to get the length of the opposite side and divide by radius to get pathlength */
					pathlength = rhabdom_radius / tan(exit_ommatidial_angle * RAD_TO_DEG_CONV);
					if (pathlength >= RHABDOM_LENGTH) {
						/* CASE 3 - BOUNCE OFF THE BASE */
						double pass, passlength;
						if (pathlength == RHABDOM_LENGTH) {
							pass = rhabdom_slope;
						}
						else if (pathlength > RHABDOM_LENGTH) {
							pass = RHABDOM_LENGTH / cos(exit_ommatidial_angle * RAD_TO_DEG_CONV);
						}
						if (pass > RHABDOM_LENGTH) {
							passlength = pass;
						}
						else if (pass < RHABDOM_LENGTH) {
							passlength = RHABDOM_LENGTH;
						}
						if (tapetal_pigment == 0 || shielding_pigment > 0) {
							fprintf(output, "%0.6f,", (pass * facet_adjustment));
						}
						else if (tapetal_pigment > 0) {
							fprintf(output, "%0.6f,", ((pass + passlength) * facet_adjustment));
						}
					}
					else if (pathlength > (RHABDOM_LENGTH - shielding_pigment) || pathlength > (RHABDOM_LENGTH - tapetal_pigment) || exit_ommatidial_angle < critical_angle) {
						/* CASE 2 - REFLECTION FROM THE EDGE */
						double pass1, pass2, passlength;
						pass1 = rhabdom_radius / sin(exit_ommatidial_angle * RAD_TO_DEG_CONV);
						pass2 = (RHABDOM_LENGTH - pathlength) / cos(exit_ommatidial_angle * RAD_TO_DEG_CONV);
						if (pass2 > pass1) {
							pass2 = pass1;
						}
						if ((pass1 + pass2) > RHABDOM_LENGTH) {
							passlength = pass1 + pass2;
						}
						else {
							passlength = RHABDOM_LENGTH;
						}
						/* check tapetum and pigment sizes */
						if (shielding_pigment > (RHABDOM_LENGTH - pathlength)) {
							fprintf(output, "%0.6f,", (pass1 * facet_adjustment));
						}
						else if (tapetal_pigment == 0 || shielding_pigment > 0) {
							fprintf(output, "%0.6f,", ((pass1 + pass2) * facet_adjustment));
						}
						else if (tapetal_pigment > 0) {
							fprintf(output, "%0.6f,", ((pass1 + pass2 + passlength) * facet_adjustment));
						}
					}
					else {
						/* CASE 1 - NO REFLECTION - LIGHT PASSES THROUGH RHABDOM */
						double pass;
						/* calculate length of pass */
						pass = rhabdom_radius / sin(exit_ommatidial_angle * RAD_TO_DEG_CONV);
						fprintf(output, "%0.6f,", (rhabdom_radius / sin(exit_ommatidial_angle * RAD_TO_DEG_CONV)));
						exit_ommatidial_angle = exit_ommatidial_angle + ommatidial_angle; /* increase exit angle */
						increase_acceptance_angle = 1; /* set increase acceptance angle as true */
						if ((RHABDOM_LENGTH - pathlength) > tapetal_pigment || (RHABDOM_LENGTH - pathlength) > shielding_pigment) {
							continue;
						}
					}
				}
				current_facet++; /* increment current facet */
				inter_ommatidial_angle = inter_ommatidial_angle + ommatidial_angle; /* increase inter ommatidial angle */
				increase_acceptance_angle = 0; /* set increase acceptance angle as false */
	
				/* set parameters for next incident facet */
				/* finish row */
				fprintf(output, "998\n"); /* write to file to signify the end of the row */

				/* account for refraction at the cornea */
				if (inter_ommatidial_angle > 60) {
					printf("UNREAL ANGLE AT CORNEA\n"); /* let user know we have an unreal angle at the cornea */
					fprintf(output, "UNREAL ANGLE AT CORNEA\n"); /* write to file that we have an unreal angle at the cornea */
				}
				else if (inter_ommatidial_angle < 15) {
					exit_ommatidial_angle = (inter_ommatidial_angle * 0.9494) + 0.004667;					
				}
				else if (inter_ommatidial_angle < 35) {
					exit_ommatidial_angle = (inter_ommatidial_angle * 0.9407) + 0.1648;					
				}			
				else if (inter_ommatidial_angle < 50) {
					exit_ommatidial_angle = (inter_ommatidial_angle * 0.9196) + 0.8676;
				}
				else if (inter_ommatidial_angle < 60) {
					exit_ommatidial_angle = (inter_ommatidial_angle * 0.8677) + 3.38;
				}

				/* light loss at the cone due to angle of incidence */
				double cone_angle, hypotenuse;
				cone_angle = FACET_WIDTH / tan(exit_ommatidial_angle * RAD_TO_DEG_CONV);
				if (cone_angle > FACET_WIDTH * 2.0) {
					facet_adjustment = cos(inter_ommatidial_angle * RAD_TO_DEG_CONV) * FACET_WIDTH;
				}
				else {
					hypotenuse = ((2.0 * cone_angle) - (2.0 * FACET_WIDTH));
					facet_adjustment = (sin(inter_ommatidial_angle * RAD_TO_DEG_CONV) * hypotenuse);				
				}
				facet_adjustment = (facet_adjustment / FACET_WIDTH);

				/* account for change in angle between adjacent rhabdoms */
				double adjacent_rhabdom_angle;
				adjacent_rhabdom_angle = (number_of_facets / BLUR_CIRCLE_EXTENT);
				for (int i=1; i<=BLUR_CIRCLE_EXTENT; i++){
					if (current_facet > (adjacent_rhabdom_angle * i)) {
						/* miss this rhabdom */
						exit_ommatidial_angle = exit_ommatidial_angle + ommatidial_angle;
						fprintf(output, "0,");
					}
				}
				/* finish row */
				if (current_facet >= number_of_facets) {
					fprintf(output, "998\n"); /* write to file to signify the end of the row */
				}
			}
			/* increment shielding pigment by 10% of rhabdom length */
			shielding_pigment = shielding_pigment + (RHABDOM_LENGTH / 10.0);
			inter_ommatidial_angle = 0.0; /* reset inter ommatidial angle */
			current_facet = 1; /* reset current facet */
			fprintf(output, "999\n"); /* signify the end of the current pair of tapetum/shielding pigment values */
		}
		/* increment tapetal pigment by 10% of rhabdom length */
		tapetal_pigment = tapetal_pigment + (RHABDOM_LENGTH / 10.0);
		shielding_pigment = 0.0; /* reset shielding pigment */
	}
	/* close files */
	fclose(output);
}

/* calculate resolution and sensitivity from pathlengths */
void calculate_ressens(void) {
	/* open files to save values */
	char* filename;
	asprintf(&filename, "%s_resolution.csv", SPECIES_NAME);
	FILE *outres = fopen(filename, "w");
	asprintf(&filename, "%s_sensitivity.csv", SPECIES_NAME);
	FILE *outsens = fopen(filename, "w");
	if (outres == NULL || outsens == NULL)
	{
	    printf("Error opening file!\n");
	    exit(1);
	}

	/* calculate resolution and sensitivity from pathlengths */


	/* close files */
	fclose(outres);
	fclose(outsens);
}

