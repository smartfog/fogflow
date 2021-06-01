FROM golang:1.16 AS builder
WORKDIR /go/src/
ADD /discovery discovery/
ADD /common common/
ADD go.mod .
ADD go.sum .
WORKDIR /go/src/discovery
RUN go get;
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -a -installsuffix cgo  -o discovery

FROM scratch
COPY --from=builder /go/src/discovery/discovery  /
CMD ["/discovery"]
