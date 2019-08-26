# Monobpf

monitoring tools

~~Cf [ebpf.md](./ebpf.md) for background & notes about eBPF~~ (outdated, most recent version on overleaf)

## :fire: [flame_xrp](flame_xrp): Flamegraphs generated on the connectors

This folder contains the protocol and the resulting flamegraphs generated on the connector. See the [Readme there](./flame_xrp/README.md)



## [bpf-http-extractor](bpf-http-extractor)

This folder contains the code attempting to parse the HTTP headers & method names.

It’s only getting short paquets without any HTTP payload.

## exporter_*.yml file

[ebpf_exporter][ebpf_exporter] configuration file

[ebpf_exporter]: https://github.com/cloudflare/ebpf_exporter/

## bcc_tuto

Some file created while following the bcc tutorial

## ebpf-sample-server

Apache server, to create a service to be monitored

## idp-connector docker

Some configuration and instruction for a ilp connector

## Monitoring connector

Attempt to monitor another container from another container that would have injected eBPF bytecode.

## sample_node.js

Attempt to create an ILP connector inside docker (to then be monitored)
