FROM golang:1.20-alpine3.18 as build
WORKDIR /build
ENV \
    TERM=xterm-color \
    TIME_ZONE="UTC" \
    CGO_ENABLED=1 \
    GOFLAGS="-mod=vendor"

RUN echo "## Prepare deps" && \
    apk add --no-cache --update tzdata gcc libc-dev && \
    cp /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime && \
    echo "${TIME_ZONE}" > /etc/timezone && date

COPY go.* ./
COPY bin/run/main.go ./main.go
COPY vendor/ ./vendor
COPY pkg ./pkg

RUN go env
RUN go version
RUN echo "  ## Build" && go build -o app .

FROM alpine:3.18
WORKDIR /app
COPY --from=build /build/app ./app
COPY --from=build /etc/localtime /etc/localtime
USER nobody:nobody
CMD ["./app"]