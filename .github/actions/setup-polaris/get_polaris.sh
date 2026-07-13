#!/bin/bash
if [[ -z "$INPUT_VERSION" ]]; then
  echo "Missing polaris version information"
  exit 1
fi
# Release assets embed the version number without the leading "v"
# (e.g. tag v10.2.0 -> polaris_10.2.0_linux_amd64.tar.gz).
POLARIS_ASSET_VERSION="${INPUT_VERSION#v}"
POLARIS_URL=https://github.com/FairwindsOps/polaris/releases/download/$INPUT_VERSION/polaris_${POLARIS_ASSET_VERSION}_linux_amd64.tar.gz
polaris version | grep "$INPUT_VERSION" &> /dev/null
if [ $? == 0 ]; then
   echo "Polaris $INPUT_VERSION is already installed! Exiting gracefully."
   exit 0
else
  echo "Installing polaris to path from " $POLARIS_URL
fi
TARGET_FILE="polaris.tar.gz"
curl -LJ -o $TARGET_FILE $POLARIS_URL
mkdir polaris
tar -xzf $TARGET_FILE -C polaris
rm $TARGET_FILE
echo "polaris" >> $GITHUB_PATH
echo "version=$INPUT_VERSION" >> $GITHUB_OUTPUT