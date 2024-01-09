# Kragins


## Build
```bash
make build
```

## Run
```bash
## copy config.example.yaml to _debug/kragins/config.yaml
## config your own postgres
./_debug/kragins/kragins serve

## test your service
curl http://localhost:30080/kragins/hello
```


## Release
```bash
make prod
```
