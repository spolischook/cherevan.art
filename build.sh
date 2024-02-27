#!/bin/bash

set -euo pipefail

# Install the npm dependencies
npm install

# Run the build command
npm run build

# Run the image processing script
node ./images_process.js

# Run Hugo with the specified config file
hugo -v --config=./hugo.yaml
