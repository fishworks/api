api
===

[![wercker status](https://app.wercker.com/status/78d0e14acc95b790f0bcec6023599cfe/m "wercker status")](https://app.wercker.com/project/bykey/78d0e14acc95b790f0bcec6023599cfe)

This is an experimental proof-of-concept for migrating Deis' controller component to Go. Please note that this is a work in progress, so things are subject to change.

# Compiling from Source

```bash
$ make && make install
```

# Usage

To run by default at `tcp://0.0.0.0:8080`:

```bash
$ api
```

To run at a different address:

```bash
$ api --addr tcp://127.0.0.1:4567
```

Or, on a unix socket!

```bash
$ api --addr unix:///var/run/api.sock
```
