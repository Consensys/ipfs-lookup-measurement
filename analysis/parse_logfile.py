#!/usr/bin/env python3

import re

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

if __name__ == "__main__":
    for m in parseLines(fileName):
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

"""