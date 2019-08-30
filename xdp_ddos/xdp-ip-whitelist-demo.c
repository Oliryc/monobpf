// xdp-ip-whitelist.py Drop packet coming from ips not in a whitelist
//
// Based on https://github.com/iovisor/bcc/blob/master/examples/
// networking/xdp/xdp_drop_count.py,
// Copyright (c) 2016 PLUMgrid
// Copyright (c) 2016 Jan Ruth
// Licensed under the Apache License, Version 2.0 (the "License")
// See http://apache.org/licenses/LICENSE-2.0
#define KBUILD_MODNAME "xdp_ddos"
#include <uapi/linux/bpf.h>
#include <linux/in.h>
#include <linux/if_ether.h>
#include <linux/if_packet.h>
#include <linux/if_vlan.h>
#include <linux/ip.h>
#include <linux/ipv6.h>

// IP not recognised in header
#define NOIP 0

#define WHITE4SIZE 6

// IP whitelist [ "192.168.1.187", "192.168.1.98", "192.168.1.1",
//                "192.168.1.150", "192.168.1.131", "192.168.1.68" ]  
static int ip4white[] = { 3137448128, 1644275904, 16885952,
                          2516691136, 2197924032, 1140959424};

/* Returns IPv4 address, in network byte order */
static inline int get_ipv4(void *data, u64 nh_off, void *data_end) {
    struct iphdr *iph = data + nh_off;

    if ((void*)&iph[1] > data_end)
        return NOIP;
    return iph->saddr;
}

int xdp_prog1(struct CTXTYPE *ctx) {

    void* data_end = (void*)(long)ctx->data_end;
    void* data = (void*)(long)ctx->data;

    struct ethhdr *eth = data;

    // Drop packets by default
    int rc = XDP_DROP;
    long *value;
    uint16_t h_proto;
    uint64_t nh_off = 0;
    uint32_t index;

    nh_off = sizeof(*eth);

    if (data + nh_off  > data_end)
        return rc;

    h_proto = eth->h_proto;

    if (h_proto == htons(ETH_P_IP)) {
        // Allow packet to pass if its IP is in the whitelist
        int ip = get_ipv4(data, nh_off, data_end);

        if (ip == NOIP) {
          return XDP_DROP;
        }
        #pragma unroll
        for (int i = 0; i < WHITE4SIZE; i++) {
          if (ip4white[i] == ip) {
            return XDP_PASS;
          }
        }
        return XDP_DROP;
    } else if (h_proto == htons(ETH_P_ARP)) {
      return XDP_PASS;
    } else {
      return XDP_DROP;
    }
    return XDP_DROP;
}
