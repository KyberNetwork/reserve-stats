## HTTP Gateway

This is gateway which redirect request to server go the right component

## Build

```shell
    cd cmd
    go build -v .
```

## Test 

```shell
    go test -v ./...
```

## Run

```shell
    cd cmd
    ./cmd
```

**Available flags**

**help**: show command manual
**write-access-key**: the key for header signature
**write-secret-key**: the secret for header signature
**listen**: the host where the component run on
**trade-logs-url**: url where gateway redirect to trade logs component
**user-url**: url where gateway redirect to user component