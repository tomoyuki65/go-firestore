FROM golang:1.26.1-alpine3.23

WORKDIR /go/src

COPY ./src .

# go.modがあれば依存関係をインストール
RUN if [ -f ./go.mod ]; then \
      go install; \
    fi

# 開発用のライブラリをインストール
RUN go install github.com/air-verse/air@v1.64.5
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

EXPOSE 8081