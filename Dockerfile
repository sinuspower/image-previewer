FROM golang:1.15-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/image-previewer .

FROM scratch
WORKDIR /root/
COPY --from=builder /app/bin/image-previewer ./bin/image-previewer
CMD ["./bin/image-previewer"]