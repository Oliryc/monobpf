#!/usr/bin/python
# -*- coding: utf-8 -*-
#
# xdp_ip_whitelist.py Drop packet coming from ips not in a whitelist
#
# Copyright (c) 2019 SnT
# Based on https://github.com/iovisor/bcc/blob/master/examples/networking/xdp/xdp_drop_count.py,
# Copyright (c) 2016 PLUMgrid
# Copyright (c) 2016 Jan Ruth
# Licensed under the Apache License, Version 2.0 (the "License")

from bcc import BPF
import pyroute2
import time
import sys
import socket, struct

# Like blockedIp = ['10.244.3.24']
blockedIp = [
]
debug = 0
flags = 0

def usage():
    print("Usage: {0} [-S] <ifdev>".format(sys.argv[0]))
    print("       -S: use skb mode\n")
    print("e.g.: {0} eth0\n".format(sys.argv[0]))
    exit(1)

if len(sys.argv) < 2 or len(sys.argv) > 3:
    usage()

if len(sys.argv) == 2:
    device = sys.argv[1]

if len(sys.argv) == 3:
    if "-S" in sys.argv:
        # XDP_FLAGS_SKB_MODE
        flags |= 2 << 0

    if "-S" == sys.argv[1]:
        device = sys.argv[2]
    else:
        device = sys.argv[1]

mode = BPF.XDP
#mode = BPF.SCHED_CLS

if mode == BPF.XDP:
    ret = "XDP_DROP"
    ctxtype = "xdp_md"
else:
    ret = "TC_ACT_SHOT"
    ctxtype = "__sk_buff"

# load BPF program
bpf_src = ''
with open("xdp_ip_whitelist.bpf") as bpf_file:
    bpf_src = bpf_file.read()
    ip4array = map(str,
        [socket.htonl(struct.unpack("!L", socket.inet_aton(ip))[0])
         for ip in blockedIp])
    bpf_src = bpf_src.replace("__IP4ARRAY__", ", ".join(ip4array))
    bpf_src = bpf_src.replace("__IP4ARRAYSIZE__", str(len(ip4array)))
    if debug:
        print("C code of BPF program:")
        print(bpf_src)
b = BPF(text = bpf_src,
        cflags=["-w", "-DRETURNCODE=%s" % ret, "-DCTXTYPE=%s" % ctxtype])

fn = b.load_func("xdp_prog1", mode)

if mode == BPF.XDP:
    print("XDP Mode")
    b.attach_xdp(device, fn, flags)
else:
    print("TC Fallback")
    ip = pyroute2.IPRoute()
    ipdb = pyroute2.IPDB(nl=ip)
    idx = ipdb.interfaces[device].index
    ip.tc("add", "clsact", idx)
    ip.tc("add-filter", "bpf", idx, ":1", fd=fn.fd, name=fn.name,
          parent="ffff:fff2", classid=1, direct_action=True)

dropcnt = b.get_table("dropcnt")
prev = [0] * 256
print("Accepting packets only from the following IP addresses {}, hit CTRL+C to stop".format(blockedIp))
while 1:
    try:
        time.sleep(1)
    except KeyboardInterrupt:
        print("Removing filter from device")
        break;

if mode == BPF.XDP:
    b.remove_xdp(device, flags)
else:
    ip.tc("del", "clsact", idx)
    ipdb.release()
