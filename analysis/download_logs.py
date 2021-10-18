#!/usr/bin/env python3

import os
from subprocess import run

def setLokiEnv():
    os.putenv("LOKI_ADDR", "http://54.66.169.156:3100/")

def downloadLogs():
    for i in range(1,6):
        cmd = """ logcli query --limit=987654321 --since=24h --output=jsonl '{host="node%d"}' 2>/dev/null >%d.log """ % (i, i)
        run(cmd, shell=True)

if __name__ == "__main__":
    setLokiEnv()
    downloadLogs()
