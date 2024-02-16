#!/bin/bash

# NOTES
# checksum=$(shasum -a 256 /path/to/agent) sed -i "s:^CHECKSUMS=$:CHECKSUMS=\"$checksum\":g" agent.sh && ./agent.sh

# Generate checksums
CHECKSUMS=

# Verify checksums
echo "$CHECKSUMS" | shasum -a 256 -c -s
ret=$?
if [ $ret != 0 ]; then
  echo 'Invalid checksums'
  exit 1
fi

# Deploy agent
# TBD: FIXME

exit 0
