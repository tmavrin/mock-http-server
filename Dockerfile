FROM golang:1.22-alpine AS build-env
RUN apk --no-cache add git
WORKDIR /go/src/github.com/tmavrin/mock-http-server/
RUN ls
COPY . .
RUN go build -o service


FROM alpine
RUN apk update && apk add ca-certificates
WORKDIR /app
COPY ./config.json /app/
COPY --from=build-env /go/src/github.com/tmavrin/mock-http-server/service /app/
ENTRYPOINT /app/service