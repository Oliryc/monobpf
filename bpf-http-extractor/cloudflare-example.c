// From https://www.netdevconf.org/2.1/slides/apr6/bertin_Netdev-XDP.pdf
#define KBUILD_MODNAME "foo"
#include <linux/bpf.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <net/ip.h>
#include <linux/tcp.h>

BPF_TABLE("array", int, long, dropcnt, 256);

static inline int match_p0f(void *data, void *data_end) {
  struct ethhdr *eth_hdr;
  struct iphdr  *ip_hdr;
  struct tcphdr *tcp_hdr;
  u8 *tcp_opts;
  eth_hdr = (struct ethhdr *)data;
  if (eth_hdr + 1 > (struct ethhdr *)data_end)        
    return XDP_ABORTED;
  if (!(eth_hdr->h_proto == htons(ETH_P_IP)))
    return XDP_PASS;
  ip_hdr = (struct iphdr *)(eth_hdr + 1);
  if (ip_hdr + 1 > (struct iphdr *)data_end)
    return XDP_ABORTED;
  // if (!((ip_hdr->frag_off & IP_MBZ) == 0))
  tcp_hdr = (struct tcphdr*)((u8 *)ip_hdr + ip_hdr->ihl * 4);

  u16 s1 = tcp_hdr->source;
  u16 s2 = tcp_hdr->source;
  bpf_trace_printk("tcp_hdr->source==: %p\n", &tcp_hdr->source);
  bpf_trace_printk("tcp_hdr->source: %p\n", &tcp_hdr->source);

  if (tcp_hdr + 1 > (struct tcphdr *)data_end)
    return XDP_ABORTED;
  if (!(tcp_hdr->dest == htons(1234)))
    return XDP_PASS;
  if (!(tcp_hdr->doff == 10))
    return XDP_PASS;
  if (!((htons(ip_hdr->tot_len) - (ip_hdr->ihl * 4) - (tcp_hdr->doff * 4)) == 0))
    return XDP_PASS;

  tcp_opts = (u8 *)(tcp_hdr + 1);
  if (tcp_opts + (tcp_hdr->doff - 5) * 4 > (u8 *)data_end)
    return XDP_ABORTED;
  if (!(htons(tcp_hdr->window) == htons(*(u16 *)(tcp_opts + 2)) * 0xa))
    return XDP_PASS;
  if (!(*(u8 *)(tcp_opts + 19) == 6))
    return XDP_PASS;
  if (!(tcp_opts[0] == 2))
    return XDP_PASS;
  if (!(tcp_opts[4] == 4))
    return XDP_PASS;
  if (!(tcp_opts[6] == 8))
    return XDP_PASS;
  if (!(tcp_opts[16] == 1))
    return XDP_PASS;
  if (!(tcp_opts[17] == 3))
    return XDP_PASS;

  return XDP_DROP;
}

int xdp_prog1(struct CTXTYPE *ctx) {
    bpf_trace_printk("Hi\n");
    void* data_end = (void*)(long)ctx->data_end;
    void* data = (void*)(long)ctx->data;

    return match_p0f(data, data_end);
}
