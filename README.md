## GuardMQ : Open Source Security proxy for peering between MQTT brokers

GuardMQ is a security proxy that manage peering connections between brokers written in GO.

At the moment it has support for:
- MQTT endpoint with TLS and MTLS
- Middleware skeleton
- Peering session between two brokers and a set of topics

### Installation

This project is in Go, so you have to instal the Go environment into your system: [https://golang.org/doc/install](https://golang.org/doc/install)
Check before that if your distribution integrates golang packages or not.

Then :
```
go build -o guardmq cmd/main.go
cp guardmq /usr/local/bin/
mkdir /etc/guardmq/
cp conf/guardmq.yml /etc/guardmq/
cp guardmq.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable guardmq.service
systemctl start guardmq.service
```

### Docker installation

You can build the Docker image for GuardMQ like this :

```
docker build -f deployment/Dockerfile -t guardmq .
```

You have a config file example in conf/ use it to run a container :

```
docker run -it -v ./conf/guardmq.yml:/usr/share/proxy/conf/guardmq.yml -p1883:1883 guardmq
```
