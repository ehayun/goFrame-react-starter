#!/bin/bash

set -e

echo "Building Tzlev application..."

# Build template assets
echo "Building template assets..."
npm install
npm run build:assets

# Build frontend
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Build backend
echo "Building backend..."
go mod tidy
go build -o bin/tzlev main.go

echo "Build complete!"
echo "Run with: ./bin/tzlev"
