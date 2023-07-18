FROM golang:1.18.5-buster as build-env

LABEL maintainer="lukegoooo@gmail.com"

ENV GO111MODULE=off
WORKDIR /go/src/git.autops.xyz/autops/tulip

COPY . .

RUN make dist

FROM d.autops.xyz/base:3.12

COPY --from=build-env /go/src/git.autops.xyz/autops/tulip/tulip /app/tulip

# Run the app by default when the container starts
CMD /app/tulip
