FROM golang:1.13-alpine AS build-env


ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOSUBDB=off \
    GOPROXY=direct \
    GOOS=linux
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH

RUN apk add --no-cache --update wget git gcc curl bash make


WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

FROM alpine
WORKDIR /app
COPY --from=build-env /app/build/comparator /app/comparator
COPY --from=build-env /app/config.json /app/config.json
EXPOSE 9091
CMD [comparator]