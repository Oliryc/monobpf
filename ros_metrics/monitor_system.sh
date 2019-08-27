#!/bin/bash 
#export COUNTER=1
mkdir results
for COUNTER in {1..15}
do
    for NODES in 1 5 10 20 40 60 80
    do
    	echo "$COUNTER + $NODES"
    	perf stat -o results/expirement1.res$COUNTER.node$NODES -e cycles,instructions,cache-references,cache-misses,bus-cycles -d ./exec_metrics.sh $NODES 2>&1
    
    done
done
