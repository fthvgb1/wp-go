#!/bin/bash

# copy plugintt directory to other dir which as same level with wp-go and remove .dev suffix
# note the go version and build tool flag must same to server build
# replace wp-go => ../wp-go in go.mod
#  -gcflags all="-N -l" --race can be used in ide debug
# wp-go config add xx plugin
go mod tidy
go build -buildmode=plugin -o xx.so main.go