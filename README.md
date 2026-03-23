# Go言語からCloud Firestoreを使うためのサンプル
Docker環境でGo言語（Golang）からCloud Firestoreを使う方法についてのサンプルです。  
  
<br>
  
## 動作要件
- Go: 1.26.1
- Echo: v4.15.0
- firebase: 15.11.0
  
<br>
  
## ローカル開発環境構築
### 1. 環境変数ファイルをリネーム
```
cp ./.env.example ./.env
```  
  
### 2. コンテナのビルドと起動
```
docker compose build --no-cache
docker compose up -d
```  
  
### 3. コンテナの停止・削除
```
docker compose down
```  
> ※ボリュームも合わせて削除したい場合は、オプション「-v」を付けて実行して下さい。（例：docker compose down -v）  
  
<br>
  
## Firebase Emulator Suiteの開き方
ローカルサーバー起動後、ブラウザで「http://localhost:4000」を開く  
  
<br>
  
## コード修正後に使うコマンド
ローカルサーバー起動中に以下のコマンドを実行可能です。  
  
### 1. go.modの修正
```
docker compose exec api go mod tidy
```  
  
### 2. フォーマット修正
```
docker compose exec api go fmt ./...
```  
  
### 3. コード解析チェック
```
docker compose exec api staticcheck ./...
```  
  
<br>
  
## サンプルAPIのエンドポイント
Base URL: http://localhost:8081  
  
- GET /  
ルート（Hello World）  
  
- POST /api/v1/messages  
  メッセージを作成  
  
- GET /api/v1/messages  
  全てのメッセージを取得  
  
- GET /api/v1/messages/:id  
  対象メッセージ取得  
  
- PUT /api/v1/messages/:id  
  対象メッセージ更新  
  
- DELETE /api/v1/messages/:id  
  対象メッセージ削除
  
<br>
  
## 参考記事  
[・Go言語（Golang）でGoogle Cloud Firestore（NoSQL）の使い方｜Docker環境構築＋CRUD APIサンプル](https://golang.tomoyuki65.com/golang-firestore-docker-crud-api)  
  