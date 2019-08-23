#!/bin/bash 
export COUNTER=1
perf stat -e cycles,instructions,cache-references,cache-misses,bus-cycles -d ./exec_metrics.sh  

