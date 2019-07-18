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
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	bpf "github.com/iovisor/gobpf/bcc"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	device            = "lo"
	snapshotLen int32 = 10240
	promiscuous       = false
	timeout           = -1 * time.Second
	handle      *pcap.Handle
	useOne      = false
)

func main() {
	if useOne {
		packetTrace1()
	} else {
		packetTrace2()
	}
}

type chownEvent struct {
	SeqNum      uint64
	SrcIP       uint32
	DstIP       uint32
	ReturnValue int32
	Filename    [256]byte
}

func packetTrace1() {
	filesrc, err := ioutil.ReadFile("packet_trace.bpf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load xdp source %v\n", err)
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
	fn, err := module.Load("xdp_prog1", C.BPF_PROG_TYPE_XDP, 1, 65536)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load xdp prog: %v\n", err)
		os.Exit(1)
	}
	err = module.AttachXDP(device, fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach xdp prog: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Attached BPF programm to '%v'\n", device)

	defer func() {
		if err := module.RemoveXDP(device); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove XDP from %s: %v\n", device, err)
			os.Exit(1)
		}
	}()
	headers := bpf.NewTable(module.TableId("headers"), module)
	headSize := bpf.NewTable(module.TableId("head_size"), module)
	// Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println("Packet!")
		headerList := headers.Iter()
		hsList := headSize.Iter()
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp, _ := tcpLayer.(*layers.TCP)
		fmt.Fprintf(os.Stdout, "New packet! %d \n", tcp.Seq)
		for headerList.Next() {
			key, leaf := headerList.Key(), headerList.Leaf()
			leafVal, err := headers.LeafBytesToStr(leaf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to convert leaf_val to str: %v", err)
				os.Exit(1)
			}
			leafInt, err := strconv.ParseUint(leafVal, 0, 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to convert to int %d\n", err)
				os.Exit(1)
			}
			// fmt.Fprintf(os.Stdout, "seq: ", tcp.Seq, ", leaf_val: ", leaf_int, "\n")
			if uint32(leafInt) == tcp.Seq {
				fmt.Fprintf(os.Stdout, "seq: %d\n", tcp.Seq)
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
				newLeaf, _ := headers.LeafStrToBytes("0")
				err = headers.Set(key, newLeaf)
				if err != nil {
					fmt.Fprint(os.Stderr, "Failed to update leaf", err)
					os.Exit(1)
				}
				// Should only be one
				for hsList.Next() {
					key, leaf := hsList.Key(), hsList.Leaf()
					leafStr, _ := headers.LeafBytesToStr(leaf)
					leafInt, _ := strconv.Atoi(leafStr)
					leafInt = leafInt - 1
					newLeaf, _ := headers.LeafStrToBytes(fmt.Sprint(leafInt))
					err = headSize.Set(key, newLeaf)
					fmt.Fprint(os.Stderr, "Failed to update leaf", err)
				}
			}
		}
	}
}
func packetTrace2() {
	filesrc, err := ioutil.ReadFile("packet_trace2.bpf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load xdp source %v\n", err)
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
	fn, err := module.Load("xdp_prog1", C.BPF_PROG_TYPE_XDP, 1, 65536)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load xdp prog: %v\n", err)
		os.Exit(1)
	}
	err = module.AttachXDP(device, fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach xdp prog: %v\n", err)
		os.Exit(1)
	}
	table := bpf.NewTable(module.TableId("packet_events"), module)

	channel := make(chan []byte)

	perfMap, err := bpf.InitPerfMap(table, channel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init perf map: %s\n", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	defer func() {
		if err := module.RemoveXDP(device); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove XDP from %s: %v\n", device, err)
			os.Exit(1)
		}
	}()
	go func() {
		var event chownEvent
		for {
			data := <-channel
			err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &event)
			if err != nil {
				fmt.Printf("failed to decode received data: %s\n", err)
				continue
			}
			fmt.Printf("Got Packet with: SeqNum %d SrcIP %d DstIP %d\n",
				event.SeqNum, event.SrcIP, event.DstIP)
		}
	}()

	perfMap.Start()
	<-sig
	perfMap.Stop()

}
