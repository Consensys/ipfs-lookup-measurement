#!/usr/bin/env python3

import os
from subprocess import run

num_nodes = 5

def setLokiEnv():
    monitor_addr = "http://localhost:3100"
    with open('nodes-list.out') as f:
        monitor_addr = f.readline()
        monitor_addr = monitor_addr.split("\"")[1]
        monitor_addr = "http://{}:3100/".format(monitor_addr)
    os.putenv("LOKI_ADDR", monitor_addr)

def downloadLogs():
    for i in range(1,num_nodes+1):
        cmd = """ logcli query --limit=987654321 --since=24h --output=jsonl '{host="node%d"}' 2>/dev/null >%d.log """ % (i, i)
        run(cmd, shell=True)

if __name__ == "__main__":
    setLokiEnv()
    downloadLogs()
