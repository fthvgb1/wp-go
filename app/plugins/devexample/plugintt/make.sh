#/bin/bash

# copy plugintt to other dir and remove .dev suffix
# note the go version and build tool flag must same to server build
# eg: -gcflags all="-N -l" may used in ide debug
go build -buildmode=plugin -o xx.so main.go