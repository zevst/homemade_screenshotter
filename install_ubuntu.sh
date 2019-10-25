#!/bin/bash
INSTALL_FLD=/opt/hms;
sudo ln -s $(pwd) $INSTALL_FLD \
    && ln homemade-screenshotter.desktop ~/.local/share/applications \
    && go build -ldflags "-X main.InstallFld=$INSTALL_FLD" -o scrn scrn.go