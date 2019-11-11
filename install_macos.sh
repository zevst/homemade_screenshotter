#!/bin/bash
INSTALL_FLD=/opt/hms;
sudo rm -rf $INSTALL_FLD \
    && ln -s $(pwd) $INSTALL_FLD \
    && go build -ldflags "-X main.InstallFld=$INSTALL_FLD" -o scrn scrn.go