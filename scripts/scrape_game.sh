#!/usr/bin/env bash
if [[ $# -ne 1 ]]; then
	echo "Usage: $0 <game_url>"
	exit 1
fi

RAW_URL=$1
GAME_ID=$(printf ${RAW_URL} | grep -o -E '[0-9a-f\-]{36}')

echo $GAME_ID
