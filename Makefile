# Define your TailwindCSS input and output files
TAILWIND_INPUT = ./static/css/input.css
TAILWIND_OUTPUT = ./static/css/output.css

# Default target that runs when you just run 'make'
all: build

# The build target runs the TailwindCSS command
tw:
	tailwindcss -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT) --watch

# prepare my files to be fed to ai for questions
gpt:
	cat *.go > gpt.txt