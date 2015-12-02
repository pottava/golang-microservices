Golang MicroServices
---

## Description

Go言語による マイクロサービスサンプル

## Micro Services

Name               | Container Name  | Port (dev only)  | Description
------------------ | --------------- | ---------------- | ------------------
app-authentication | auth            | 8080             | OAuth認証
app-aws            | aws             | (8081)           | AWSリソース操作 API
app-dbio           | dbio            | (8082)           | データベースIO API
app-webui          | web             | 80               | Web UI
-                  | dynamodb        |                  | DynamoDB Local

## Usage

### 1. 下準備
　  
AWS, Twitterでアプリ用 Credentialsを取得  
https://apps.twitter.com/app/  

### 2. Docker Quickstart Terminal を起動

### 3. 環境変数をセット

```
export AWS_REGION=ap-northeast-1
export AWS_ACCESS_KEY_ID=?
export AWS_SECRET_ACCESS_KEY=?
export APP_TWITTER_CONSUMER_KEY=?
export APP_TWITTER_CONSUMER_SECRET=?
export APP_TWITTER_CONSUMER_CALLBACK=http://192.168.99.100:8001/twitter/callback
```

### 4. git cloneしたフォルダに移動し、コンテナを起動

```
cd ~/src/github.com/pottava/golang-micro-services
docker-compose --x-networking up -d
```

### 5. ブラウザでアクセス（IPアドレスは VMのものを指定）

[http://192.168.99.100/](http://192.168.99.100/)

### 開発するときは

```
cd app-webui/
npm install
gulp
```
