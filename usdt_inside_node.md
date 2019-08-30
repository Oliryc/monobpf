# USDT & Node

Face aux limites d’eBPF pour accéder à des informations de profiling
précises, on veut explorer les possibilités avec des sondes personnalisés (cf `./usdt_node`).

- Notons que ceci nous fait nous écarter du monitoring non invasif.

- Pas besoin de recompiler node

## Installation des probes

Automatiser ?

## Questions

- Avantages à utiliser USDT plutôt que des outils existants ? 
  - Plus de perf ? peut-être, vu qu’on pourrait n’activer que les sondes voulues. (sonde activée : ≃80µs, non activée, ≃10µs)
  - Plus de précisions ? Sans doute, pas de sampling
- Impact négatif sur les perf d’une sonde qui ne serait pas activée ?
- Le sampling est-il vraiment un problème dans l’éco-système actuel de mesure de perf ? (en particulier dans l’écosystème javascript)
- Est-ce qu’on n’est pas en train de refaire des outils existants en moins bien ?

## Ressources

- http://www.brendangregg.com/blog/2016-10-12/linux-bcc-nodejs-usdt.html
  - L’avantage d’USDT serait d’être intégré au reste des outils eBPF (notamment pour voir les interractions avec le reste du système)…
- https://nodejs.org/en/docs/guides/simple-profiling/
- http://www.brendangregg.com/perf.html#TimedProfiling & https://jvns.ca/blog/2016/03/12/how-does-perf-work-and-some-questions/

## Réu

- Démo USDT sur connecteur
- Flamegraph
  - https://www.alxolr.com/articles/squeeze-node-js-performance-with-flame-graphs (envoyé par mail) et https://nodejs.org/en/docs/guides/diagnostics-flamegraph/ (official documentation) -> 0x
  - sur le connecteur, problème dans la transmission du SIGINT : https://github.com/davidmarkclements/0x/issues/136 avec pm2
- Impact perf ?
