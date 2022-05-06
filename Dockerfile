FROM golang:1.18 as builder

WORKDIR /src

COPY . .

RUN go build -o yv cmd/yv/*.go

FROM alpine:3.12 as runner

COPY --from=builder /src/yv /bin/yv

ENTRYPOINT [ "/bin/yv" ]