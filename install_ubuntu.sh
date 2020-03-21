#!/bin/bash
. .env

rm -rf "$INSTALL_FLD" && cp -r "$(pwd)" "$INSTALL_FLD"