// From https://www.netdevconf.org/2.1/slides/apr6/bertin_Netdev-XDP.pdf
#define KBUILD_MODNAME "foo"
#include <linux/bpf.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <net/ip.h>
#include <linux/tcp.h>

#define LISTEN_PORT 80

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
  if (ip_hdr + 1 > (struct iphdr *)data_end) {
    bpf_trace_printk("=== Abort 2\n");
    return XDP_ABORTED;
  }

  tcp_hdr = (struct tcphdr*)((u8 *)ip_hdr + ip_hdr->ihl * 4);
  if (tcp_hdr + 1 > (struct tcphdr *)data_end) {
    bpf_trace_printk("=== Abort 3\n");
    return XDP_ABORTED;
  }

  u16 sport = htons(tcp_hdr->source);
  u16 dport = htons(tcp_hdr->dest);
  // u16 sport = tcp_hdr->source;
  // u16 dport = tcp_hdr->dest;
  bpf_trace_printk("tcp_hdr->doff: %u\n", tcp_hdr->doff);

  if (sport == LISTEN_PORT || dport == LISTEN_PORT) {
    char *http_hdr = (char *)((u8 *)tcp_hdr + tcp_hdr->doff * 4);
    bpf_trace_printk("==============================, %u\n", sizeof(*http_hdr));
    if (http_hdr + 2 > (char *)data_end) {
      bpf_trace_printk("=== Abort 4\n");
      return XDP_ABORTED;
    }

    bpf_trace_printk("tcp_hdr[1]: %u\n", http_hdr[1]);
    // bpf_trace_printk("tcp_hdr[2]: %u\n", http_hdr[2]);
    // bpf_trace_printk("tcp_hdr[3]: %u\n", http_hdr[3]);
    // bpf_trace_printk("tcp_hdr[4]: %u\n", http_hdr[4]);
    // bpf_trace_printk("tcp_hdr[5]: %u\n", http_hdr[5]);
  } else {
    bpf_trace_printk("Ignored tcp_hdr->source: %u, tcp_hdr->dest: %u\n", sport, dport);
  }

  bpf_trace_printk("=== Pass 5\n");
  return XDP_PASS;
}

int xdp_prog1(struct CTXTYPE *ctx) {
    bpf_trace_printk("Hi\n");
    void* data_end = (void*)(long)ctx->data_end;
    void* data = (void*)(long)ctx->data;

    int rc = match_p0f(data, data_end);
    bpf_trace_printk("Bye\n");
    return rc;
}
