#!/usr/bin/env bash
# Aseguramos que CompileDaemon esté en el PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Vigilar toda la raíz del proyecto, compilar y lanzar el server
CompileDaemon \
  -directory="." \
  -build="go build -o tmp/server ./cmd/server/main.go" \
  -command="./tmp/server"
