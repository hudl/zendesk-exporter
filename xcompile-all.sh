#!/bin/bash
source $GOPATH/src/github.com/davecheney/golang-crosscompile/crosscompile.bash
go-darwin-386 build -o ./bin/zd-exporter-darwin-386
go-darwin-amd64 build -o ./bin/zd-exporter-darwin-amd64
go-windows-386 build -o ./bin/zd-exporter-windows-386.exe
go-windows-amd64 build -o ./bin/zd-exporter-windows-amd64.exe
go-openbsd-amd64 build -o ./bin/zd-exporter-openbsd-amd64
go-openbsd-386 build -o ./bin/zd-exporter-openbsd-386
go-freebsd-386 build -o ./bin/zd-exporter-freebsd-386
go-freebsd-amd64 build -o ./bin/zd-exporter-freebsd-amd64
go-freebsd-arm build -o ./bin/zd-exporter-freebsd-arm
go-linux-386 build -o ./bin/zd-exporter-linux-386
go-linux-amd64 build -o ./bin/zd-exporter-linux-amd64
go-linux-arm build -o ./bin/zd-exporter-linux-arm
