FROM golang:alpine AS builder

RUN apk --update add --no-cache --virtual .build-deps \
      ca-certificates \
      tzdata \
    && update-ca-certificates \
    && rm -rf /tmp/*.apk /var/cache/apk/*

# Create unprivileged user
ENV USER=appuser
ENV UID=1000

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --no-create-home \
    --home "/nonexistent" \
    --gecos "" \
    --shell "/sbin/nologin" \
    --uid "${UID}" \
    "${USER}"

COPY . /src/

WORKDIR /src

RUN CGO_ENABLED=0 go build -ldflags '-w -s -extldflags "-static"' -o /src/api


FROM scratch

LABEL maintainer="Jordan Cartwright"

# Copy from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Add the static executable
COPY --from=builder /src/api /api

# Use an unprivileged user
USER appuser:appuser

ENTRYPOINT ["/api"]

CMD []

EXPOSE 8080
