# Recherches sur eBPF

## Objectifs

Monitoring de processus dans des containers Docker

1. Lister toutes les informations accessibles via eBPF
  - monitoring non invasif
  - implications de sécurité
    - voir quel appels bloquer

2. Dans un second temps, afficher des graphes ([Prometheus](https://prometheus.io/))

3. Voir l’intégration avec Kubernetes

## Pistes

- [ebpf_exporter][3] pourrait être utilisé (cf [exemples][ebpf_exporter_example]) (mais à
  configurer, sachant que la configuration nécessite d’écrire un peu de code C
  pour eBPF)
  - TODO Déterminer les événements pertinents pour le monitoring à partir de
    [cette liste][https://github.com/iovisor/bcc/blob/master/docs/reference_guide.md#events--arguments] et des suivantes dans la page
  - TODO Essayer de monitorer un événement dans un autre processus/container à partir d’un [exemple][seccomp-bpf]


## Informations accessibles via eBPF et outils dont on pourrait s’inspirer

Sélectionnés dans [cette liste][bcc-tools]

### Disque

- tools/[biosnoop](tools/biosnoop.py): Trace block device I/O with PID and latency. [Examples](tools/biosnoop_example.txt).
- tools/[filetop](tools/filetop.py): File reads and writes by filename and process. Top for files. [Examples](tools/filetop_example.txt).

### TCP

- tools/[tcpconnlat](tools/tcpconnlat.py): Trace TCP active connection latency (connect()). [Examples](tools/tcpconnlat_example.txt).
- tools/[tcptop](tools/tcptop.py): Summarize TCP send/recv throughput by host. Top for TCP. [Examples](tools/tcptop_example.txt).
- tools/[tcplife](tools/tcplife.py): Trace TCP sessions and summarize lifespan. [Examples](tools/tcplife_example.txt).
- tools/[tcptracer](tools/tcptracer.py): Trace TCP established connections (connect(), accept(), close()). [Examples](tools/tcptracer_example.txt).

### Exécution & CPU

- tools/[cachetop](tools/cachetop.py): Trace page cache hit/miss ratio by processes. [Examples](tools/cachetop_example.txt).
- tools/[cpudist](tools/cpudist.py): Summarize on- and off-CPU time per task as a histogram. [Examples](tools/cpudist_example.txt)
- tools/[exitsnoop](tools/exitsnoop.py): Trace process termination (exit and fatal signals). [Examples](tools/exitsnoop_example.txt).

- tools/[funclatency](tools/funclatency.py): Time functions and show their latency distribution. [Examples](tools/funclatency_example.txt).
- tools/[funcslower](tools/funcslower.py): Trace slow kernel or user function calls. [Examples](tools/funcslower_example.txt).

- tools/[opensnoop](tools/opensnoop.py): Trace open() syscalls. [Examples](tools/opensnoop_example.txt).
- tools/[pidpersec](tools/pidpersec.py): Count new processes (via fork). [Examples](tools/pidpersec_example.txt).
- tools/[profile](tools/profile.py): Profile CPU usage by sampling stack traces at a timed interval. [Examples](tools/profile_example.txt).
- tools/[syscount](tools/syscount.py): Summarize syscall counts and latencies. [Examples](tools/syscount_example.txt).
- tools/[trace](tools/trace.py): Trace arbitrary functions, with filters. [Examples](tools/trace_example.txt).
- tools/[ucalls](tools/lib/ucalls.py): Summarize method calls or Linux syscalls in high-level languages. [Examples](tools/lib/ucalls_example.txt).
- tools/[uflow](tools/lib/uflow.py): Print a method flow graph in high-level languages. [Examples](tools/lib/uflow_example.txt).
- tools/[ugc](tools/lib/ugc.py): Trace garbage collection events in high-level languages. [Examples](tools/lib/ugc_example.txt).
- tools/[ustat](tools/lib/ustat.py): Collect events such as GCs, thread creations, object allocations, exceptions and more in high-level languages. [Examples](tools/lib/ustat_example.txt).
- tools/[uthreads](tools/lib/uthreads.py): Trace thread creation events in Java and raw pthreads. [Examples](tools/lib/uthreads_example.txt).

### Illustration des implications de sécurité

- tools/[bashreadline](tools/bashreadline.py): Print entered bash commands system wide. [Examples](tools/bashreadline_example.txt).
- tools/[sslsniff](tools/sslsniff.py): Sniff OpenSSL written and readed data. [Examples](tools/sslsniff_example.txt).
- tools/[ttysnoop](tools/ttysnoop.py): Watch live output from a tty or pts device. [Examples](tools/ttysnoop_example.txt).

## Ressources

### https://medium.com/@andrewhowdencom/adventures-with-ebpf-and-prometheus-6a59dd170b26

- Prometheus <-> Exporter <-> eBPF
- Nombreux exporter existants, peut-être que
  [ebpf_exporter][3] fait une partie de ce qu’on veut
  

### [Restriction des appels systèmes][seccomp-bpf]

- exemple de code  -> point de départ
- restriction des appels système (déjà utilisé dans une certaine mesure dans Docker d’après [eBPF & Kubernetes][5])

### https://blog.cloudflare.com/bpf-the-forgotten-bytecode/

- Limitations du bytecode BPF : 4096 instructions
  - eBPF est aussi affecté par ces limitations

### [eBPF & Kubernetes][5]

- eBPF sert surtout pour du monitoring réseau dans kubernetes, mais d’autres
  applications sont envisagées (TODO Les comprendre)

### https://medium.com/@andrewhowdencom/coming-to-grips-with-ebpf-4a5434591167

TODO

[seccomp-bpf]: https://blog.yadutaf.fr/2014/05/29/introduction-to-seccomp-bpf-linux-syscall-filter/
[3]: https://github.com/cloudflare/ebpf_exporter
[ebpf_exporter_example]: https://github.com/cloudflare/ebpf_exporter#examples
[5]: https://kubernetes.io/blog/2017/12/using-ebpf-in-kubernetes/
[bcc-tools]: https://github.com/iovisor/bcc/blob/master/README.md#tools
