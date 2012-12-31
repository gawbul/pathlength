***
# 
# **Superposition Eye Path-length Program**
# 
***
# 

A program implementing a ray tracing model to determine the interactions of various parameters and their impacts on path's of POLs through reflecting superposition compound eye.

Original QBASIC version by Magnus L Johnson and Genevre Parker, 1995

Python rewrite by Stephen P Moss, 2012

http://about.me/gawbul

gawbul@gmail.com

# 
***
# 
# USAGE
# 
***
# 

Run the program using (**N.B. Implemented using Python 2.7**):

	python pathlen.py

** Use the following to check your python version**:

	$ python -V
	Python 2.7.3

If you have a version less than 2.7, you can get the latest python version for your system from http://www.python.org/.

The default settings are located in the main function:

	# main handler subroutine
	def main():
		# check what the program arguments are and assign appropriate variables
		opts_array = handle_options(sys.argv[1:])
		input_file, graphics_opt = opts_array
		
		# check whether the user provide an input filename
		if input_file:
			# process file
			process_input_file(input_file, graphics_opt)
			sys.exit()
		else:
			# just continue with inline parameters below
			pass
		
		# if not using an input file for the parameters you can set them manually as follows
		# setup acanthephyra_eye as new SuperpositionEye object - with relevant parameters passed	
		# using Acanthephyra pelagica measurments
		# also test using Systellaspis debilis and Nephrops norvegicus
		acanthephyra_eye = SuperpositionEye("acanthephyra", 127, 15.8, 2480, 22.5, 870, 1.34, 1.37, 1, 0.0) 

		# run the model	
		acanthephyra_eye.run_model(graphics_opt)
		
		# summarise the data
		acanthephyra_eye.summarise_data()


			
To setup a new eye object you need to do the following before executing the program:

	eye_object = SuperpositionEye("genus", 127, 15.8, 2480, 22.5, 870, 1.34, 1.37, 1, 0) 

Where the parameters equal:

	"genus"	=	A prefix for the output filenames e.g. organism genus name
	127 	=	Rhabdom Length
	15.8 	=	Rhabdom Width
	2480 	=	Eye Diameter
	22.5 	=	Facet Width
	870		=	Aperture Diameter
	1.34	=	Cytoplasm Refractive Index
	1.37	=	Rhabdom Refractive Index
	1		=	Blur Circle Extent
	0		=	Proximal Rhabdom Angle (used to create pointy-ended rhabdoms)

**N.B.: The genus name is NOT case sensitive. It is always converted to lowercase to avoid file access issues.**

Then you need to simply call the run_model method in order to execute the model.

	eye_object.run_model()

This outputs two files (where genus is the name you give when setting up the object):

	genus_output_one.txt	=	Each record is separated by 999 in the text and contains the length of the reflective tapetum and shielding pigment initially, followed by the path length values for each rhabdom the light passes through, starting at the axial rhabdom.
	
	genus_output_two.txt	=	Description needed

In order to summarise the data one is required to call the summarise_data method:

	eye_object.summarise_data()
	
This outputs three files (where genus is the name you give when setting up the object):

	genus_summary_one.txt	= Description needed
	genus_summary_res.txt	= Resolution output
	genus_summary_sen.txt	= Sensitivity output

# 
***
# 
# Command line options
# 
***
# 

The program allows you to input certain command line options when executing the program:

	e.g. python pathlen.py -v

The options that are available currently are:

	f	=	file (also --file)
	g	=	graphics (also --graphics)
	h	=	help (also --help)
	v	=	version (also --version)

These options have the following effects:

	file	=	Allows the user to give a filename containing parameters in comma separated value format, with individual sets of parameters on separate lines. The program will parse each line of the file in turn, running the model for each set of parameters.
	graphics =	Allows the user to view a realtime visual representation of the parameters being calculated (**currently not implemented**).
	help	=	Allows the user to view the usage information.
	version	=	Allows the user to view the version of the program.


By providing an input file, you can implement a workflow, testing various different eye parameters and thus different hypotheses. The file should be in the following format and follows the same structure as using the object within the program, as shown above:

	acanthephyra,127,15.8,2480,22.5,870,1.34,1.37,1,0
	acanthephyra_bce5,127,15.8,2480,22.5,870,1.34,1.37,5,0
	acanthephyra_bce10,127,15.8,2480,22.5,870,1.34,1.37,10,0
