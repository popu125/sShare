#!/bin/bash
mkdir -p dist/utils
rm -rf dist/*
for arch in "386" "amd64"; do
	echo "Start building $arch."
	export GOARCH=$arch
	export GOOS=linux
	go build -o dist/sShare -ldflags="-s -w" github.com/popu125/sShare 
	go build -o dist/utils/simplepage-gen -ldflags="-s -w" github.com/popu125/sShare/utils/simplepage-gen 
	cp config.example.json dist/config.json
	cp utils/simplepage-gen/index.tpl dist/utils/index.tpl
	upx dist/sShare dist/utils/simplepage-gen
	echo "Packing $arch."
	tar -C dist/ -czf build-$GOARCH.tar.gz $(ls dist)
done
echo "Build Done."