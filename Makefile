# Define the go command
GOCMD = go

# Phony targets are not actual files; they are just a name for a set of commands
.PHONY: tidy

tidy:
	$(GOCMD) mod tidy
