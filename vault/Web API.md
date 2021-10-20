# Web API for logserve

The api will accept PUT requests to `/log` with logging information to log data and GET requests from `/record`.

## Requests to

Requests are sent to `/log` with PUT in plain text.

- [ ] Logging
    - `<timestamp> LOG <log information>`
        Logs the log information
- [ ] Key-value
    - `<timestamp> KEY <key> VALUE`
        Records the new value for the key
    - `<timestamp> KEY <key> REMOVE`
        Removes the key
- [ ] Benchmark
    - `<timestamp> BENCH <name> <start> <finish>`
        Benchmarks the event

## Getting data from the website

Requests are sent to `/record` with GET.

- [ ] Display log information
- [ ] Display Key-value info
- [ ] Ability to hide log info or key-value info or both
- [ ] Ability to open them in different browser windows
