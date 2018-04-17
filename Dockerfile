FROM ysitd/dep AS builder

WORKDIR /go/code.ysitd.cloud/component/art/gallery

COPY .  /go/code.ysitd.cloud/component/art/gallery

RUN dep ensure -vendor-only && \
    go build -v main.go

FROM alpine:3.6

COPY --from=builder /go/code.ysitd.cloud/component/art/gallery/gallery /

CMD ["/gallery"]
