#!/bin/bash

set -euo pipefail

# Install the npm dependencies
npm install

# Run the build command
npm run build

# List the files in the current directory
ls -l

# Run Hugo with the specified config file
hugo -v --config=./hugo.yaml

# Run the image processing script
node ./images_process.js

# List the files in the public directory
tree public