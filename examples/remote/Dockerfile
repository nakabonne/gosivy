FROM golang:1.17-alpine AS build-env

COPY . /app
WORKDIR /app/examples/remote

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -o remote \
 && mv ./remote /


FROM alpine

COPY --from=build-env /remote /usr/local/bin/

EXPOSE 9090

CMD ["remote"]
