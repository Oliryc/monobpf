# Recherches sur eBPF

## Objectifs

Monitoring de processus dans des containers Docker

1. Lister toutes les informations accessibles via eBPF
  - monitoring non invasif
  - implications de sécurité
    - voir quel appels bloquer

2. Dans un second temps, afficher des graphes ([Prometheus](https://prometheus.io/))

3. Voir l’intégration avec Kubernetes

### Suite aux retours

- Collecter un maximum niveau réseau (type de paquet) puis diskio
- Voir quelle granularité on peut avoir (jusqu’au module nodeJS ?)

## Pistes

- [ebpf_exporter][ebpf_exporter] pourrait être utilisé (cf [exemples][ebpf_exporter_example]) (mais à
  configurer, sachant que la configuration nécessite d’écrire un peu de code C
  pour eBPF)
  - Essayer de monitorer un événement dans un autre processus/container à partir d’un [exemple][seccomp-bpf] : existe déjà dans les scripts du dossier tools de bcc
  - Installer et faire fonctionner [ebpf_exporter][ebpf_exporter].
- Métriques à collecter :
  - TODO Nœud nodeJS : https://github.com/iovisor/bcc/blob/master/tools/tplist.py
  - TODO Suivi du GC Java
  - TODO Outil stackcount
- TODO Regarder les exemples :
  - TODO https://github.com/cloudflare/ebpf_exporter/blob/master/examples/eaddrinuse.yaml -> détection d’attaque (parce qu’à court de socket libre?)
  - TODO https://github.com/cloudflare/ebpf_exporter/blob/master/examples/runqlat.yaml mesure de latence
  - TODO https://github.com/cloudflare/ebpf_exporter/blob/master/examples/ipcstat.yaml mesure d’utilisation du CPU plus fidèle
    - TODO Voir si on peut la restreindre à un processus

- TODO Outils pour explorer le contenu des paquets : [xdpcap][xdpcap], cf [contexte][xdpcap-context]

- TODO Faire un schéma récapitulatif des outils utilisés (container docker, prometheus, hôte…)

- TODO Regarder [les exemples du dépôt bcc][bcc-example]

## Informations accessibles via eBPF et outils dont on pourrait s’inspirer

Sélectionnés dans [cette liste][bcc-tools]

### Disque

- tools/[biosnoop](tools/biosnoop.py): Trace block device I/O with PID and latency. [Examples](tools/biosnoop_example.txt).
- tools/[filetop](tools/filetop.py): File reads and writes by filename and process. Top for files. [Examples](tools/filetop_example.txt).

### TCP : on accède à une structure de socket et on peut se hooker sur les principaux appels systèmes pour suivre la connexion

* tools/[tcpconnlat](tools/tcpconnlat.py): Trace TCP active connection latency (connect()). [Examples](tools/tcpconnlat_example.txt).
* tools/[tcptop](tools/tcptop.py): Summarize TCP send/recv throughput by host. Top for TCP. [Examples](tools/tcptop_example.txt).
* tools/[tcplife](tools/tcplife.py): Trace TCP sessions and summarize lifespan. [Examples](tools/tcplife_example.txt).
* tools/[tcptracer](tools/tcptracer.py): Trace TCP established connections (connect(), accept(), close()). [Examples](tools/tcptracer_example.txt).

### Exécution & CPU

- tools/[cachetop](tools/cachetop.py): Trace page cache hit/miss ratio by processes. [Examples](tools/cachetop_example.txt).
- tools/[cpudist](tools/cpudist.py): Summarize on- and off-CPU time per task as a histogram. [Examples](tools/cpudist_example.txt)
- tools/[exitsnoop](tools/exitsnoop.py): Trace process termination (exit and fatal signals). [Examples](tools/exitsnoop_example.txt).

- tools/[funclatency](tools/funclatency.py): Time functions and show their latency distribution. [Examples](tools/funclatency_example.txt).
  - Plusieurs problèmes ce sont posés avec l’outil funclatency : je n’ai pas encore trouvé comment attacher des fonctions arbitraire d’un programme. Par contre, on peut s’attacher aux fonctions de la bibliothèque standard C, ça fonctionne bien ça. L’autre problème est que la documentation indique que les résultats peuvent êtres erronés pour les fonctions récursives et imbriquées ainsi que dans les cas de multithreding (pour cause d’état partagé).
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

- Accès aux argument et au type de retour des fonctions de la bibliothèque C standard

- L’exécution de bytecode eBPF nécessite des privilèges plus élevés (en dehors
  des cas simples où on considère qu’il n’y a pas de risque de sécurité) : on
  peut quand même désactiver cette possibilité avec `sysctl
  kernel.unprivileged_bpf_disabled`. Cf aussi [ce commit](https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=1be7f75d1668d6296b80bf35dcf6762393530afc)

- Les appels peuvent être [restreints][seccomp-bpf] assez finemenet par le processus père.

## Terminologie

- BPF_TABLE : macro bas niveau, utilisée plutôt par les macros suivantes
- BPF_HASH : tableau associatif
- BPF_ARRAY : tableau, accès et mise à jour rapide
- BPF_HISTOGRAMME : réalisation d’histogrammes avec la méthode increment pour modifier l’effectif d’une catégorie

## [ebpf_exporter][ebpf_exporter]

### Installation

- Le mieux est d’utiliser le ppa fournit avec [les instructions d’installation][bcc-install]
- Nécessite Docker et demande libbcc en version 0.9, mais la version 0.10 est sortie depuis… attention, si les versions ne concordent pas, on a une segfault. Appliquer le patch ./ebpf_exporter.patch
- Problème pour l’ajout de clé aussi, bloqué par le pare feu semble-t-il. Utiliser hkp://keyserver….com:80

### Utilisation

- Avec l’option --debug, il est possible de prévisualiser les données à l’adresse http://localhost:9435/tables 

- TODO Faire une PR pour soumettre ces changements ?

### Avec prometheus

- Voir le fichier ./prometheus.yml
``` bash
$ sudo ./release/ebpf_exporter-1.1.0-13-ge548839-dev/ebpf_exporter --config.file ./examples/shrinklat.yaml
$ ~/Downloads/prometheus-2.10.0.linux-amd64/prometheus --config.file=./ebpf/prometheus.yml
```

- Les graphiques sont à l’adresse http://localhost:9090/graph

### Limites

- L’export des données est limité :
  - Les BPF_HASH doivent avoir des valeurs de type u64 (alors que le type des
    clés est libre), contrainte qu’on n’a pas en utilisant directement bcc.
    Pourrait se contourner en extrayant les valeurs d’intérêt dans un BPF_HASH
    séparé. Cette limitation est suggéré par la [doc][ebpf_exporter-doc-hash] «
    Maps coming from the kernel are binary encoded. Values are always u64, but
    keys can be primitive types like u64 or structs. »

[ebpf_exporter-doc-hash]: https://github.com/cloudflare/ebpf_exporter/#labels

  - Il semble impossible de récupérer des valeurs avec [perf_submit][bcc-doc-perf_submit]. Pas d’exemple ou d’explication trouvé.

[ebpf-doc-perf_submit]: https://github.com/iovisor/bcc/blob/master/docs/reference_guide.md#3-perf_submit

## Déploiement dans docker

### Ressources

- https://prometheus.io/docs/prometheus/latest/getting_started/ : exemple docker compose
- https://hub.docker.com/r/prom/prometheus : image docker officielle

### « Coller » deux container (par exemple pour ajouter un débugger)

- https://stackoverflow.com/questions/31007934/strace-to-monitor-dockerized-application-activity, https://medium.com/@rothgar/how-to-debug-a-running-docker-container-from-a-separate-container-983f11740dc6, https://gist.github.com/justincormack/f2444fbdf210b05d4f7baabe6fcd219a

- TODO Tester

## Programmation avec bcc

### Attention

- les indirections ne sont pas toujours bien gérées par le compilateur, il faut parfois passer par une variable intermédiaire

### Problème : exécuter un bcctool dans un container

Pistes :

 - https://www.youtube.com/watch?v=Yrjk4W-F9iY
 - http://www.adelzaalouk.me/2017/security-bpf-docker-cillium/#security-policies-using-ebpf

 Finalement, on exécutera les bcctools en local dans un premier temps, on a quand même accès à un certain nombre d’informations

## Ressources

### https://medium.com/@andrewhowdencom/adventures-with-ebpf-and-prometheus-6a59dd170b26

- Prometheus <-> Exporter <-> eBPF
- Nombreux exporter existants, peut-être que
  [ebpf_exporter][ebpf_exporter] fait une partie de ce qu’on veut
  

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

[bcc]: http://iovisor.github.io/bcc/
[bcc-install]: https://github.com/iovisor/bcc/blob/master/INSTALL.md
[bcc-example]: https://github.com/iovisor/bcc/blob/master/examples/networking/
[seccomp-bpf]: https://blog.yadutaf.fr/2014/05/29/introduction-to-seccomp-bpf-linux-syscall-filter/
[ebpf_exporter]: https://github.com/cloudflare/ebpf_exporter
[ebpf_exporter_example]: https://github.com/cloudflare/ebpf_exporter#examples
[5]: https://kubernetes.io/blog/2017/12/using-ebpf-in-kubernetes/
[bcc-tools]: https://github.com/iovisor/bcc/blob/master/README.md#tools
[xdpcap]: https://github.com/cloudflare/xdpcap
[xdpcap-context]: https://blog.cloudflare.com/xdpcap/

