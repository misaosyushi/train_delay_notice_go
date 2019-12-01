# 電車遅延通知Line bot

## 対象の路線の遅延情報を取得し、Lineにpush通知するサービス

### How to use

- 通知したい運営会社、路線名（`targetCompany`/`targetName`）を設定
- goプロジェクトをビルドし、zipに固めてAWS Lambdaにデプロイ
 - バイナリのファイル名はハンドラと同じ名前にする
```
GOOS=linux GOARCH=amd64 go build -o hoge
```
- Messaging APIの`channel secret`, `channel token`, `user id`を環境変数にセット
- CloudWatch Eventsなどでcronを仕込めば定期実行が可能です
