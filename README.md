# Port Proxy

A Fast and Simple util tool to tunnel the local TCP or UDP port to another port. 

[![Build Status](https://travis-ci.org/enbiso/portproxy.svg?branch=master)](https://travis-ci.org/enbiso/portproxy)

## Examples

Simple example of tunnel port 80 to port 8080 locally
```
portproxy --source 127.0.0.1:80 --dest 127.0.0.1:8080
```

To publicly expose the web server which is only accesible in local
```
portproxy --source 127.0.0.1:80 --dest :8080
```

Tunnel UDP traffic in port 6000 to 8000
```
portproxy --source 127.0.0.1:6000 --dest :8000 --protocol udp
```
