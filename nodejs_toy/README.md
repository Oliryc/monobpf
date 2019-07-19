Simple nodejs programm to show illustrate usage of somme PCP/BCC/Vector tool.


## Slow, everytime

This is an illustration of [this blog post](https://www.future-processing.pl/blog/on-problems-with-threads-in-node-js/).

The first example illustrate the fact that in nodeJS, asynchronous call are done with a pool of thread, each blocking for the duration of the call.

Compile the shared library

```
make
```

Run an see how the first $UV_THREADPOOL_SIZE calls are done in about one second and how this delay all other calls.

```
UV_THREADPOOL_SIZE=4 LD_PRELOAD=./scandir.so node main.js
UV_THREADPOOL_SIZE=8 LD_PRELOAD=./scandir.so node main.js
UV_THREADPOOL_SIZE=16 LD_PRELOAD=./scandir.so node main.js
```

## Slow, randomly (Work In Progress)

Some API[^1] are based on websocket.

[^1]: https://cryptowat.ch/docs/websocket-api

This is a more complete example where the call to scandir is slow randomly. The purpose is to debug this with the BCC PMDA for PCP and Vector.

TODO
- List PCP module used
- Explain running client server

```
npm install websocket
```

### Vector

http://localhost:3000/?q=%5B%7B%22p%22:%22http%22,%22h%22:%22localhost:44323%22,%22hs%22:%22localhost%22,%22ci%22:%22_all%22,%22cl%22:%5B%22bcc-biolatency%22,%22bcc-gclatency%22,%22bcc-tcplife%22,%22bcc-http-analysis%22%5D%7D%5D

### To test

  ./go/bin/ws ws://10.186.108.12:7777
  > {"type": "echo", "payload": 1}

with https://github.com/hashrocket/ws
