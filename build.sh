#!/usr/bin/env bash

dep ensure

echo "Compiling functions to bin/handlers/ ..."

rm -rf bin/

cd handlers/
for f in */; do
  filename="$f${f:0:-1}.go"
  echo $filename
  if GOOS=linux go build -o "../bin/handlers/${f:0:-1}" $filename; then
    echo "✓ Compiled $filename"
  else
    echo "✕ Failed to compile $filename!"
    exit 1
  fi
done

echo "Done."