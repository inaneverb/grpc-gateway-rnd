
# Keep in mind

### Single proto descriptor per route rule

`proto_descriptor` field in Envoy's YAML config accepts ***single*** string (filepath).
So that's why we either:
- Have one descriptor file from multiple (or all) Proto services (one big gen command)
- Have multiple sections of Envoy rules, each for their own Proto service and HTTP route (is it possible?)

Additional R&D required:
- Does the https://buf.build/ supports proto descriptor files? From multiple sources?
- Envoy's config divide & conquer separation principles (multiple sections with its own `proto_descriptor` fields and maybe route rules?)

Additional Envoy's examples of Protobuf annotations you may use: 
https://github.com/envoyproxy/envoy/blob/3424968/test/proto/bookstore.proto

# Alternative

- https://github.com/grpc-ecosystem/grpc-gateway - generate HTTP bindings from Proto files
like you do from OpenAPI format. Might be dangerous because is not a black box but rather
fully extendable solution you can write the code for. But easier than Envoy + you don't need
a separate API layer and a technology to support.

# Links

- https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/grpc_json_transcoder_filter
- https://oneuptime.com/blog/post/2026-01-08-grpc-gateway-rest-transcoding/view
- https://stackoverflow.com/questions/40367413
- https://github.com/grpc-ecosystem/grpc-gateway/blob/main/examples/internal/proto/examplepb/a_bit_of_everything.proto