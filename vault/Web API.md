# Web API for logserve

The api will accept POST requests to `/log` with logging information to log data and GET requests from `/data`.

## Requests to

Logs are sent to `/log` with POST in plain text. Logs are sent in form of queries:

- [x] Logging
    - `LOG <log information>`
        Logs the log information
- [ ] Key-value
    - `KEY <key> VALUE`
        Records a new value for the key
    - `KEY <key> REMOVE`
        Removes the key
- [ ] Benchmark
    - `BENCH <name> <start> <finish>`
        Benchmarks the event

## Getting data from the website

Requests are sent to `/data` with GET.

- [x] Getting state in json form
        `/data/json`
        returns in form `{keyvalues, log: [{timestamp, query}]}`
- [x] Getting logs since a timestamp
        `/data/since?t=<timestamp>`
        returns all the logs in raw form since `timestamp` in form `[{timestamp, query}]`
- [ ] Static page returning from `/data`
- [ ] Display log information
- [ ] Display Key-value info
- [ ] Ability to hide log info or key-value info or both
- [ ] Ability to open them in different browser windows
