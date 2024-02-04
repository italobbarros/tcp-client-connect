#!/bin/bash

# Adicionando o diretório ao PATH
FILE="tcpclient"
# Criando a pasta de build se não existir
mkdir -p ./build

# Build para Linux
GOOS=linux GOARCH=amd64 go build -o ./build/$FILE

# Build para Windows
GOOS=windows GOARCH=amd64 go build -o ./build/$FILE.exe

echo "Binários foram construídos e estão localizados em ./build."


BINARY="$(pwd)/build/$FILE"
# Adiciona a variável ao final do arquivo de perfil
sudo cp $BINARY /usr/local/bin/
sudo chmod +x /usr/local/bin/$FILE

echo "Alterações aplicadas."


