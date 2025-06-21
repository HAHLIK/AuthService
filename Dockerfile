FROM golang:1.24-alpine AS builder

WORKDIR /usr/local/src

COPY src/go.mod  ./

RUN go mod download

COPY src/ ./
RUN go build -o bin/app sso/cmd/sso/main.go

FROM alpine

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/sso/env/.env usr/local/src/sso/config/local.yaml /

CMD [ "/app" ]