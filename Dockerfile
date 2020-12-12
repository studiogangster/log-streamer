
FROM golang:1.13 as build


ARG CGO_ENABLED=0
ARG GO111MODULE="on"
ARG GOPROXY=""

WORKDIR /go/src/log_reader
COPY .  .


RUN CGO_ENABLED=0  go build   .

RUN ls
RUN chmod 777 ./log_reader

CMD   [ "./log_reader" ]


FROM alpine

COPY --from=build /go/src/log_reader/log_reader /log_reader

ENTRYPOINT    [ "/log_reader" ]