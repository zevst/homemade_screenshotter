#!/bin/bash

. .env.dist
test -f .env && . .env

echo "building hmsc app ..."
go build \
  -ldflags "-X main.UploadUrl=$UPLOAD_URL -X main.TmpFolder=$TMP_FOLDER -X main.AccessKey=$ACCESS_KEY" \
  -o "$GOPATH/bin/hmsc" .