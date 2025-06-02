FROM  golang:1.24.2-alpine3.21 

WORKDIR /go/src/app/

# back以下の中身だけを /go/src/app/ 配下にコピー
COPY ./. .

RUN apk upgrade --update && \
    apk --no-cache add git gcc musl-dev

RUN go install github.com/go-delve/delve/cmd/dlv@latest 