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
	"fmt"
	bpf "github.com/iovisor/gobpf/bcc"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func getTopicsDemo() (error, []string) {
	app := "/opt/ros/melodic/bin/rostopic"
	arg0 := "list"
	cmd := exec.Command(app, arg0)
	stdout, err := cmd.Output()
	stdstring := string(stdout)
	topics := strings.Split(stdstring, "\n")
	if err != nil {
		println(err.Error())
		return nil, topics
	}
	return nil, topics
}

func MonitorROS(muTopics *sync.Mutex, topicList []string, stopChan chan struct{}) {
	var (
		device = "lo"
	)
	filesrc, err := ioutil.ReadFile("ros_metric.bpf")
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "Failed to load xdp source %v\n", err)
		os.Exit(1)
	}
	source := string(filesrc)
	ret := "XDP_PASS"
	ctxtype := "xdp_md"
	module := bpf.NewModule(source, []string{
		"-w",
		"-DRETURNCODE=" + ret,
		"-DCTXTYPE=" + ctxtype,
	})
	defer module.Close()
	fn, err := module.Load("session_monitor", C.BPF_PROG_TYPE_XDP, 1, 65536)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "Failed to load xdp prog: %v\n", err)
		os.Exit(1)
	}
	err = module.AttachXDP(device, fn)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "Failed to attach xdp prog: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := module.RemoveXDP(device); err != nil {
			_, err = fmt.Fprintf(os.Stderr, "Failed to remove XDP from %s: %v\n", device, err)
			os.Exit(1)
		}
	}()
	fmt.Println("May be dropping packets, hit CTRL+C to stop. " +
		"See output of `sudo cat /sys/kernel/debug/tracing/trace_pipe`")
	session := bpf.NewTable(module.TableId("session"), module)
	for {
		select {
		default:
			err, topics := getTopicsDemo()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to get Topics list %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("topics %v\n", topics)
			muTopics.Lock()
			copy(topicList, topics)
			muTopics.Unlock()
			time.Sleep(time.Second)
			for it := session.Iter(); it.Next(); {
				key, leaf := it.Key(), it.Leaf()
				key_str, err := session.KeyBytesToStr(key)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Failed to convert to str", err)
					os.Exit(1)
				}
				leaf_str, err := session.LeafBytesToStr(leaf)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Failed to convert to str", err)
					os.Exit(1)
				}
				fmt.Printf("%s", key_str)
				fmt.Printf("%s", leaf_str)
			}
		case <-stopChan:
			return
		}
	}
	/*headers := bpf.NewTable(module.TableId("headers"), module)
	head_size := bpf.NewTable(module.TableId("head_size"), module)
	topics = bpf.NewTable(module.TableId("topics"), module)
	//Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		header_list := headers.Iter()
		hs_list := head_size.Iter()
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp, _ := tcpLayer.(*layers.TCP)
		fmt.Fprintf(os.Stdout, "New packet!", tcp.Seq, "\n")
		for header_list.Next() {
			key, leaf := header_list.Key(), header_list.Leaf()
			leaf_val, err := headers.LeafBytesToStr(leaf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to convert to str", err)
				os.Exit(1)
			}
			leaf_int, err := strconv.ParseUint(leaf_val, 0, 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to convert to int", err)
				os.Exit(1)
			}
			// fmt.Fprintf(os.Stdout, "seq: ", tcp.Seq, ", leaf_val: ", leaf_int, "\n")
			if uint32(leaf_int) == tcp.Seq {
				fmt.Fprintf(os.Stdout, "seq: ", tcp.Seq)
				applicationLayer := packet.ApplicationLayer()
				if applicationLayer != nil {
					fmt.Println("Application layer/Payload found.")
					fmt.Printf("%s\n", applicationLayer.Payload())

					// Search for a string inside the payload
					if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
						fmt.Println("HTTP found!")
					}
				}

				// Check for errors
				if err := packet.ErrorLayer(); err != nil {
					fmt.Println("Error decoding some part of the packet:", err)
				}
				new_leaf, _ := headers.LeafStrToBytes("0")
				err = headers.Set(key, new_leaf)
				if err != nil {
					fmt.Fprint(os.Stderr, "Failed to update leaf", err)
					os.Exit(1)
				}
				// Should only be one
				for hs_list.Next() {
					key, leaf := hs_list.Key(), hs_list.Leaf()
					leaf_str, _ := headers.LeafBytesToStr(leaf)
					leaf_int, _ := strconv.Atoi(leaf_str)
					leaf_int = leaf_int - 1
					new_leaf, _ := headers.LeafStrToBytes(fmt.Sprint(leaf_int))
					err = head_size.Set(key, new_leaf)
					fmt.Fprint(os.Stderr, "Failed to update leaf", err)
				}
			}
		}
	}*/
}
