FROM golang:1.14.2-alpine3.11 as go-builder

RUN apk add --no-cache gcc g++

WORKDIR $GOPATH/src/operametrix/mqtt

COPY go.mod go.sum ./

RUN go mod verify

COPY middleware middleware
COPY proxy proxy
COPY cmd cmd
COPY session session
COPY conf conf
COPY notify notify

RUN go build -o guardmq cmd/main.go

# Final stage
FROM alpine:3.11

WORKDIR "/usr/share/proxy"

LABEL maintainer="Nicolas Gonzalez <ngonzalez@operametrix.fr>"

RUN apk add --update openssl

COPY --from=go-builder /go/src/operametrix/mqtt/guardmq .
COPY --from=go-builder /go/src/operametrix/mqtt/conf ./conf/
ENV PATH="/usr/share/proxy:${PATH}"

EXPOSE 1883

ENTRYPOINT [ "./guardmq" ]
