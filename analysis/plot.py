import datetime
import json
import os
import matplotlib.pyplot as plt


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


# Results for pvd
pvd_latencies = []
pvd_dht_walk_latencies = []
pvd_put_latencies = []
pvd_put_latencies_succeed = []
pvd_put_latencies_failed = []
pvd_agents = dict()
pvd_dht_walk_agents = dict()
pvd_put_agents = dict()
pvd_put_agents_succeed = dict()
pvd_put_agents_failed = dict()
# Results for ret
ret_latencies = []
ret_dht_walk_latencies = []
ret_get_latencies = []
ret_agents = dict()
ret_dht_walk_agents = dict()
ret_get_agents = dict()
# Error messages.
errs = []
# Error getting provider record after a successful put
agents_1 = dict()
errs_1 = []
# Get empty provider record after a successful put
agents_2 = dict()


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
        if data["line"].startswith("Start providing"):
            pvd_start = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            pvd_on = True
        elif data["line"].startswith("Finish providing"):
            end = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            pvd_put_latencies.append(millis_interval(pvd_search, end))
            pvd_latencies.append(millis_interval(pvd_start, end))
            pvd_on = False
        elif data["line"].startswith("Start retrieving"):
            ret_start = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            ret_on = True
        elif data["line"].startswith("Done retrieving"):
            end = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
            ret_get_latencies.append(millis_interval(ret_search, end))
            ret_latencies.append(millis_interval(ret_start, end))
            ret_on = False
        else:
            if pvd_on:
                agent = trim_agent(data["line"].split("(")[-1].split(")")[0])
                if data["line"].startswith("Getting closest peers for cid"):
                    pvd_agents[agent] = pvd_agents.get(agent, 0) + 1
                    pvd_dht_walk_agents[agent] = pvd_dht_walk_agents.get(agent, 0) + 1
                elif data["line"].startswith("Succeed in putting provider record"):
                    pvd_agents[agent] = pvd_agents.get(agent, 0) + 1
                    pvd_put_agents[agent] = pvd_put_agents.get(agent, 0) + 1
                    pvd_put_agents_succeed[agent] = pvd_put_agents_succeed.get(agent, 0) + 1
                    latency = data["line"].split(":")[-1][1:]
                    pvd_put_latencies_succeed.append(get_millis(latency))
                elif data["line"].startswith("Error putting provider record"):
                    pvd_agents[agent] = pvd_agents.get(agent, 0) + 1
                    pvd_put_agents[agent] = pvd_put_agents.get(agent, 0) + 1
                    pvd_put_agents_failed[agent] = pvd_put_agents_failed.get(agent, 0) + 1
                    latency = data["line"].split(":")[-1][1:]
                    pvd_put_latencies_failed.append(get_millis(latency))
                    if agent == "hydra-booster":
                        errs.append(data["line"].split(")")[1].split("time taken")[0])
                elif data["line"].startswith("In total"):
                    pvd_search = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
                    pvd_dht_walk_latencies.append(millis_interval(pvd_start, pvd_search))
                elif data["line"].startswith("Error getting provider record for cid"):
                    agents_1[agent] = agents_1.get(agent, 0) + 1
                    errs_1.append(data["line"].split("after a successful put")[1])
                elif data["line"].startswith("Got 0 provider records back from"):
                    agents_2[agent] = agents_2.get(agent, 0) + 1
            if ret_on:
                agent = trim_agent(data["line"].split("(")[-1].split(")")[0])
                if data["line"].startswith("Getting providers for cid"):
                    ret_agents[agent] = ret_agents.get(agent, 0) + 1
                    ret_dht_walk_agents[agent] = ret_dht_walk_agents.get(agent, 0) + 1
                elif data["line"].startswith("Found provider for cid"):
                    ret_agents[agent] = ret_agents.get(agent, 0) + 1
                    ret_get_agents[agent] = ret_get_agents.get(agent, 0) + 1
                elif data["line"].startswith("Connected to provider"):
                    ret_search = datetime.datetime.strptime(data["timestamp"][:-9], "%Y-%m-%dT%H:%M:%S.%f")
                    ret_dht_walk_latencies.append(millis_interval(ret_start, ret_search))
                pass
            pass

# Generate graphs
plt.rc('font', size=8)
plt.rcParams["figure.figsize"] = (10,6)
# Overall pvd latency
plt.clf()
plt.hist(list(pvd_latencies), bins=20, density=False)
plt.title("Content publish overall latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/pvd_latency.png")
# DHT Walk pvd latency
plt.clf()
plt.hist(list(pvd_dht_walk_latencies), bins=20, density=False)
plt.title("Content publish DHT walk latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/pvd_dht_walk_latency.png")
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
# PVD all agents
plt.clf()
plt.barh(list(pvd_agents.keys()), list(pvd_agents.values()))
for index, value in enumerate(list(pvd_agents.values())):
    plt.text(value, index, str(value))
plt.title("All agents encountered in overall content publish")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/pvd_agents.png")
# PVD dht walk agents
plt.clf()
plt.barh(list(pvd_dht_walk_agents.keys()), list(pvd_dht_walk_agents.values()))
for index, value in enumerate(list(pvd_dht_walk_agents.values())):
    plt.text(value, index, str(value))
plt.title("All agents encountered in DHT walk content publish")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/pvd_agents_dht_walk.png")
# PVD dht put agents
plt.clf()
plt.barh(list(pvd_put_agents.keys()), list(pvd_put_agents.values()))
for index, value in enumerate(list(pvd_put_agents.values())):
    plt.text(value, index, str(value))
