#!/usr/bin/env bash

# Query some website without https, which starts to be quite rare

for url in http://www.stealmylogin.com http://bu.univ-lorraine.fr/bu404 \
  http://ovh.net http://example.com/$(date +%s-%N); do
  curl -4 $url>/dev/null
  sleep 0.3
done
