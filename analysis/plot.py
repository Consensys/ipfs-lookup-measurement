import datetime
import json
import re
import os
import matplotlib.pyplot as plt
import numpy as np


# Results for pvd
pvd_latencies = []
pvd_search_latencies = []
pvd_put_latencies = []
pvd_put_latencies_succeed = []
pvd_put_latencies_failed = []
pvd_agents = dict()
pvd_agents_succeed = dict()
pvd_agents_failed = dict()
# Results for ret
ret_latencies = []
ret_search_latencies = []
ret_get_latencies = []
ret_agents = dict()

# Helper function to trim agent version
def trim_agent(agent):
    if agent.startswith("/"):
        agent = agent[1:]
    if agent.startswith("go-ipfs"):
        return "go-ipfs"
    elif agent.startswith("hydra-booster"):
        return "hydra-booster"
    elif agent.startswith("storm"):
        return "storm"
    elif agent.startswith("ioi"):
        return "ioi"
    elif agent.startswith("n.a."):
        return "n.a."
    else:
        return "others"

def get_millis(duration):
    if duration.endswith("µs"):
        return float(duration[:-2]) / 1000
    elif duration.endswith("ms"):
        return float(duration[:-2])
    elif duration.endswith("s"):
        if "m" in duration:
            return int(duration[0]) * 60 * 1000 + get_millis(duration[2:])
        return float(duration[:-2]) * 1000


def millis_interval(start_time, end_time):
    """start and end are datetime instances"""
    diff = end_time - start_time
    millis = diff.days * 24 * 60 * 60 * 1000
    millis += diff.seconds * 1000
    millis += diff.microseconds / 1000
    return millis

for file in ["1.log", "2.log", "3.log", "4.log", "5.log"]:
    pvd_on = False
    pvd_start = None
    pvd_search = None
    ret_on = False
    ret_start = None
    ret_search = None
    if not os.path.isfile(file):
        break
    for line in reversed(list(open(file))):
        data = json.loads(line)
        if data["line"].startswith("Start retrieving"):
            ret_start = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            ret_on = True
        elif data["line"].startswith("Got provider"):
            ret_search = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            ret_search_latencies.append(millis_interval(ret_start, ret_search))
        elif data["line"].startswith("Done retrieving"):
            end = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            ret_get_latencies.append(millis_interval(ret_search, end))
            ret_latencies.append(millis_interval(ret_start, end))
            ret_on = False
        elif data["line"].startswith("Start providing"):
            pvd_start = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            pvd_on = True
        elif data["line"].startswith("In total"):
            pvd_search = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            pvd_search_latencies.append(millis_interval(pvd_start, pvd_search))
        elif data["line"].startswith("Finish providing"):
            end = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            pvd_put_latencies.append(millis_interval(pvd_search, end))
            pvd_latencies.append(millis_interval(pvd_start, end))
            pvd_on = False
        elif "(" in data["line"]:
            if ret_on:
                agent = data["line"].split("(")[-1].split(")")[0]
                if agent in ret_agents:
                    ret_agents[agent] += 1
                else:
                    ret_agents[agent] = 1
            elif pvd_on:
                agent = data["line"].split("(")[-1].split(")")[0]
                if agent in pvd_agents:
                    pvd_agents[agent] += 1
                else:
                    pvd_agents[agent] = 1
                if data["line"].startswith("Succeed in"):
                    if agent in pvd_agents_succeed:
                        pvd_agents_succeed[agent] += 1
                    else:
                        pvd_agents_succeed[agent] = 1
                    latency = data["line"].split(":")[-1][1:]
                    pvd_put_latencies_succeed.append(get_millis(latency))
                elif data["line"].startswith("Error putting"):
                    if agent in pvd_agents_failed:
                        pvd_agents_failed[agent] += 1
                    else:
                        pvd_agents_failed[agent] = 1
                    latency = data["line"].split(":")[-1][1:]
                    pvd_put_latencies_failed.append(get_millis(latency))

