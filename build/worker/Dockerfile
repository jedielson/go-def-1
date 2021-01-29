FROM golang:1.15-alpine AS base 
RUN apk add --no-cache ca-certificates git make

FROM base AS dependencies 
WORKDIR /workspace
COPY go.* ./
COPY Makefile .
RUN make deps

FROM dependencies AS builder
ENV GOOS=linux
ENV GOARCH=amd64
COPY . .
RUN make build

FROM scratch as runner 
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /workspace/bin/worker /bin/worker
ENTRYPOINT ["./bin/worker"]