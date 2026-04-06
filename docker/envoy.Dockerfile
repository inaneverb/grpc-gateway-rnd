FROM envoyproxy/envoy:dev

COPY /config/envoy /envoy-custom

RUN ls /envoy-custom

CMD ["-c /envoy-custom/json-grpc-transcode.yaml"]