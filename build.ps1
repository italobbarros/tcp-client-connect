# Criando a pasta de build se não existir
New-Item -ItemType Directory -Force -Path .\build

# Build para Linux
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o .\build\tcpclient

# Build para Windows
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o .\build\tcpclient.exe


$scriptDir = Get-Location
$env:tcpclient = "$scriptDir\build\tcpclient"
$env:Path += ";$scriptDir\build"
Write-Host "Binários foram construídos e estão localizados em .\build."
Write-Host "variavel de ambiente: $env:tcpclient"