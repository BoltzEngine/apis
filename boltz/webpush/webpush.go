// Package gcm implements communication between master and slave for WebPush.
package webpush

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// VAPID represents Voluntary Application Server Identification for Web Push.
type VAPID struct {
	Subject    string // トークン生成者のURI(mailto: or https:)
	PublicKey  []byte
	PrivateKey []byte
}

// CredentialはWebPushサーバにアクセスする資格情報をあらわす。
type Credential struct {
	// VAPID
	VAPID *VAPID
	// (テスト用)サーバ証明書の正当性を確認をしない
	InsecureSkipVerify bool
}

// Token represents BoltzEngine specific WebPush token.
type Token struct {
	Version   int    `json:"v"`
	URL       string `json:"endpoint"`
	PublicKey string `json:"p256dh"`
	AuthToken string `json:"auth"`
}

func (token *Token) String() (string, error) {
	s, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

type Urgency int

const (
	VeryLow Urgency = iota
	Low
	Normal
	High
)

type Message struct {
	Token      *Token
	Payload    string
	TimeToLive int // in second
	Urgency    Urgency
	Topic      string // like FCM's collapse_key
}

type Request struct {
	Credential *Credential
	Messages   []*Message
	BandWidth  int32 // 1秒あたりの通知数(0以下なら無制限)
}

type ProtocolError struct {
	StatusCode int
}

func (e *ProtocolError) InvalidToken() bool {
	return e.StatusCode == http.StatusNotFound || e.StatusCode == http.StatusGone
}

// FailedMessageは送信失敗したメッセージとその理由をあらわす。
// 必ず、ErrorString、Detailはどれか1つだけセットされる。
// なのでDetailを判定し、nilならErrorStringをエラーの理由として扱うこと。
type FailedMessage struct {
	// WebPushとは関係のない場所で発生したエラー(例えば"no such host")
	ErrorString string
	// WebPushプロトコルにおけるエラーの場合にセット
	Detail *ProtocolError
	// リクエストしたメッセージ
	Message *Message
}

// Responseはリクエストに対するスレーブからの応答をあらわす。
type Response struct {
	// 送信失敗したメッセージと理由。
	// すべて成功した場合は空の配列。
	FailedMessages []*FailedMessage
}

func LoadVAPIDKey(filename string) ([]byte, error) {
	v, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, base64.RawURLEncoding.DecodedLen(len(v)))
	n, err := base64.RawURLEncoding.Decode(buf, v)
	return buf[:n], err
}
