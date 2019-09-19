#!/bin/bash 
#export COUNTER=1
mkdir results
export NODES=60

for COUNTER in {1..15}
do
    for SAMPLES in 1 5 10 50 100 200 500 1000
    do
    	echo "$COUNTER + $NODES"
    	#perf stat -o results/control.res$COUNTER.node$NODES -e cycles,instructions,cache-references,cache-misses,bus-cycles -d ./control.sh $NODES 2>&1
    	#perf stat -o results/expirement1.res$COUNTER.node$NODES -e cycles,instructions,cache-references,cache-misses,bus-cycles -d ./exec_metrics.sh $NODES 2>&1
    	perf stat -o results/break_bw.res$COUNTER.node$SAMPLES -e cycles,instructions,cache-references,cache-misses,bus-cycles -d ./break_topic_bw.sh $NODES $SAMPLES 2>&1
    	#perf stat -o results/topic_delay.res$COUNTER.node$NODES -e cycles,instructions,cache-references,cache-misses,bus-cycles -d ./get_topic_delay.sh $NODES 2>&1
    	#perf stat -o results/topic_hz.res$COUNTER.node$NODES -e cycles,instructions,cache-references,cache-misses,bus-cycles -d ./get_topic_hz.sh $NODES 2>&1
    	perf stat -o results/break_all.res$COUNTER.node$SAMPLES -e cycles,instructions,cache-references,cache-misses,bus-cycles -d ./break_topic_all.sh $NODES $SAMPLES 2>&1
   	killall -9 roscore rosout rosmaster python2 
    done
done
