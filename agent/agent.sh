#!/bin/bash

# Usage: Generate checksum
# checksum=$(shasum -a 256 /path/to/agent) sed -i "s:^CHECKSUM=$:CHECKSUM=\"$checksum\":g" agent.sh

# Usage: Deploy agent
# ./agent.sh "$ARTIFACT_USER" "$ARTIFACT_PASS" "$ARTIFACT_URL" "$AGENT_PATH"

# Generate checksum
CHECKSUM=

# Verify checksum
echo "$CHECKSUM" | shasum -a 256 -c -s
ret=$?
if [ $ret != 0 ]; then
  echo 'Invalid checksum'
  exit 1
fi

# Deploy agent
curl -f -s -u"$1":"$2" -L "$3" -o "$4"
ret=$?
if [ $ret != 0 ]; then
  echo 'Missing agent'
  exit 2
fi

chmod +x "$4"

exit 0
