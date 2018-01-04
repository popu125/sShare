#!/bin/bash
mkdir -p dist/utils
go build -o dist/sShare github.com/popu125/sShare 
go build -o dist/utils/simplepage-gen github.com/popu125/sShare/utils/simplepage-gen 
cp config.example.json dist/config.json
cp utils/simplepage-gen/index.tpl dist/utils/index.tpl
upx -9 dist/*
tar czf build-$GOARCH.tar.gz dist/*