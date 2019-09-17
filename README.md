# BoltzEngine API 接続定義ファイル

## このリポジトリについて

[超高速プッシュ通知エンジン BoltzEngine](https://www.fenrir-inc.com/jp/boltzengine/) への接続定義ファイルを公開するリポジトリです。

## 使用方法

### gRPC でご利用の場合

本リポジトリの rpc/ 階層以下にある proto ファイルを protoc コマンドで各言語のクライアントコードに変換してご利用いただけます。

### Go から gRPC でご利用の場合

1. [protoc](https://github.com/protocolbuffers/protobuf/releases) をインストールしてください。
開発環境が Mac の場合は `brew install protobuf` をご利用頂けます。
開発環境が Debian / Ubuntu の方は `apt-get install protobuf-compiler` をご利用頂けます。
2. `protoc-gen-go` をインストールしてください。 `go get -u github.com/golang/protobuf/protoc-gen-go` でインストール可能です。
3. `make` を実行すると、 `*.pb.go` ファイルが最新に更新されます。

#### ProtoPackageIsVersion3 のエラーが出る場合
v3.2.0 より、リポジトリに一緒に登録される *.pg.go ファイルたちは Protocol Buffers v3.0 で生成されています。
こちらを用いると

```
xxxxx.pb.go:xx:xx: undefined: proto.ProtoPackageIsVersion3
```

のようなエラーが出力される場合があります。
その場合は、以下の Issue が参考になるかも知れません。

* [Rev of proto-gen-go to ProtoPackageIsVersion3 causing breakage](https://github.com/golang/protobuf/issues/763)

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

## License

```
Copyright 2018 Fenrir Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
