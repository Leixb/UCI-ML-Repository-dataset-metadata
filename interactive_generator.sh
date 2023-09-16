#!/usr/bin/env bash

FILE="${FILE:-./data.json}"
INPUT_FILE="${1:-$FILE}"

UCIML="${UCIML:-./uciml}"

$UCIML --serve "${INPUT_FILE}" >/dev/null 2>&1 &
pid=$!

sleep 1 # wait for server to start

mkdir -p generated

curl -H 'Accept: text/plain' http://localhost:8080/datasets | \
    fzf --multi --preview "curl -s -H 'Accept: text/markdown' http://localhost:8080/datasets/{2} | bat -l md --color always" | \
    cut -d' ' -f2 | tee generated/selection.txt | \
    xargs "$UCIML" --verbose --julia --toml "${INPUT_FILE}"

kill $pid
