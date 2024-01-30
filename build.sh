#!/bin/bash

# Criando a pasta de build se não existir
mkdir -p ./build

# Build para Linux
GOOS=linux GOARCH=amd64 go build -o ./build/tcpclient

# Build para Windows
GOOS=windows GOARCH=amd64 go build -o ./build/tcpclient.exe

echo "Binários foram construídos e estão localizados em ./build."
