# Controller
## Start ipfs for local testing
`ipfs daemon --init`
## Build
`make build`
## Integration testing with local ipfs
`make itest`
## Running simple nodes experiment
```
# editing nodes-list.out file with ipfs nodes, hostname:port
# run experiment
./controller

# or
./controller -l <nodes list file>
```
