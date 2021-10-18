#!/usr/bin/env python3

import re
import time
import datetime

fileName = "log-1.txt"

def parseLines(fileName):
    with open(fileName) as logFile:
        for line in logFile.readlines():
            m = re.search("^"+
            "(?P<logtime>[\d-]{10}\W[\d:]+)"+
            "\W+"+

            "(?P<logkey>"+
             "(?:Start providing cid)"+
            "|(?:Finish providing cid)"+
            "|(?:Start retrieving content for)"+
            "|(?:Done retrieving content for)"+
            ")\W+"+
            
            "(?P<cid>\w+)", line)
            
            if m:
                yield m.groupdict()

def str2time(s):
    return time.mktime(datetime.datetime.strptime(s, "%Y-%m-%d %H:%M:%S").timetuple())

def yieldElapsed(fileName):
    events = dict()
    for m in parseLines(fileName):
        cid = m["cid"]
        logkey = m["logkey"]
        logtime = str2time(m["logtime"])
        if logkey == "Start providing cid":
            events[cid] = dict(start=logtime, event="providing")
        elif logkey == "Finish providing cid":
            start = events[cid]["start"]
            events[cid].update(elapsed=logtime-start)
            y = events[cid]
            del(events[cid])
            yield y
        elif logkey == "Start retrieving content for":
            events[cid] = dict(start=logtime, event="retrieving")
        elif logkey == "Done retrieving content for":
            start = events[cid]["start"]
            events[cid].update(elapsed=logtime-start)
            y = events[cid]
            del(events[cid])
            yield y

if __name__ == "__main__":
    for m in parseLines(fileName):
        print(m)
    print("Elapsed:")
    for m in yieldElapsed(fileName):
        print(m)

"""
example output:
{'logtime': '2021-10-18 12:48:31', 'logkey': 'Start providing cid', 'cid': 'QmRx5dXiQHECSdj8Ftwo4MtcUSCJDnxRLtAEAakx7NXyfK'}
{'logtime': '2021-10-18 12:48:43', 'logkey': 'Finish providing cid', 'cid': 'QmRx5dXiQHECSdj8Ftwo4MtcUSCJDnxRLtAEAakx7NXyfK'}
{'logtime': '2021-10-18 12:49:32', 'logkey': 'Start providing cid', 'cid': 'QmewAtyyoEJ9GQb6KvPXfCUcr5CqGrX1A1da4esredPjG3'}
{'logtime': '2021-10-18 12:50:20', 'logkey': 'Finish providing cid', 'cid': 'QmewAtyyoEJ9GQb6KvPXfCUcr5CqGrX1A1da4esredPjG3'}
{'logtime': '2021-10-18 12:51:57', 'logkey': 'Start providing cid', 'cid': 'QmUjSYUyreuC9zrMqhExTbGPM8gBGahRW4BChwPgnPceoc'}
{'logtime': '2021-10-18 12:52:57', 'logkey': 'Finish providing cid', 'cid': 'QmUjSYUyreuC9zrMqhExTbGPM8gBGahRW4BChwPgnPceoc'}
{'logtime': '2021-10-18 12:54:01', 'logkey': 'Start retrieving content for', 'cid': 'QmQo4hQdKyqjfLtcCMDRTZMsxZ5jF25VEtVDhLQqP5H8EU'}
{'logtime': '2021-10-18 12:54:03', 'logkey': 'Done retrieving content for', 'cid': 'QmQo4hQdKyqjfLtcCMDRTZMsxZ5jF25VEtVDhLQqP5H8EU'}
{'logtime': '2021-10-18 12:54:33', 'logkey': 'Start retrieving content for', 'cid': 'QmVeFwXAmPGFktmar785MqjKGEgrcRWwxrty3TEsekzYQ8'}
{'logtime': '2021-10-18 12:54:35', 'logkey': 'Done retrieving content for', 'cid': 'QmVeFwXAmPGFktmar785MqjKGEgrcRWwxrty3TEsekzYQ8'}
{'logtime': '2021-10-18 12:54:54', 'logkey': 'Start retrieving content for', 'cid': 'QmXX5QuryUbLBwMxYNFnr1gu6fHZRHvTPzRV2EnJcBw9uL'}
{'logtime': '2021-10-18 12:54:55', 'logkey': 'Done retrieving content for', 'cid': 'QmXX5QuryUbLBwMxYNFnr1gu6fHZRHvTPzRV2EnJcBw9uL'}
{'logtime': '2021-10-18 12:55:35', 'logkey': 'Start retrieving content for', 'cid': 'QmVRBagUJMHgskxqXSJVsZt5GBsZy59uNvUyQtmaxYwTaD'}
{'logtime': '2021-10-18 12:55:36', 'logkey': 'Done retrieving content for', 'cid': 'QmVRBagUJMHgskxqXSJVsZt5GBsZy59uNvUyQtmaxYwTaD'}
Elapsed:
{'start': 1634521711.0, 'event': 'providing', 'elaps': 12.0}
{'start': 1634521772.0, 'event': 'providing', 'elaps': 48.0}
{'start': 1634521917.0, 'event': 'providing', 'elaps': 60.0}
{'start': 1634522041.0, 'event': 'retrieving', 'elaps': 2.0}
{'start': 1634522073.0, 'event': 'retrieving', 'elaps': 2.0}
{'start': 1634522094.0, 'event': 'retrieving', 'elaps': 1.0}
{'start': 1634522135.0, 'event': 'retrieving', 'elaps': 1.0}
"""