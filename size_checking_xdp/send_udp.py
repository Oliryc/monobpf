#!/usr/bin/env python3

import socket
from sys import argv, exit
from time import sleep

UDP_IP = "127.0.0.1"
UDP_PORT = 5005
MESSAGE = "Hello, World!\n"

if __name__ == '__main__':
    print("UDP target IP:", UDP_IP)
    print("UDP target port:", UDP_PORT)
    print("message:", MESSAGE)

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) # UDP
    sock.sendto(bytes(MESSAGE, "utf-8"), (UDP_IP, UDP_PORT))

    if (len(argv) < 3):
        print("usage: {} min_packet_length max_packet_length".format(argv[0]))
        exit(1)

    min_packet_length = int(argv[1])
    max_packet_length = int(argv[2])

    for i in range(min_packet_length,max_packet_length+1):
        l=i
        msg="\0"*l
        print(l)
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) # UDP
        sock.sendto(bytes(msg, "utf-8"), (UDP_IP, UDP_PORT))
        sleep(0.1)
