# Checking size of XDP-based tcpdump-like tool

## Files

- `xdp_sample_pkts_*`: tcpdump-like tool. Completed version of https://github.com/xdp-project/xdp-tutorial/tree/master/tracing04-xdp-tcpdump. Place these files in tracing04-xdp-tcpdump directory to compile.

- `send_udp.py`: sends udp packet of variable size. The code is very short and self documenting

- `samples.pcap`: file produced by `xdp_sample_pkts_user`, when performing an HTTP request and sending UDP packet of up to ≃5100 byte. Can be opened with Wireshark.

## Sample commands

```
./send_udp.py 5000 6000
sudo ./xdp_sample_pkts_user -d lo -S -F
wireshark samples.pcap
```
