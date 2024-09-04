FROM --platform=$BUILDPLATFORM golang:1.22-alpine as builder

RUN apk update && apk add --no-cache git ca-certificates tzdata 

ARG TARGETOS
ARG TARGETARCH

COPY . .

RUN go mod download

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags="-w -s" -o /mock-http-server ./cmd/main.go

FROM scratch AS final

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY ./example/mock-server.json /config/mock-server.json
COPY --from=builder /mock-http-server /mock-http-server

WORKDIR /

ENTRYPOINT ["/mock-http-server"]