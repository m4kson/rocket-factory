FROM golang:1.25.5 AS builder

RUN apt-get update \
    && apt-get install -y git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY inventory/go.mod inventory/go.sum ./inventory/

ENV GOWORK=off
RUN cd inventory && go mod download

COPY platform ./platform
COPY shared ./shared
COPY inventory ./inventory

ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

RUN chmod +x grpc-health-probe

RUN cd inventory && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/app-inventory ./cmd/main.go

FROM alpine:3.21.3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

COPY --from=builder /app/app-inventory .

COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

EXPOSE 50051

ENTRYPOINT ["./app-inventory"]