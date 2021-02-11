FROM tuplestream/golang:latest AS build

WORKDIR /collector
ADD . .
RUN ./build.bash -release

FROM debian:stable-slim AS runtime

WORKDIR /runtime

COPY --from=build /collector/bin/collector .

CMD [ "/runtime/collector" ]
