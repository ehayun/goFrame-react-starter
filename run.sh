#!/bin/bash

# Development run script

# Always build assets from source (assets/ → public/)
echo "Building assets..."
npm run build:assets

# Run the application
echo "Starting Tzlev server..."
go run main.go
