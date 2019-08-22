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
	"time"
)

func secureRos(stopChan chan struct{}) {
	var (
		device = "lo"
	)
	filesrc, err := ioutil.ReadFile("secure_ros.bpf")
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
	fn, err := module.Load("secure_ros", C.BPF_PROG_TYPE_XDP, 1, 65536)
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
	for {
		select {
		default:
			fmt.Printf("Filtering online")
			time.Sleep(time.Second)

		case <-stopChan:
			return
		}
	}
}
