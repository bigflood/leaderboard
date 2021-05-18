FROM golang:1.16.4 AS build

ENV CGO_ENABLED=0

COPY . /go/src/leaderboard/

WORKDIR /go/src/leaderboard

RUN go build  ./cmd/leaderboard-server

FROM alpine

COPY --from=build /go/src/leaderboard/leaderboard-server /

ENTRYPOINT ["/leaderboard-server"]
