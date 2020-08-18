FROM golang:1.15-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/image-previewer .

FROM alpine/make
WORKDIR /root/
COPY --from=builder /app/bin/image-previewer ./bin/image-previewer
COPY --from=builder /app/Makefile .
CMD ["make", "run"]