plt.title("All agents encountered in PUT content publish")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/pvd_agents_put.png")
# PVD dht put succeed agents
plt.clf()
plt.barh(list(pvd_put_agents_succeed.keys()), list(pvd_put_agents_succeed.values()))
for index, value in enumerate(list(pvd_put_agents_succeed.values())):
    plt.text(value, index, str(value))
plt.title("All agents encountered in PUT Succeed content publish")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/pvd_agents_put_succeed.png")
# PVD dht put failed agents
plt.clf()
plt.barh(list(pvd_put_agents_failed.keys()), list(pvd_put_agents_failed.values()))
for index, value in enumerate(list(pvd_put_agents_failed.values())):
    plt.text(value, index, str(value))
plt.title("All agents encountered in PUT Failed content publish")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/pvd_agents_put_failed.png")
# Overall ret latency
plt.clf()
plt.rc('font', size=8)
plt.hist(list(ret_latencies), bins=20, density=False)
plt.title("Content fetch overall latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/ret_latency.png")
# DHT Walk ret latency
plt.clf()
plt.hist(list(ret_dht_walk_latencies), bins=20, density=False)
plt.title("Content fetch DHT walk latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/ret_dht_walk_latency.png")
# Get ret latency
plt.clf()
plt.hist(list(ret_get_latencies), bins=20, density=False)
plt.title("Content fetch GET latency")
plt.xlabel("Latency (ms)")
plt.ylabel("Count")
plt.savefig("./figs/ret_get_latency.png")
# Ret all agents
plt.clf()
plt.barh(list(ret_agents.keys()), list(ret_agents.values()))
for index, value in enumerate(list(ret_agents.values())):
    plt.text(value, index, str(value))
plt.title("All agents encountered in overall content fetch")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/ret_agents.png")
# Ret dht walk agents
plt.clf()
plt.barh(list(ret_dht_walk_agents.keys()), list(ret_dht_walk_agents.values()))
for index, value in enumerate(list(ret_dht_walk_agents.values())):
    plt.text(value, index, str(value))
plt.title("All agents encountered in DHT walk content fetch")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/ret_agents_dht_walk.png")
# Ret dht get agents
plt.clf()
plt.barh(list(ret_get_agents.keys()), list(ret_get_agents.values()))
for index, value in enumerate(list(ret_get_agents.values())):
    plt.text(value, index, str(value))
plt.title("All agents encountered for getting a provider record")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/ret_agents_get.png")
# Hydra-booster error type, PVD put
errdict = dict()
for err in errs:
    for e in err.split("*")[1:]:
        if "127.0.0.1" not in e:
            if "i/o timeout" in e:
                errdict["i/o timeout"] = errdict.get("i/o timeout", 0) + 1
            elif "timeout: no recent network activity" in e:
                errdict["timeout: no recent \nnetwork activity"] = errdict.get("timeout: no recent \nnetwork activity", 0) + 1
            elif "dial backoff" in e:
                errdict["dial backoff"] = errdict.get("dial backoff", 0) + 1
            elif "context deadline exceeded" in e:
                errdict["context deadline \nexceeded"] = errdict.get("context deadline \nexceeded", 0) + 1
            else:
                print("uncaptured error", e)
plt.clf()
plt.barh(list(errdict.keys()), list(errdict.values()))
for index, value in enumerate(list(errdict.values())):
    plt.text(value, index, str(value))
plt.title("Error in putting provider record to hydra-booster nodes")
plt.xlabel("Count")
plt.ylabel("Error type")
plt.savefig("./figs/pvd_agents_put_failed_hydra_booster_err.png")
# Agents error getting provider record back
plt.clf()
plt.barh(list(agents_1.keys()), list(agents_1.values()))
for index, value in enumerate(list(agents_1.values())):
    plt.text(value, index, str(value))
plt.title("Agents for having error in getting provider record back after a successful put")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/pvd_agents_get_error.png")
# The error of getting provider record back
errdict = dict()
for err in errs_1:
    if "stream reset" in err:
        errdict["stream reset"] = errdict.get("stream reset", 0) + 1
    elif "timeout: no recent network activity" in err:
        errdict["timeout: no recent \nnetwork activity"] = errdict.get("timeout: no recent \nnetwork activity", 0) + 1
    elif "timed out reading response" in err:
        errdict["timed out \nreading response"] = errdict.get("timed out \nreading response", 0) + 1
    elif "Application error 0x0" in err:
        errdict["Application error \n0x0"] = errdict.get("Application error \n0x0", 0) + 1
    elif "failed to dial" in err:
        errdict["failed to dial"] = errdict.get("failed to dial", 0) + 1
    else:
        print("uncaptured error", err)
plt.clf()
plt.barh(list(errdict.keys()), list(errdict.values()))
for index, value in enumerate(list(errdict.values())):
    plt.text(value, index, str(value))
plt.title("The error in getting provider record back after a successful put")
plt.xlabel("Count")
plt.ylabel("Error type")
plt.savefig("./figs/pvd_agents_get_error_type.png")
# Empty provider record after a successful put
plt.clf()
plt.barh(list(agents_2.keys()), list(agents_2.values()))
for index, value in enumerate(list(agents_2.values())):
    plt.text(value, index, str(value))
plt.title("Agents for having empty provider record back after a successful put")
plt.xlabel("Count")
plt.ylabel("Agent")
plt.savefig("./figs/pvd_agents_get_empty.png")
