#!/usr/bin/env bash

OLD=0
NEW=0

for i in {0..99}; do
	WINNER=$(battlesnake play -W 11 -H 11 -n old -u http://localhost:8081 -n new -u http://localhost:8080 2>&1 >/dev/null | grep '[DONE]' | sed -E 's/^.*\. (.+) is the winner\./\1/')

	if [[ ${WINNER} == 'old' ]]; then
		let OLD++
	elif [[ ${WINNER} == 'new' ]]; then
		let NEW++
	else
		echo "Unexpected winner ${WINNER}"
		exit 1
	fi

	echo $WINNER
done

echo "Old: ${OLD}, New: ${NEW}"
