#!/bin/bash

cd kmpress
export GOPATH="`pwd`/lib"
mkdir -p "`pwd`/lib" 2>/dev/null
# Works despite the error.
go get ./...
go build kmpress.go
cd ..

cd voronoi
make
cd ..

cd runner
npm install
cd ..

cd voronoi-runner
npm install
cd ..
