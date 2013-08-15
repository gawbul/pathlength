***
# **Superposition Eye Path-length Program**
# 
***
# 

A program implementing a ray tracing model to determine the interactions of various parameters and their impacts on the pathlengths of POLs through a reflecting superposition compound eye.

Original QBASIC version by Magnus L Johnson and Genevre Parker, 1995

Python rewrite by Stephen P Moss, 2012-2013

http://about.me/gawbul

gawbul@gmail.com

#Installation

Install the following dependencies:

# Usage

Run the program using:

	python pathlength.py

**N.B.: Implemented using Python 2.7**

**Use the following to check your python version**:

	$ python -V
	Python 2.7.3

If you have a version less than 2.7, you can get the latest python version for your system from http://www.python.org/.

The default settings are located in the main function:

	# main handler subroutine
	def main():
		# check what the program arguments are and assign appropriate variables
		opts_array = handle_options(sys.argv[1:])
		(input_file, graphicsopt) = opts_array
	
		# check whether the user provide an input filename
		if input_file:
			# process file
			process_input_file(input_file, graphicsopt)
			sys.exit()
		else:
			# just continue with inline parameters below
			pass
	
		# show startup information
		startup()

		# track how long it takes
		start = time.time()
		
		# if not using an input file for the parameters you can set them manually as follows
		# setup nephrops_eye as new SuperpositionEye object - with relevant parameters passed	
		# using Nephrops norvegicus flat lateral measurments
		# see README file or GitHub for information on parameters
		print "Setting up new superposition eye object..."
		nephrops_eye = SuperpositionEye("nephrops", 180, 25, 7800, 50, 3200, 1.34, 1.37, 18, 0) 

		# run the model
		print "Running the ray tracing model (please wait)..."
		nephrops_eye.run_model(graphicsopt)
	
		# summarise the data
		print "Outputting summary data..."
		nephrops_eye.summarise_data()
	
		# how long did we take?
		end = time.time()
		took = end - start
		print "\nFinished in %s seconds.\n" % timedelta(seconds=took)

To setup a new eye object you need to do the following before executing the program:

	eye_object = SuperpositionEye("genus", 180, 25, 7800, 50, 3200, 1.34, 1.37, 18, 0) 

Where the parameters equal:

	"genus"	=	A prefix for the output filenames e.g. organism genus name
	180 	=	Rhabdom Length
	25 	=	Rhabdom Width
	7800 	=	Eye Diameter
	50 	=	Facet Width
	3200		=	Aperture Diameter
	1.34	=	Cytoplasm Refractive Index
	1.37	=	Rhabdom Refractive Index
	18		=	Blur Circle Extent
	0		=	Proximal Rhabdom Angle (used to create pointy-ended rhabdoms)

**N.B.: The genus name is NOT case sensitive. It is always converted to lowercase to avoid file access issues.**

Then you need to simply call the run_model method in order to execute the model.

	eye_object.run_model()

This outputs two files (where genus is the name you give when setting up the object):

	genus_output_one.txt	=	Each record is separated by 999 in the text and contains the length of the reflective tapetum and shielding pigment initially, followed by the path length values for each rhabdom the light passes through, starting at the axial rhabdom.
	genus_output_two.txt	=	**Description needed**

In order to summarise the data one is required to call the summarise_data method:

	eye_object.summarise_data()
	
This outputs three files (where genus is the name you give when setting up the object):

	genus_summary_one.txt	= **Description needed**
	genus_summary_res.txt	= Resolution output
	genus_summary_sen.txt	= Sensitivity output

#Command line options

The program allows you to input certain command line options when executing the program:

	e.g. python pathlength.py -v

The options that are available currently are:

	f	=	file (also --file)
	g	=	graphics (also --graphics)
	c	=	citation (also --citation)
	h	=	help (also --help)
	v	=	version (also --version)

These options have the following effects:

	file		=	Allows the user to give a filename containing parameters in comma separated value format, with individual sets of parameters on separate lines. The program will parse each line of the file in turn, running the model for each set of parameters.
	graphics	=	Allows the user to view graphical output. *** not yet implemented ***
	citation	=	Allows the user to view the citation information.
	help		=	Allows the user to view the usage information.
	version		=	Allows the user to view the version of the program.

By providing an input file, you can implement a workflow, testing various different eye parameters and thus different hypotheses. The file should be in the following format and follows the same structure as using the object within the program, as shown above:

	nephropsfl,180,25,7800,50,3200,1.34,1.37,18,0
	nephropspl,180,25,7800,50,3200,1.34,1.37,18,12.5
	nephropsfa,180,25,6760,50,3060,1.34,1.37,10,0
	nephropspa,180,25,6760,50,3060,1.34,1.37,10,12.5
