#!/bin/bash

# Usage: Generate checksum
# shasum -a 256 agent

# Usage: Deploy agent
# ./agent.sh "$ARTIFACT_USER" "$ARTIFACT_PASS" "$ARTIFACT_URL" "$ARTIFACT_PATH", "$AGENT_EXEC", "$AGENT_PATH_AGENT_EXEC"

# Install jq
jq --version > /dev/null
ret=$?
if [ "$ret" -ne 0 ]; then
  sudo apt install -y jq > /dev/null
fi

# Fetch checksum
CHECKSUM=$(curl -f -s -u "$1":"$2" "$3/api/storage/$4/$5" | jq '.checksums.sha256' | tr -d '"')

# Verify checksum
echo "$CHECKSUM $5" > "$5".checksum
sha256sum --ignore-missing --status -c "$5".checksum
ret=$?
rm -rf "$5".checksum
if [ $ret -eq 0 ]; then
  echo 'Checksum pass'
  exit 0
fi

# Deploy agent
curl -f -s -u"$1":"$2" -L "$3/$4/$5" -o "$6"
ret=$?
if [ $ret != 0 ]; then
  echo 'Missing agent'
  exit 1
fi

chmod +x "$6"

exit 0
