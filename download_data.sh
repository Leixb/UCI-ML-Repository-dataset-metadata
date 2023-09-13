#!/usr/bin/env bash

curl  -H 'Content-Type: application/json' \
    'https://archive.ics.uci.edu/api/trpc/donated_datasets.search?batch=1&input=%7B%220%22%3A%7B%22json%22%3A%7B%22Area%22%3A%5B%5D%2C%22Keywords%22%3A%5B%5D%2C%22Types%22%3A%5B%5D%2C%22Task%22%3Anull%2C%22NumInstances%22%3Anull%2C%22NumAttributes%22%3Anull%2C%22AttributeTypes%22%3Anull%2C%22skip%22%3A10%2C%22take%22%3A1000%2C%22sort%22%3A%22desc%22%2C%22orderBy%22%3A%22NumHits%22%2C%22search%22%3A%22%22%7D%2C%22meta%22%3A%7B%22values%22%3A%7B%22Task%22%3A%5B%22undefined%22%5D%2C%22NumInstances%22%3A%5B%22undefined%22%5D%2C%22NumAttributes%22%3A%5B%22undefined%22%5D%2C%22AttributeTypes%22%3A%5B%22undefined%22%5D%7D%7D%7D%7D' | \
    jq '.response[0].result.data.json.datasets'
