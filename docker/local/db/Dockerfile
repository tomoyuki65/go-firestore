FROM ubuntu:26.04

WORKDIR /app

# 必要パッケージのインストール
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      curl \
      gnupg \
      openjdk-21-jre-headless \
      ca-certificates && \
    # Node.js v24をインストール
    curl -fsSL https://deb.nodesource.com/setup_24.x | bash - && \
    apt-get update && \
    apt-get install -y --no-install-recommends nodejs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Firebase CLIのインストール
RUN npm install -g firebase-tools

# 設定ファイルのコピー
COPY firebase.json .

EXPOSE 4000 8080 9150
