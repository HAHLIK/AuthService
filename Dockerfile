FROM golang:1.24-alpine AS builder

WORKDIR /usr/local/src

COPY src/go.mod  ./

RUN go mod download

COPY src/ ./
RUN go build -o bin/app sso/cmd/sso/main.go

COPY config/ env/ ./

FROM alpine

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/.env usr/local/src/local.yaml /

CMD [ "/app" ]