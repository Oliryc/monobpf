#!/usr/bin/env python3
# -*- coding: utf-8 -*-

# Convert ip in byte representation of XDP to human readable version

import sys, ipaddress

ipn = int(sys.argv[1])
print(str(ipaddress.ip_address(ipn)))