# Sort dictionary
ret_agents = {k: v for k, v in sorted(ret_agents.items(), key=lambda item: item[1], reverse=True)}
pvd_agents = {k: v for k, v in sorted(pvd_agents.items(), key=lambda item: item[1], reverse=True)}

# Generate graphs.
# Overall pvd latency
plt.rc('font', size=8)
plt.hist(list(pvd_latencies), bins=20, density=False)
plt.title("Content publish overall latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/pvd_latency.png")
# PVD search latency
plt.clf()
plt.hist(list(pvd_search_latencies), bins=20, density=False)
plt.title("Content publish search latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/pvd_search_latency.png")
# PVD put latency
plt.clf()
plt.hist(list(pvd_put_latencies), bins=20, density=False)
plt.title("Content publish put latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/pvd_put_latency.png")
# PVD put latency succeed
plt.clf()
plt.hist(list(pvd_put_latencies_succeed), bins=20, density=False)
plt.title("Content publish put latency (succeed)")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/pvd_put_latency_succeed.png")
# PVD put latency failed
plt.clf()
plt.hist(list(pvd_put_latencies_failed), bins=20, density=False)
plt.title("Content publish put latency (failed)")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/pvd_put_latency_failed.png")
# PVD across agents
plt.clf()
total = 0
trim_agents = dict()
for agent, count in pvd_agents.items():
    agent = trim_agent(agent)
    total += count
    if agent in trim_agents:
        trim_agents[agent] += count
    else:
        trim_agents[agent] = count
plt.pie(trim_agents.values(), labels=trim_agents.keys(), autopct="%.1f%%")
plt.title("Agent came across in content publish, total {} records".format(total))
plt.savefig("./figs/pvd_all_agents.png")
# PVD across agents succeed
plt.clf()
total = 0
trim_agents = dict()
for agent, count in pvd_agents_succeed.items():
    agent = trim_agent(agent)
    total += count
    if agent in trim_agents:
        trim_agents[agent] += count
    else:
        trim_agents[agent] = count
plt.pie(trim_agents.values(), labels=trim_agents.keys(), autopct="%.1f%%")
plt.title("Agent came across succeed in storing provider record, total {} records".format(total))
plt.savefig("./figs/pvd_succeed_agents.png")
# PVD across agents failed
plt.clf()
total = 0
trim_agents = dict()
for agent, count in pvd_agents_failed.items():
    agent = trim_agent(agent)
    total += count
    if agent in trim_agents:
        trim_agents[agent] += count
    else:
        trim_agents[agent] = count
plt.pie(trim_agents.values(), labels=trim_agents.keys(), autopct="%.1f%%")
plt.title("Agent came across fail to store provider record, total {} records".format(total))
plt.savefig("./figs/pvd_failed_agents.png")
# Overall ret latency
plt.clf()
plt.rc('font', size=8)
plt.hist(list(ret_latencies), bins=20, density=False)
plt.title("Content fetch overall latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/ret_latency.png")
# ret search latency
plt.clf()
plt.hist(list(ret_search_latencies), bins=20, density=False)
plt.title("Content fetch search latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/ret_search_latency.png")
# ret get latency
plt.clf()
plt.hist(list(ret_get_latencies), bins=20, density=False)
plt.title("Content fetch get latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/ret_get_latency.png")
# ret across agents
plt.clf()
total = 0
trim_agents = dict()
for agent, count in ret_agents.items():
    agent = trim_agent(agent)
    total += count
    if agent in trim_agents:
        trim_agents[agent] += count
    else:
        trim_agents[agent] = count
plt.pie(trim_agents.values(), labels=trim_agents.keys(), autopct="%.1f%%")
plt.title("Agent came across in content fetch, total {} records".format(total))
plt.savefig("./figs/ret_all_agents.png")