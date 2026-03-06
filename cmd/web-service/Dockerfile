FROM rust:latest AS builder
WORKDIR /app

RUN --mount=type=cache,target=/usr/local/cargo/registry \
    true

COPY Cargo.toml Cargo.lock ./
COPY cmd/web-service/Cargo.toml ./cmd/web-service/

RUN mkdir -p cmd/web-service && echo "fn main() {}" > cmd/web-service/main.rs


RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/app/target \
    cargo build --manifest-path cmd/web-service/Cargo.toml --bin web-service --release && \
    rm -f target/release/deps/web_service*

COPY . .

RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/app/target \
    cargo build --manifest-path cmd/web-service/Cargo.toml --bin web-service --release && \
    cp target/release/web-service /app/web-service
#   ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
#   копируем из кэш-тома в обычный слой образа

FROM debian:bookworm-slim AS runtime
WORKDIR /app
COPY --from=builder /app/web/templates ./web/templates
COPY --from=builder /app/web-service ./web-service

CMD ["/app/web-service"]
