#!/bin/bash

# Build the Docker image
docker build -t simplealloc .

# Run the Docker container with memory constraints
docker run --name simplealloc_container --memory=6m --memory-swap=6m simplealloc

# Copy the run file from the container to the host machine
docker cp simplealloc_container:/app/runs ./runs

# Clean up the container
docker rm simplealloc_container