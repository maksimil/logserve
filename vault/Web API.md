# Web API for logserve

The api will accept POST requests to `/log` with logging information to log data and GET requests from `/data`.

## Requests to

Logs are sent to `/log` with POST in plain text. Logs are sent in form of queries. Values should not include the = sign.

- [x] Logging
    - `LOG <log info>`  
         Logs the log information
- [x] Key-value
    - `KEY_SET key=<key> value=<value>`  
         Records a new value for the key
    - `KEY_REMOVE key=<key>`  
         Removes the key
- [ ] Benchmark
    - `BENCH name=<name> start=<start> finish=<finish>`  
         Benchmarks the event

## Getting data from the website

Requests are sent to `/data` with GET.

- [x] Getting json of the state
        `/data/json`  
         returns in form `{keyvalues, log: [{timestamp, query}]}`
- [x] Getting logs since a timestamp  
         `/data/since?t=<timestamp>`  
         returns all the logs in raw form since `timestamp` in form `[{timestamp, query}]`
- [ ] add group argument
- [ ] add filtering in groups
- [x] Static page returning from `/data`
- [ ] Display log information in form of `<timestamp> <type> <values from data>`
- [ ] Display key-value info
- [ ] Ability to hide log info or key-value info or both
- [ ] Ability to open them in different browser windows
