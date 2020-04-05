#!/bin/bash

echo "set up hmsc desktop entry ..."
cp hmsc.png ~/.local/share/icons/
envsubst < hmsc.desktop > ~/.local/share/applications/hmsc.desktop