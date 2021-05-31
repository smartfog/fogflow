FROM golang:1.16 AS builder
WORKDIR /go/src/
ADD /worker worker/
ADD /common common/
ADD go.mod .
ADD go.sum .
WORKDIR /go/src/worker
RUN go get;
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -a -installsuffix cgo  -o worker

FROM scratch
COPY --from=builder /go/src/worker/worker  /
CMD ["/worker"]
