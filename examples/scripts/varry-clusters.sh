#!/bin/bash

set -e

here="`dirname $0`"
kmpress="$here/../../kmpress"
results="$here/../results"
image="$here/../images/bird_large.tiff"

for i in {2..100}; do
  $kmpress -clusters $i -i $image -o "$results/bird-seq-$i.tiff"
  echo "Completed: $i."
done
