#!/bin/bash

set -e

here="`dirname $0`"
kmpress="$here/../../kmpress/kmpress"
results="$here/../results/vary-clusters"
image="$here/../images/bird_large.tiff"

mkdir -p "$results" 2>/dev/null

for i in {2..100}; do
  $kmpress -clusters $i -i $image -o "$results/bird-seq-$i.tiff"
  echo "Completed: $i."
done
