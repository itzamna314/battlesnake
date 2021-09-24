#!/usr/bin/env bash
if [[ $# -ne 2 ]]; then
	echo "Usage: $0 <game_url> <name>"
	exit 1
fi

RAW_URL=$1
NAME=$2

GAME_ID=$(printf ${RAW_URL} | grep -o -E '[0-9a-f\-]{36}')
TURN=$(printf ${RAW_URL} | sed -r 's/^.*turn=([0-9]+)/\1/')

GAME_JSONL=$(websocat -E "wss://engine.battlesnake.com/games/${GAME_ID}/events")

FRAME=$(echo ${GAME_JSONL} | jq "select(.Data.Turn == ${TURN})")
RULES=$(echo ${GAME_JSONL} | jq 'select(.Type == "game_end")')

OUT=$(cat <<EOF
{
  "game": {
    "id": "unit-test",
    "ruleset": $(echo ${RULES} | jq '.Data.Ruleset'),
	"timeout": $(echo ${RULES} | jq '.Data.SnakeTimeout')
  },
  "turn": ${TURN},
  "board": {
    "height": $(echo ${RULES} | jq '.Data.Height'),
    "width": $(echo ${RULES} | jq '.Data.Width'),
    "food": $(echo ${FRAME} | jq '.Data.Food'),
	"snakes": $(echo ${FRAME} | jq ' [ .Data.Snakes[] | select(.Death == null) | . += {"head": .Body[0], "length": .Body | length } ]'),
    "hazards": $(echo ${FRAME} | jq '.Data.Hazards')
  },
  "you": $(echo ${FRAME} | jq '.Data.Snakes[] | . += {"head": .Body[0], "length": .Body | length } | select(.Name == "Cobra Kai Never Die")') 
}
EOF
)

echo $OUT | jq -c '.' > "frames/${NAME}.json"
