/* include required header files */
#include <stdio.h>
#include <math.h>
#include <float.h>

/********************************/
/* ONLY CHANGE THESE PARAMETERS */
/********************************/

/* declare global constant eye parameters */
const char* SPECIES_NAME = "astacodes";				/* The species name */
const double RHABDOM_LENGTH = 84.0;					/* The length of the rhabdom */
const double RHABDOM_WIDTH = 16.0;					/* The width of the rhabdom */
const double EYE_DIAMETER = 890.0;					/* The diameter of the eye */
const double FACET_WIDTH = 32.0;					/* The width of the ommatidia's facet */
const double APERTURE_DIAMETER = 445.0;				/* The diameter of the aperture of the eye */
const double CYTOPLASM_REFRACTIVE_INDEX = 1.34;		/* The refractive index of the eye's cytoplasm */
const double RHABDOM_REFRACTIVE_INDEX = 1.37;		/* The refractive index within the eye */
const double BLUR_CIRCLE_EXTENT = 18.0;				/* The diameter of the blur circle */
const double PROXIMAL_RHABDOM_ANGLE = 0;			/* Allow for a different angle at the tips of the rhabdoms */

/*************************/
/* NO CHANGES BELOW HERE */
/*************************/

/* define constant values */
#ifndef M_PI
	#define M_PI 3.14159265358979323846				/* Store value of pi */
#endif
const double RAD_TO_DEG_CONV = M_PI / 180.0;		/* Store conversion from radians to degrees */

/* declare model functions */
void initialise(void);
void run_model(void);

/* declare global model parameters */
double circumference_of_eye;					/* Store the circumference of the eye */
double angle_of_incidence;						/* Store the angle of incident light */
int number_of_facets;							/* Store the number of facets in the aperture */
double angle_of_tir;							/* Store the angle of total internal reflection */
double tapetal_pigment;							/* Store the extent of the tapetal pigment */
double shielding_pigment;						/* Store the extent of the shielding pigment */
double aperture_radius;							/* Store the aperture radius */
double eye_radius;								/* Store the eye radius */
double distance_to_aperture;					/* Store distance from center of the eye to the aperture */
double angle_at_center;							/* Store angle at center of the eye */
double aperture_diameter;						/* Store the aperture diameter */
double optical_axis;							/* Store the optical axis */
int facet_num;									/* Store the current facet number */

/* main subroutine */
int main() {
	initialise();
	printf("Calculating resolution and sensitivity for %s...\n", SPECIES_NAME);
	
	run_model();
}

/* initialise the model parameters */
void initialise(void) {
	/* initial calculations as defined in Magnus' PhD thesis (see page 191) */
	circumference_of_eye = M_PI * EYE_DIAMETER;				/* Calculate the circumference of the eye */
	angle_of_incidence = 0.0;								/* Set incident light to come straight down through the rhabdom */
	number_of_facets = APERTURE_DIAMETER / FACET_WIDTH;		/* Calculate the number of facets */
	angle_of_tir = 0.00;									/* Calculate the angle of total internal reflection */
	tapetal_pigment = 0.0;									/* Initialise tapetal pigment to zero */
	shielding_pigment = 0.0;								/* Initialise shielding pigment to zero */
	
	/* additional calculations required */
	aperture_radius =  APERTURE_DIAMETER / 2.0;						/* Calculate the aperture radius */
	eye_radius = EYE_DIAMETER / 2.0;								/* Calculate the eye radius */
	optical_axis = (FACET_WIDTH / circumference_of_eye) * 360.0;	/* Calculate optical axis */
	/* pythagorus theorum to work out distance to aperture and angle at center */
	distance_to_aperture = sqrt((eye_radius * eye_radius) - (aperture_radius * aperture_radius));
	angle_at_center = atan(aperture_radius / distance_to_aperture);
	aperture_diameter = circumference_of_eye * (angle_at_center / (2 * M_PI));
	/* can also do by converting radians to degrees first */
	/* angle_at_center = atan(aperture_radius / distance_to_aperture) / RAD_TO_DEG_CONV;
	aperture_diameter = (angle_at_center / 360.0) * circumference_of_eye; */
	
	facet_num = 1;
	
	printf("%0.3f\n",  aperture_diameter);
}


/*
self.rhabdom_radius = self.rhabdom_width / 2 # rhabdom radius
self.old_rhabdom_length = self.rhabdom_length # old rhabdom length
self.max_rhabdom_length = self.rhabdom_length # store rhabdom length for main loop
self.inter_ommatidial_angle = 0	# inter-ommatidial angle
self.current_facet = 0 # current facet

# angle of total internal reflection (rhabdoms)
#self.adj = math.sqrt((self.rhabdom_ri ** 2) - (self.cytoplasm_ri ** 2))
#self.cb = math.atan(self.cytoplasm_ri / self.adj)
#self.critical_angle = 90 - (self.cb / self.conv) # critical angle below which light is totally internally reflected within rhabdom
self.snells_law = math.asin(self.cytoplasm_ri / self.rhabdom_ri) / self.conv # calculate angle for total internal reflection using Snell's law
self.critical_angle = 90 - self.snells_law # critical angle below which light is totally internally reflected within rhabdom
self.mx = math.sqrt((self.rhabdom_length ** 2) + (self.rhabdom_radius ** 2)) # mx???
*/


/* Run the model */
void run_model(void) {
	/* iterate through shielding pigment lengths until they are the full length of the rhabdom */
	/* need to add 1 to rhabdom length due to problems with floating point precision */
	while (shielding_pigment <= RHABDOM_LENGTH + 1) {
		/* iterate through tapetal pigment lengths until they are the full length of the rhabdom */
		/* need to add 1 to rhabdom length due to problems with floating point precision */
		while (tapetal_pigment <= RHABDOM_LENGTH + 1) {
			
			//printf("P: %0.3f\n", shielding_pigment);
			//printf("T: %0.3f\n", tapetal_pigment);

			/* iterate over each facet in turn */
			while (facet_num <= number_of_facets) {
				//printf("%d\n", facet_num);
				facet_num++;
			}
			/* increment tapetal pigment by 10% of rhabdom length */
			tapetal_pigment = tapetal_pigment + (RHABDOM_LENGTH / 10.0);
			facet_num = 1; /* reset facet number */
		}
		/* increment shielding pigment by 10% of rhabdom length */
		shielding_pigment = shielding_pigment + (RHABDOM_LENGTH / 10.0);
		tapetal_pigment = 0.0; /* reset tapetal pigment */
	}
}




