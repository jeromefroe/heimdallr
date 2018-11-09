FROM golang:1.11 AS build
WORKDIR /go/src/github.com/jeromefroe/heimdallr

RUN go get github.com/twitchtv/retool
COPY Makefile tools.json Gopkg.toml Gopkg.lock ./
RUN make dep-install

COPY cmd cmd
COPY pkg pkg
RUN CGO_ENABLED=0 make build

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/heimdallr /bin/heimdallr
