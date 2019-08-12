package exporter

import (
	"fmt"
	"github.com/performancecopilot/speed"
	"log"
	"os"
	"strconv"
	"time"
)

func exportMetrics() {
	metric, err := speed.NewPCPCounter(
		0,
		"count",
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
	for i := 0; i < *timelimit; i++ {
		metric.Up()
		metricChatter.Up()
		metricChatter.Up()
		metricChatter.Up()
		time.Sleep(time.Second)
	}
	topic_list := topics.Iter()
	for topic_list.Next() {
		key, leaf := topic_list.Key(), topic_list.Leaf()
		topic_name, err := topics.KeyBytesToStr(key)
		if err != nil {
			{
				fmt.Fprintf(os.Stderr, "Failed to convert to str", err)
				os.Exit(1)
			}
		}
		metricTemp, err := speed.NewPCPCounter(
			0,
			topic_name,
			"BW of topic",
		)
		leaf_val, err := topics.LeafBytesToStr(leaf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to convert to str", err)
			os.Exit(1)
		}
		leaf_int, err := strconv.ParseUint(leaf_val, 0, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to convert to int", err)
			os.Exit(1)
		}
		metricTemp.Set(int64(leaf_int))
	}
}
