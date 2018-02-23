#!/usr/bin/env bash
for r in {1..8}; do
  for c in {1..8}; do
    printf "Result for row %d, col %d: " $r $c
    go run *.go < words.txt $r $c
  done
done
