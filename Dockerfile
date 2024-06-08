FROM golang:1.22-bookworm as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o feedy

FROM gcr.io/distroless/static-debian12
COPY --from=builder /build/feedy /usr/bin/feedy
ENTRYPOINT [ "/usr/bin/feedy" ]
