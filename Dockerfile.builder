
FROM golang:1.13 as build


ARG CGO_ENABLED=0
ARG GO111MODULE="on"
ARG GOPROXY=""

WORKDIR /go/src/readlogs
COPY .  .


RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o readlogs . \
    && CGO_ENABLED=0 GOOS=darwin go build -a -ldflags "-s -w" -installsuffix cgo -o readlogs-darwin  \
    && GOARM=6 GOARCH=arm CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o readlogs-armhf . \
    && GOARCH=arm64 CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o readlogs-arm64 . \
    && GOOS=windows CGO_ENABLED=0 go build -a -ldflags "-s -w" -installsuffix cgo -o readlogs.exe .

FROM scratch

ARG PLATFORM
COPY --from=build /go/src/readlogs/readlogs$PLATFORM ./fwatchdog
