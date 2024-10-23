# Set environment variables for cross-compilation
$env:GOOS = "linux"
$env:GOARCH = "amd64"

# Define variables
$binaryName = "mangobyte-site"   # The name of the Go binary
$localBinaryPath = ".\bin\$binaryName"  # Path to the built binary on local machine
$localStaticPath = ".\ui\static\"       # Path to the static folder on local machine
$localEnvFile = ".\.env"                # Path to the .env file on local machine
$stagingServer = "staging-server"       # Staging server IP/hostname
$remoteDir = "/home/mangobyte-site"     # Directory on the staging server
$remoteBinaryPath = "$remoteDir/$binaryName"  # Full path of the binary on the staging server
$remoteStaticPath = "$remoteDir/ui/"       # Directory on the staging server for static files
$remoteEnvFilePath = "$remoteDir/.env"     # Path to the .env file on the staging server

# Step 1: Build Go binary for Linux
Write-Host "Building Go binary for Linux..."
go build -o $localBinaryPath .\cmd\web\

if ($LASTEXITCODE -ne 0) {
    Write-Error "Go build failed! Aborting deployment."
    exit 1
}

# Step 2: Stop the service before deploying
Write-Host "Stopping the mangobyte-site service on the staging server..."
ssh $stagingServer "sudo systemctl stop mangobyte-site"

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to stop the service on the staging server."
    exit 1
}

# Step 3: Transfer the binary to the staging server using SCP
Write-Host "Transferring binary to staging server..."
scp $localBinaryPath "${stagingServer}:${remoteDir}"

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to transfer binary! Aborting deployment."
    exit 1
}

# Step 4: Transfer the static folder to the staging server
Write-Host "Transferring static assets to staging server..."
scp -r $localStaticPath "${stagingServer}:${remoteStaticPath}"

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to transfer static assets! Aborting deployment."
    exit 1
}

# Step 5: Transfer the .env file to the root of the remote directory
Write-Host "Transferring .env file to staging server..."
scp $localEnvFile "${stagingServer}:${remoteEnvFilePath}"

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to transfer .env file! Aborting deployment."
    exit 1
}

# Step 6: Make the binary executable on the staging server
Write-Host "Making the binary executable..."
ssh $stagingServer "chmod +x $remoteBinaryPath"

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to make the binary executable!"
    exit 1
}

# Step 7: Restart the service after deployment
Write-Host "Restarting the mangobyte-site service on the staging server..."
ssh $stagingServer "sudo systemctl restart mangobyte-site"

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to restart the service on the staging server."
    exit 1
}

Write-Host "Deployment complete and app is running on the staging server!"
