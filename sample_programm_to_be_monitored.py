#!/usr/bin/env python3

import time, gc

if __name__ == '__main__':
    while True:
        time.sleep(3)
        gc.collect()
        print("Collected")
