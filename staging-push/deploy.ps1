# Set environment variables for cross-compilation
$env:GOOS = "linux"
$env:GOARCH = "amd64"

# Define variables
$binaryName = "mangobyte-site"         # The name of the Go binary
$localBinaryPath = "./bin/$binaryName" # Path to the built binary on local machine
$localStaticPath = "./ui/static/"      # Path to the static folder on local machine
$localEnvFile = "./.env"               # Path to the .env file on local machine
$localTLS = "./tls/"
$localScrape = "./scraper/"
$localUvLock = "./uv.lock"                   # Path to the uv lock file
$stagingServer = "staging-server"            # Staging server IP/hostname
$remoteDir = "/home/mangobyte-site"          # Directory on the staging server
$remoteBinaryPath = "$remoteDir/$binaryName" # Full path of the binary on the staging server
$remoteStaticPath = "$remoteDir/ui/"         # Directory on the staging server for static files
$remoteEnvFilePath = "$remoteDir/.env"       # Path to the .env file on the staging server

# Step 1: Build Go binary for Linux
Write-Host "Building Go binary for Linux..."
go build -o "$localBinaryPath" ./cmd/web/
if ($LASTEXITCODE -ne 0) {
    Write-Host "Go build failed! Aborting deployment."
    exit 1
}

# Step 2: Stop the service before deploying
Write-Host "Stopping the mangobyte-site service on the staging server..."
ssh $stagingServer "sudo systemctl stop mangobyte-site"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to stop the service on the staging server."
    exit 1
}

# Step 3: Transfer the binary to the staging server using SCP
Write-Host "Transferring binary to staging server..."
scp "$localBinaryPath" "${stagingServer}:$remoteDir"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to transfer binary! Aborting deployment."
    exit 1
}

# Step 4: Transfer the static folder to the staging server
Write-Host "Transferring static assets to staging server..."
scp -r "$localStaticPath" "${stagingServer}:$remoteStaticPath"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to transfer static assets! Aborting deployment."
    exit 1
}

# Step 5: Transfer the .env file to the root of the remote directory
Write-Host "Transferring .env file to staging server..."
scp "$localEnvFile" "${stagingServer}:$remoteEnvFilePath"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to transfer .env file! Aborting deployment."
    exit 1
}

# Step 6: Transfer TLS certificates to the staging server
Write-Host "Transferring TLS certificates to staging server..."
scp -r "$localTLS" "${stagingServer}:$remoteDir"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to transfer TLS certificates! Aborting deployment."
    exit 1
}

# Step 7: Transfer scraper folder and uv.lock file to staging server
Write-Host "Transferring scraper content and uv.lock file to staging server..."
scp -r "$localScrape" "${stagingServer}:$remoteDir"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to transfer scraper directory! Aborting deployment."
    exit 1
}
scp "$localUvLock" "${stagingServer}:$remoteDir/uv.lock"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to transfer uv.lock! Aborting deployment."
    exit 1
}

# Step 8: Install `uv` and dependencies using `requirements.txt`
Write-Host "Setting up Python environment and installing dependencies using 'uv' on the staging server..."
ssh $stagingServer "sudo apt-get update && sudo apt-get install -y python3-pip python3-dev python3-venv && if ! command -v uv &> /dev/null; then pip3 install --user uv; fi && cd $remoteDir && if [ ! -d '.venv' ]; then uv venv .venv; fi && uv pip install -r requirements.txt"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to set up Python environment or install dependencies using 'uv'!"
    exit 1
}

# Step 9: Make the binary executable on the staging server
Write-Host "Making the binary executable..."
ssh $stagingServer "chmod +x $remoteBinaryPath"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to make the binary executable!"
    exit 1
}

# Step 10: Restart the service after deployment
Write-Host "Restarting the mangobyte-site service on the staging server..."
ssh $stagingServer "sudo systemctl restart mangobyte-site"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to restart the service on the staging server."
    exit 1
}

Write-Host "Deployment complete and app is running on the staging server!"
