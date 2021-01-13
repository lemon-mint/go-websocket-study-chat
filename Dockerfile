FROM golang:alpine as build

RUN apk update
RUN apk add git upx

ADD . /build
WORKDIR /build

RUN go build -ldflags="-s -w" -v -o wschat main.go
RUN upx --lzma /build/wschat


FROM alpine:latest
COPY --from=build /build /app
EXPOSE 54791
WORKDIR /app
CMD /app/wschat
