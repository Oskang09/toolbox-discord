FROM golang:1.17.11 as builder
COPY . /go/toolbox-discord
WORKDIR /go/toolbox-discord
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags timetzdata -ldflags="-w -s" -o main .

FROM alpine
COPY --from=builder /go/toolbox-discord /go/toolbox-discord
WORKDIR /go/toolbox-discord
ENTRYPOINT ./main