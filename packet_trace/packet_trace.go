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
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	bpf "github.com/iovisor/gobpf/bcc"
	"io/ioutil"
	"log"
	"os"
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

	defer func() {
		if err := module.RemoveXDP(device); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove XDP from %s: %v\n", device, err)
			os.Exit(1)
		}
	}()

	headers := bpf.NewTable(module.TableId("headers"), module)
	head_size := bpf.NewTable(module.TableId("head_size"), module)

	// Open device
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
		fmt.Fprintf(os.Stdout, "New packet!",tcp.Seq,"\n")
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
	}
}
