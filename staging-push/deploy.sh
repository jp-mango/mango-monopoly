#!/bin/bash

# Set environment variables for cross-compilation
export GOOS="linux"
export GOARCH="amd64"

# Define variables
binaryName="mangobyte-site"         # The name of the Go binary
localBinaryPath="./bin/$binaryName" # Path to the built binary on local machine
localStaticPath="./ui/static/"      # Path to the static folder on local machine
localEnvFile="./.env"               # Path to the .env file on local machine
localTLS="./tls/"
localScrape="./scraper/"
stagingServer="staging-server"            # Staging server IP/hostname
remoteDir="/home/mangobyte-site"          # Directory on the staging server
remoteBinaryPath="$remoteDir/$binaryName" # Full path of the binary on the staging server
remoteStaticPath="$remoteDir/ui/"         # Directory on the staging server for static files
remoteEnvFilePath="$remoteDir/.env"       # Path to the .env file on the staging server

# Step 1: Build Go binary for Linux
echo "Building Go binary for Linux..."
go build -o "$localBinaryPath" ./cmd/web/

if [ $? -ne 0 ]; then
	echo "Go build failed! Aborting deployment."
	exit 1
fi

# Step 2: Stop the service before deploying
echo "Stopping the mangobyte-site service on the staging server..."
ssh "$stagingServer" "sudo systemctl stop mangobyte-site"

if [ $? -ne 0 ]; then
	echo "Failed to stop the service on the staging server."
	exit 1
fi

# Step 3: Transfer the binary to the staging server using SCP
echo "Transferring binary to staging server..."
scp "$localBinaryPath" "${stagingServer}:${remoteDir}"

if [ $? -ne 0 ]; then
	echo "Failed to transfer binary! Aborting deployment."
	exit 1
fi

# Step 4: Transfer the static folder to the staging server
echo "Transferring static assets to staging server..."
scp -r "$localStaticPath" "${stagingServer}:${remoteStaticPath}"

if [ $? -ne 0 ]; then
	echo "Failed to transfer static assets! Aborting deployment."
	exit 1
fi

# Step 5: Transfer the .env file to the root of the remote directory
echo "Transferring .env file to staging server..."
scp "$localEnvFile" "${stagingServer}:${remoteEnvFilePath}"

if [ $? -ne 0 ]; then
	echo "Failed to transfer .env file! Aborting deployment."
	exit 1
fi

# Step 6: Transfer TLS certificates to the staging server
echo "Transferring TLS certificates to staging server..."
scp -r "$localTLS" "${stagingServer}:${remoteDir}"

if [ $? -ne 0 ]; then
	echo "Failed to transfer TLS certificates! Aborting deployment."
	exit 1
fi

# Step 7: transfer scraper folder and its content to staging server
echo "Transferring scraped content and scripts to staging"
scp -r "$localScrape" "${stagingServer}:${remoteDir}"

if [ $? -ne 0 ]; then
	echo "Failed to transfer scraper directory! Aborting deployment."
	exit 1
fi

# Step 8: Make the binary executable on the staging server
echo "Making the binary executable..."
ssh "$stagingServer" "chmod +x $remoteBinaryPath"

if [ $? -ne 0 ]; then
	echo "Failed to make the binary executable!"
	exit 1
fi

# Step 9: Restart the service after deployment
echo "Restarting the mangobyte-site service on the staging server..."
ssh "$stagingServer" "sudo systemctl restart mangobyte-site"

if [ $? -ne 0 ]; then
	echo "Failed to restart the service on the staging server."
	exit 1
fi

echo "Deployment complete and app is running on the staging server!"
