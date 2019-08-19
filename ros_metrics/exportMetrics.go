package main

import (
	"fmt"
	"github.com/performancecopilot/speed"
	"log"
	"sync"
	"time"
)

func exportMetrics(muTopics *sync.Mutex, topicList []string, stopChan chan struct{}) {
	metric, err := speed.NewPCPSingletonMetric(
		0,
		"count",
		speed.Uint64Type,                                                // type
		speed.CounterSemantics,                                         // semantics
		speed.OneUnit,                                                  // unit
		"A Simple Metric",
	)
	if err != nil {
		log.Fatal("Could not create counter, error: ", err)
	}
	client, err := speed.NewPCPClient("rostopic")
	if err != nil {
		log.Fatal("Could not create client, error: ", err)
	}
	err = client.Register(metric)
	if err != nil {
		log.Fatal("Could not register metric, error: ", err)
	}
	client.MustStart()
	defer client.MustStop()
	metricChatter, err := speed.NewPCPCounter(
		0,
		"/chatter",
		"BW of topic",
	)
	clientBW, err := speed.NewPCPClient("rostopic_bw")
	if err != nil {
		log.Fatal("Could not create client, error: ", err)
	}
	err = clientBW.Register(metricChatter)
	if err != nil {
		log.Fatal("Could not register metric, error: ", err)
	}
	clientBW.MustStart()
	defer clientBW.MustStop()
	fmt.Println("The metric should be visible as rostopic.topic_counter")
	for {
		select {
		default:
			muTopics.Lock()
			localList := make([]string, len(topicList))
			copy(localList, topicList)
			muTopics.Unlock()
			localList = deleteEmpty(localList)
			listlen := len(localList)
			fmt.Printf("Current list %v\n", localList)
			fmt.Printf("Current metric val? %d\n", listlen)
			err = metric.Set(uint64(listlen))
			if err != nil {
				log.Fatal("Could not set metric, error: ", err)
			}
			time.Sleep(time.Second)
		case <-stopChan:
			return
		}
	}

}
func deleteEmpty (s []string) []string {
    var r []string
    for _, str := range s {
        if str != "" {
            r = append(r, str)
        }
    }
    return r
}
