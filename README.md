# Albion status checker

Simple albion server status checker with discord webhooks.

Running in docker:

```console
docker build -t checker .
docker run --rm -d -e WEBHOOK_URL={INSERT YOUR DISCORD WEBHOOK} --name checker checker
```

Run yourself:

```console
go get -v -d ./...
go build .
export WEBHOOK_URL={INSERT YOUR DISCORD WEBHOOK}
./albion-status-checker
```
