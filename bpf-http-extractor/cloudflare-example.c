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
  /*
  if (eth_hdr + 1 > (struct ethhdr *)data_end)        
    return XDP_ABORTED;
  */

  ip_hdr = (struct iphdr *)(eth_hdr + 1);
  if (ip_hdr + 1 > (struct iphdr *)data_end)
    return XDP_ABORTED;

  tcp_hdr = (struct tcphdr*)((u8 *)ip_hdr + ip_hdr->ihl * 4);
  if (tcp_hdr + 1 > (struct tcphdr *)data_end)
    return XDP_ABORTED;

  bpf_trace_printk("tcp_hdr->source: %u\n", htons(tcp_hdr->source));
  bpf_trace_printk("tcp_hdr->dest: %u\n", htons(tcp_hdr->dest));

  bpf_trace_printk("tcp_hdr->doff: %u\n", tcp_hdr->doff);

  tcp_opts = (u8 *)(tcp_hdr + 1);

  return XDP_PASS;
}

int xdp_prog1(struct CTXTYPE *ctx) {
    bpf_trace_printk("Hi\n");
    void* data_end = (void*)(long)ctx->data_end;
    void* data = (void*)(long)ctx->data;

    return match_p0f(data, data_end);
}
