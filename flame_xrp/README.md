# Test traffic generation and draw flamegraph

## Parameter of the script

The `client.js` script we use below takes the following parameters:

1. number of streams,
2. amount to send per stream,
3. amount to send per packet.

*Note*: total amount per session = number of streams * amount per stream

## Reference implementation (commit ff20aff6f064945ef4b616baaa06990c2592707a (master branch at the time, tag: v22.5.0))

### Hosts

- A client connector **c1**
- A server connector **s1**
- Another “observer” connector **o1**

*c1* and *o1* are both connected to *s1*.

### Overview

1. server.js script is continously running with *s1*.
2. *o1* is launched with 0x to generate flamegraph
2. *c1* is launched with 0x to generate flamegraph
2. Run script client.js (with various parameter, see next section) on *c1*, once *c1* is connected to everything. Use the time command to know how much time it requires.

**Note**: Try to keep *o1* and *s1* running about the same time

Command to run *c1* and *o1*:
```
nvm use 12.9.0 # Or another version
DEBUG='' 0x src/index.js
```

## Rafiki (commit 32a6979424de209a5cca182f738a65157c159862 (master branch at the time))

### Hosts

Same thing.

### Overview

Idem.

Command to run *c1* and *o1*:

```
nvm use 12.9.0 # Or another version
DEBUG='' 0x build/src/start.js
```

We use Rafiki with Redis.

# Results

Results are placed in a folder named from the connector implementation, containing:
- a folder with the approximate date, nodejs version and parameters used with the script along with debugging directive
  - a file `c1.html` (resp. `o1.html`), the flamegraph for *c1* (resp *o1*)
  - a file `duration` records the duration of the execution of the client script

*Note*: If `o1.html` is absent, it is that *o1* was not running during this particular benchmark

Like so:

```
rafiki
├── 2019-08-28-v10.16.3-2_5000_1-debug_off
│  ├── c1.html
│  ├── duration
│  └── o1.html
└── 2019-08-28-v10.16.3-2_50000_1-debug_off
   ├── c1.html
   ├── duration
   └── o1.html
reference
├── 2019-08-28-v10.16.3-2_5000_1-debug_off
│  ├── c1.html
│  ├── duration
│  └── o1.html
└── 2019-08-28-v10.16.3-2_50000_1-debug_off
   ├── c1.html
   ├── duration
   └── o1.html
```

