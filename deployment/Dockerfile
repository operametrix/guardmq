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

RUN go build cmd/main.go

# Final stage
FROM alpine:3.11

WORKDIR "/usr/share/proxy"

LABEL maintainer="Nicolas Gonzalez <ngonzalez@operametrix.fr>"

COPY --from=go-builder /go/src/operametrix/mqtt/main ./bin/
COPY --from=go-builder /go/src/operametrix/mqtt/conf ./conf/

EXPOSE 1883

ENTRYPOINT [ "./bin/main" ]
