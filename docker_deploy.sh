#!/bin/bash

# Build the Docker image
echo "Building Docker image..."
docker build -t 9919952021/ltp_updater:latest .

# Check if the build was successful
if [ $? -ne 0 ]; then
  echo "Docker build failed. Exiting."
  exit 1
fi

# Log in to Docker Hub
echo "Logging in to Docker Hub..."
docker login

# Check if the login was successful
if [ $? -ne 0 ]; then
  echo "Docker login failed. Exiting."
  exit 1
fi

# Push the Docker image to Docker Hub
echo "Pushing Docker image..."
docker push 9919952021/ltp_updater:latest

# Check if the push was successful
if [ $? -ne 0 ]; then
  echo "Docker push failed. Exiting."
  exit 1
fi

echo "Docker image successfully pushed."