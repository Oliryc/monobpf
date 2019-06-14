#!/usr/bin/python
#
# nodejs_gc_tracing     Time GC calls
#                       For Linux, uses BCC, BPF. Embedded C.
#
# USAGE: nodejs_gc_tracing PID
#
# Copyright 2016 Netflix, Inc.
# Licensed under the Apache License, Version 2.0 (the "License")

# Based on https://github.com/iovisor/bcc/blob/master/examples/tracing/nodejs_http_server.py

from __future__ import print_function
from bcc import BPF, USDT
from bcc.utils import printb
import sys

if len(sys.argv) < 2:
    print("USAGE: {} PID".format(sys.argv[0]))
    exit()
pid = sys.argv[1]
debug = 1

# load BPF program
bpf_text = """
#include <uapi/linux/ptrace.h>
int do_trace(struct pt_regs *ctx) {
    uint64_t addr;
    char path[128]={0};
    bpf_usdt_readarg(2, ctx, &addr);
    bpf_probe_read(&path, sizeof(path), (void *)addr);
    bpf_trace_printk("path:%s\\n", path);
    return 0;
};
"""

# enable USDT probe from given PID
u = USDT(pid=int(pid))
u.enable_probe(probe="gc__start", fn_name="do_trace")
if debug:
    print(u.get_text())
    print(bpf_text)

# initialize BPF
b = BPF(text=bpf_text, usdt_contexts=[u])

# header
print("%-18s %-16s %-6s %s" % ("TIME(s)", "COMM", "PID", "ARGS"))

# format output
while 1:
    try:
        (task, pid, cpu, flags, ts, msg) = b.trace_fields()
    except ValueError:
        print("value error")
        continue
    except KeyboardInterrupt:
        exit()
    printb(b"%-18.9f %-16s %-6d %s" % (ts, task, pid, msg))
