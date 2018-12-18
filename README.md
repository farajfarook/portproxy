# Port Proxy

A Fast and Simple util tool to tunnel the TCP or UDP port to another port. 

[![Build Status](https://travis-ci.org/enbiso/portproxy.svg?branch=master)](https://travis-ci.org/enbiso/portproxy)

## Examples

Simple example of tunnel port 80 to port 8080 locally in PC or network
```
portproxy --source 127.0.0.1:80 --target 127.0.0.1:8080
```

To publicly expose the web server which is only accesible in local PC or network
```
portproxy --source 172.0.0.10:80 --target :8080
```

Tunnel UDP traffic in port 6000 to 8000
```
portproxy --source 172.10.20.0:6000 --target :8000 --protocol udp
```


## Docker

```
docker run -p 8000:8000/udp \
    enbiso/portproxy \
    --source 127.0.0.1:6000 --target :8000 --protocol udp
```