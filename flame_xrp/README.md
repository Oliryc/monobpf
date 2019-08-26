# Test traffic generation and draw flamegraph

## Parameter of the script

The `client.js` script we use below takes the following parameters:

1. number of streams,
2. amount to send per stream,
3. amount to send per packet.

*Note*: total amount per session = number of streams * amount per stream

## Reference implementation

### Hosts

- A client connector **c1**
- A server connector **s1**
- Another server connector **o1**

### Overview

1. server.js script is continously running with *s1*.
2. *o1* is launched with 0x to generate flamegraph
2. *c1* is launched with 0x to generate flamegraph
2. Run script client.js (with various parameter, see next section) on *c1*, once *c1* is connected to everything

**Note**: Try to keep *o1* and *s1* running about the same time

### Results

Results are placed in the folder `reference`, with:
- a folder with the approximate date and parameters used with the script
  - a file c1.html (resp. o1.html), the flamegraph for *c1* (resp *o1*)

## Rafiki

### Hosts

Same thing, without **o1**

### Overview

Idem, without *o1*

## Results

Idem, in folder `Rafiki`

