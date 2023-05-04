#/bin/bash

# note the build tool flag must same to server build .
# eg: -gcflags all="-N -l" may used in ide debug
go build -buildmode=plugin -o xx.so main.go