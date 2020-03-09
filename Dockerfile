FROM golang:1.14-buster AS build

WORKDIR /hawkeye-sidecar
RUN go get github.com/hpcloud/tail/...
ADD main.go .
RUN go build

FROM debian:stable-slim AS runtime

WORKDIR /runtime

COPY --from=build /hawkeye-sidecar/hawkeye-sidecar .

CMD [ "/runtime/hawkeye-sidecar" ]