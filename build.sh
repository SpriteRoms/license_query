#!/bin/bash

go env -w GOOS=linux GOARCH=amd64
go build -o pack_tool

./pack_tool

go env -w GOOS=js GOARCH=wasm
garble -tiny build -trimpath -o ./web/js/invoice.wasm
go env -w GOOS=windows GOARCH=amd64
