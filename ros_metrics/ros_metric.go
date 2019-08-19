package main

/*
#cgo CFLAGS: -I/usr/include/bcc/compat
#cgo LDFLAGS: -lbcc
#include <bcc/bcc_common.h>
#include <bcc/libbpf.h>
void perf_reader_free(void *ptr);
*/
import "C"

import (
	"flag"
	"sync"
	"time"
)

//const interval = time.Millisecond

var timelimit = flag.Int("time", 600, "number of seconds to run for")

func main() {
	flag.Parse()
	var (
		muTopics  sync.Mutex
		//TODO: MAGIC NUMBER BAD
		//topicList []string
		topicList = make([]string, 256)
		stopChan  = make(chan struct{})
	)

	go MonitorROS(&muTopics, topicList, stopChan)
	go exportMetrics(&muTopics, topicList, stopChan)
	for i := 0; i < *timelimit; i++ {
		time.Sleep(time.Second)
	}
	close(stopChan)

}
