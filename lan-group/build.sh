#!/bin/bash

cd "$(dirname "$0")" || exit
os_all=(windows linux)
arch=amd64
src=lan-group
srcfile=main/main.go
for os in "${os_all[@]}"; do
  [ "$os" == "windows" ] && exe=".exe" || exe=""
  file=../bin/$src/$src-$os-$arch$exe
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o "$file" $srcfile
  echo "output: $file"
done

cp -v "$file" /c/TOOL/
#  C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp
#cp -v script/lan-group-server.bat "/c/ProgramData/Microsoft/Windows/Start Menu/Programs/StartUp"
