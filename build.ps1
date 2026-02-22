# Create releases folder
New-Item -ItemType Directory -Force -Path releases | Out-Null

# Disable CGO for clean cross-compiling
$env:CGO_ENABLED="0"

Write-Host "Building Windows..."
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o releases/locked-windows-amd64.exe .
$env:GOOS="windows"; $env:GOARCH="386";   go build -o releases/locked-windows-386.exe .
$env:GOOS="windows"; $env:GOARCH="arm64"; go build -o releases/locked-windows-arm64.exe .

Write-Host "Building Linux..."
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o releases/locked-linux-amd64 .
$env:GOOS="linux"; $env:GOARCH="386";   go build -o releases/locked-linux-386 .
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o releases/locked-linux-arm64 .

Write-Host "Building macOS..."
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o releases/locked-darwin-amd64 .
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o releases/locked-darwin-arm64 .

Write-Host "Done. All builds are in the releases folder."