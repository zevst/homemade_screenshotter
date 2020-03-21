#!/bin/bash
. .env

cp -f homemade-screenshotter.desktop ~/.local/share/applications \
    && go build -ldflags "-X main.InstallFld=$INSTALL_FLD" -o scrn scrn.go