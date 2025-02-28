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
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	bpf "github.com/iovisor/gobpf/bcc"
)

var (
	device            = "lo"
	snapshotLen int32 = 1024
	promiscuous       = false
	_           error
	timeout     = 30 * time.Second
	handle      *pcap.Handle
)

func main() {
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
		}
	}()

	headers := bpf.NewTable(module.TableId("headers"), module)
	header_list := headers.Iter()
	head_size := bpf.NewTable(module.TableId("head_size"), module)
	hs_list := head_size.Iter()

	// Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println("Packet!")
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		tcp, _ := tcpLayer.(*layers.TCP)
		for header_list.Next() {
			key, leaf := header_list.Key(), header_list.Leaf()
			leaf_val, err := headers.LeafBytesToStr(leaf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to convert leaf_val to str: %v", err)
				os.Exit(1)
			}
			if fmt.Sprint(tcp.Seq) == leaf_val {
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
				headers.Set(key, new_leaf)
				// Should only be one
				for hs_list.Next() {
					key, leaf := hs_list.Key(), hs_list.Leaf()
					leaf_str, _ := headers.LeafBytesToStr(leaf)
					leaf_int, _ := strconv.Atoi(leaf_str)
					leaf_int = leaf_int - 1
					new_leaf, _ := headers.LeafStrToBytes(fmt.Sprint(leaf_int))
					head_size.Set(key, new_leaf)
				}
			}
		}
	}
}
