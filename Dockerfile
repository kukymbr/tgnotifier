FROM golang:1.23-alpine AS gobuild

RUN apk add --update --no-cache ca-certificates git

RUN mkdir -p /src
WORKDIR /src

# Download the dependencies as a separated layer to cache it separately.
COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../ .
RUN go build ./cmd/tgnotifier

FROM alpine:3.8

RUN mkdir /tgnotifier
WORKDIR /tgnotifier
COPY --from=gobuild /src/tgnotifier ./

ENTRYPOINT ["/tgnotifier/tgnotifier"]