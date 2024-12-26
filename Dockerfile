# https://klotzandrew.com/blog/smallest-golang-docker-image/
FROM golang:1.23.4 as base

RUN adduser \
    --disabled-password \
    --gecos "" \
    --no-create-home \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 59999 \
    small-user

WORKDIR /app

COPY . .

ENV GIN_MODE=release
ENV PORT=8080
ENV TZ=UTC

WORKDIR /app

RUN go mod download
RUN go mod verify

# RUN go build -o main .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /main .

FROM scratch

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=base /main .

USER small-user:small-user

CMD ["./main"]
