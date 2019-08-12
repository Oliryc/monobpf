À remettre dans Overleaf
--------------------------


### Avancement dans la construction d’une chaîne de monitoring, par indicateurs

#### Légende

Dans les tableaux suivants, on note :

- `/` pour « L’outil n’est pas utilisé », car cela n’est pas nécessaire ou parce que l’outil n’est pas encore supporté.
- `X` pour « L’outil est utilisé et fonctionne ».

#### Avec filtrage

| Indicateur | eBPF | PCP | Vector | Temps employé à | Notes |
|------------|:-----:|:------:|:----------:|-----------------------------|----|
| Utilisation Mémoire/CPU/Disque IOPS/Réseau (hôte et container) | / | X | X | Configuration de PCP, intégration de PCP et de Vector | La sélection se fait par nom de container, depuis l’interface de vector |
| Flamegraph (BCC PMDA, module profile) | X | X | / | Comprendre comment lancer manuellement l’écoute | L’affichage dans Vector est possible d’après la documentation du module, mais n’a pas été possible pendant les tests |
| GC | X | X | X | Porter le programme eBPF, faire accepter les données à Vector pour l’affichage | Une solution de contournement a été trouvée, en manipulant les données reçues depuis le code JavaScript de Vector.
| Décompte des appels systèmes (module syscount) | X | X | / | Installation de `auditd` | *Instable*, risque d’overflow sur les décomptes |
| Latence des accès/écritures disque (modules {ext4,xfs,zfs}dist (restreignable à un processus)) | X | X | X | Difficile de tester le filtrage sur un processus |
| Retransmissions TCP (module tcpretrans) | X | X | X  | Comment provoquer une retransmission TCP ? |


#### Filtrage à tester

| Indicateur | eBPF | PCP | Vector | Temps employé à | Notes |
|------------|:-----:|:------:|:----------:|-----------------------------|----|
| Connexion TCP (autres modules tcp_*) | X | X | / | Pas de difficulté particulière |
| Inspection de header avec XDP | / | / | / | Comprendre pourquoi certains paquets restaient inaccessibles |

#### Sans filtrage

| Indicateur | eBPF | PCP | Vector | Temps employé à | Notes |
|------------|:-----:|:------:|:----------:|-----------------------------|----|
| Connexion TCP (module TCPLife (avec filtrage en théorie), tcptop(sans filtrage)) | X | X | X  | Tentatives pour faire fonctionner le filtrage |
| Flamegraph | X | X | X | Compiler le PMDA Vector | eBPF n’est pas toujours indispensable |
| Commandes lancées sur le système (execsnoop) | X | X | X | Pas de difficulté particulière | Peu d’intérêt à le lancer sur le processus à monitorer, voir plutôt ce que le processus à monitorer lance |
| Latence des accès/écritures disque (biolatency) | X | X | X |  |


#### TODO

- Insérer ce document dans Overleaf

# Pour insertion dans le rapport

Légende :

- `E`: La fonctionnalité existe dans le module considéré
- `R`: J’ai réalisé la fonctionnalité au cours de mon stage
- `/` pour « L’outil n’est pas utilisé », car cela n’est pas nécessaire ou parce que l’outil n’est pas encore supporté.

| Indicateur | eBPF | PCP | Vector | Utilisation | Notes |
|------------|:-----:|:------:|:----------:|-----------------------------|----|
| Connexion TCP | X | X | X  | Tentatives pour faire fonctionner le filtrage |
| Flamegraph | X | X | X | Compiler le PMDA Vector | eBPF n’est pas toujours indispensable |
| Commandes lancées sur le système (execsnoop) | X | X | X | Pas de difficulté particulière | Peu d’intérêt à le lancer sur le processus à monitorer, voir plutôt ce que le processus à monitorer lance |
| Latence des accès/écritures disque (biolatency) | X | X | X |  |
| Utilisation Mémoire/CPU/Disque IOPS/Réseau (hôte et container) | / | X | X | Configuration de PCP, intégration de PCP et de Vector | La sélection se fait par nom de container, depuis l’interface de vector |
| Flamegraph (BCC PMDA, module profile) | X | X | / | Comprendre comment lancer manuellement l’écoute | L’affichage dans Vector est possible d’après la documentation du module, mais n’a pas été possible pendant les tests |
| GC | X | X | X | Porter le programme eBPF, faire accepter les données à Vector pour l’affichage | Une solution de contournement a été trouvée, en manipulant les données reçues depuis le code JavaScript de Vector.
| Décompte des appels systèmes (module syscount) | X | X | / | Installation de `auditd` | *Instable*, risque d’overflow sur les décomptes |
| Latence des accès/écritures disque (modules {ext4,xfs,zfs}dist (restreignable à un processus)) | X | X | X | Difficile de tester le filtrage sur un processus |
| Retransmissions TCP (module tcpretrans) | X | X | X  | Comment provoquer une retransmission TCP ? |
| Connexion TCP (autres modules tcp_*) | X | X | / | Pas de difficulté particulière |
| Inspection HTTP | / | / | / | Comprendre pourquoi certains paquets restaient inaccessibles |
| Inspection UDP/DTLS | / | / | / | Comprendre pourquoi certains paquets restaient inaccessibles |


