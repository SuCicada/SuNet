#!/bin/bash


cd "$(dirname "$0")" || exit
os=windows
arch=amd64
exe=.exe
src=lan-group
file=../bin/$src/$src-$os-$arch$exe
srcfile=main/main.go
CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o $file $srcfile

cp -v $file /c/TOOL/
#  C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp
cp -v script/lan-group-server.bat "/c/ProgramData/Microsoft/Windows/Start Menu/Programs/StartUp"
