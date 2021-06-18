#!/bin/bash
go get
go mod vendor

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    *)          machine="UNKNOWN:${unameOut}"
esac
if [ $machine == MinGw ] || [ $machine == Cygwin ]; then
    echo @off
    cd vendor/gocv.io/x/gocv
    ./win_build_opencv.cmd
else 
    (cd vendor/gocv.io/x/gocv && make install)
fi