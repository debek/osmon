#!/bin/bash

# Define the current directory as the base directory for the osmon package
BASE_DIR=$(pwd)

# Build the osmon.deb package
dpkg-deb --build "$BASE_DIR"

mv ../osmon-deb.deb "$BASE_DIR/osmon.deb"

echo "The osmon.deb package has been successfully created in $BASE_DIR"
