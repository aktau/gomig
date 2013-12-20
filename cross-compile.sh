#!/bin/bash

set -e

BIN="./bin"
APP="gomig"
APP_PATH="$BIN/$APP"

[ -d "$BIN" ] || mkdir "$BIN"

echo '# linux arm7'
GOARM=7 GOARCH=arm GOOS=linux go build -o "$APP_PATH"-arm5
echo '# linux arm5'
GOARM=5 GOARCH=arm GOOS=linux go build -o "$APP_PATH"-arm7
echo '# windows 386'
GOARCH=386 GOOS=windows go build -o "$APP_PATH"-x86.exe
echo '# windows amd64'
GOARCH=amd64 GOOS=windows go build -o "$APP_PATH"-x64.exe
echo '# darwin'
GOARCH=amd64 GOOS=darwin go build -o "$APP_PATH"-osx
echo '# freebsd'
GOARCH=amd64 GOOS=freebsd go build -o "$APP_PATH"-freebsd
