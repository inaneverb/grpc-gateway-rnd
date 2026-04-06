FROM alpine:3.20 AS builder

RUN apk add --no-cache \
    curl \
    bash \
    git \
    build-base

RUN curl https://mise.run | sh

ENV PATH="/root/.local/bin:$PATH"

WORKDIR /app

COPY . .

RUN mise trust --yes || true
RUN mise install
RUN mise exec -- just build

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/upstream /app/upstream

RUN chmod +x /app/upstream

ENTRYPOINT ["/app/upstream"]