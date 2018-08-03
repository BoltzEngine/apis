# BoltzEngine API 接続定義ファイル

## このリポジトリについて

[超高速プッシュ通知エンジン BoltzEngine](https://www.fenrir-inc.com/jp/boltzengine/) への接続定義ファイルを公開するリポジトリです。

## 使用方法

### gRPC でご利用の場合

本リポジトリの rpc/ 階層以下にある proto ファイルを protoc コマンドで各言語のクライアントコードに変換してご利用いただけます。

### Go から net/rpc 接続でご利用の場合

[![GoDoc](https://godoc.org/github.com/BoltzEngine/apis/boltz?status.svg)](https://godoc.org/github.com/BoltzEngine/apis/boltz)

本リポジトリの boltz/ 階層以下に型定義が入っており、go get で取得してご利用いただけます。

```
go get github.com/BoltzEngine/apis/boltz
```

## BoltzEngine について

**1秒で3.5万デバイスへ。降り注ぐ超高速プッシュ通知**

BoltzEngine（ボルツエンジン）は、国内トップクラスの配信速度を誇るプッシュ通知エンジンです。
オンプレミスにもクラウドにも対応する柔軟さで、自社の情報セキュリティポリシーに即した運用も可能です。
もちろん、iOS / Android の両方にメッセージを配信できます。

詳細は BoltzEngine のウェブサイトをご覧ください。

