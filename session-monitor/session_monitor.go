// session_monitor.go Flag when
// protocol type. Based on:
// https://github.com/iovisor/bcc/blob/master/examples/networking/xdp/xdp_drop_count.py
//
// Copyright (c) 2017 GustavoKatel
// Licensed under the Apache License, Version 2.0 (the "License")

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	bpf "github.com/iovisor/gobpf/bcc"
)

/*
#cgo CFLAGS: -I/usr/include/bcc/compat
#cgo LDFLAGS: -lbcc
#include <bcc/bcc_common.h>
#include <bcc/libbpf.h>
void perf_reader_free(void *ptr);
*/
import "C"

func usage() {
	fmt.Printf("Usage: %v <ifdev>\n", os.Args[0])
	fmt.Printf("e.g.: %v eth0\n", os.Args[0])
	os.Exit(1)
}

func main() {
	var device string
	filesrc, err := ioutil.ReadFile("session_monitor.bpf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load xdp source %v\n", err)
	}
	source := string(filesrc)

	if len(os.Args) != 2 {
		usage()
	}

	device = os.Args[1]

	ret := "XDP_DROP"
	ctxtype := "xdp_md"

	module := bpf.NewModule(source, []string{
		"-w",
		"-DRETURNCODE=" + ret,
		"-DCTXTYPE=" + ctxtype,
	})
	defer module.Close()

	fn, err := module.Load("session_monitor", C.SOCKET_FILTER, 1, 65536)
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
		}
	}()

	fmt.Println("Dropping packets, hit CTRL+C to stop")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	dropcnt := bpf.NewTable(module.TableId("dropcnt"), module)

	<-sig

	fmt.Printf("\n{IP protocol-number}: {total dropped pkts}\n")
	for it := dropcnt.Iter(); it.Next(); {
		key := bpf.GetHostByteOrder().Uint32(it.Key())
		value := bpf.GetHostByteOrder().Uint64(it.Leaf())

		if value > 0 {
			fmt.Printf("%v: %v pkts\n", key, value)
		}
	}
}
