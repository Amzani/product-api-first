FROM golang:1.10.3 as builder

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

# build directories
RUN mkdir -p /go/src/github.com/Amzani/product-api-first
ADD . /go/src/github.com/Amzani/product-api-first
WORKDIR /go/src/github.com/Amzani/product-api-first

# Go glide!
RUN curl https://glide.sh/get | sh
RUN glide update

# Build my app
RUN go build -o /go/src/github.com/Amzani/product-api-first/bin/server .

## Now run it.
FROM golang:1.10.3
RUN mkdir /app

COPY --from=builder /go/src/github.com/Amzani/product-api-first/bin/server /app/

WORKDIR /app
EXPOSE 5000
CMD ["./server"]