
# How to

1. https://mise.jdx.dev/
2. `mise trust`

### Docker network
```shell
docker network create grpc-gateway-network
```

### GRPC upstream

Build binary of GRPC upstream locally
```shell
just build # or "mise exec -- just build"
```

Build docker image of GRPC upstream and run it
```shell
docker build -t grpc-gateway-upstream:latest -f docker/upstream.Dockerfile . && \
docker run -d -p 9000:9000 --network grpc-gateway-network --name grpc-gateway-upstream grpc-gateway-upstream:latest -addr 0.0.0.0:9000
```

### Envoy

https://www.envoyproxy.io/
https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/grpc_json_transcoder_filter

```shell
docker build -t grpc-gateway-envoy:latest -f docker/envoy.Dockerfile . && \
docker run -d -p 9001:9001 -p 9901:9901 --network grpc-gateway-network --name grpc-gateway-envoy grpc-gateway-envoy:latest
```

### Debug network container

https://hub.docker.com/r/nicolaka/netshoot

```shell
docker run --rm -it --network grpc-gateway-network nicolaka/netshoot

```

### Play

```shell
http PUT grpc-gateway-envoy:9001/api/v1/users name='Alice' age=21
```

# Keep in mind

### Single proto descriptor per route rule

`proto_descriptor` field in Envoy's YAML config accepts ***single*** string (filepath).
So that's why we either:
- Have one descriptor file from multiple (or all) Proto services (one big gen command)
- Have multiple sections of Envoy rules, each for their own Proto service and HTTP route (is it possible?)

Additional R&D required:
- Does the https://buf.build/ supports proto descriptor files? From multiple sources?
- Envoy's config divide & conquer separation principles (multiple sections with its own `proto_descriptor` fields and maybe route rules?)
- `/grpc` envoy route prefix failed to work, switched to `/` prefix, does it not trim like in nginx?

Additional Envoy's examples of Protobuf annotations you may use: 
https://github.com/envoyproxy/envoy/blob/3424968/test/proto/bookstore.proto

# Alternative

- https://github.com/grpc-ecosystem/grpc-gateway - generate HTTP bindings from Proto files
like you do from OpenAPI format. Might be dangerous because is not a black box but rather
fully extendable solution you can write the code for. But easier than Envoy + you don't need
a separate API layer and a technology to support.
- https://www.krakend.io/docs/enterprise/backends/grpc/ - paywall, enterprise only feature

# Links

- https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/grpc_json_transcoder_filter
- https://oneuptime.com/blog/post/2026-01-08-grpc-gateway-rest-transcoding/view
- https://stackoverflow.com/questions/40367413
- https://github.com/grpc-ecosystem/grpc-gateway/blob/main/examples/internal/proto/examplepb/a_bit_of_everything.proto