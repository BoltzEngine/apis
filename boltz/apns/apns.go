// Package apns implements communication between master and slave.
package apns

import (
	"fmt"
	"strings"
	"time"
)

const (
	PrioritySentImmediately   uint8 = 10
	PrioritySentAtPowerSaving       = 5
)

const (
	PushTypeAlert        = "alert"
	PushTypeBackground   = "background"
	PushTypeVoIP         = "voip"
	PushTypeComplication = "complication"
	PushTypeFileProvider = "fileprovider"
	PushTypeMdm          = "mdm"
)

type Status uint8

const (
	Success            Status = 0
	ProcessingError           = 1
	MissingToken              = 2
	MissingTopic              = 3
	MissingPayload            = 4
	InvalidTokenSize          = 5
	InvalidTopicSize          = 6
	InvalidPayloadSize        = 7
	InvalidToken              = 8
	Shutdown                  = 10
	None                      = 255
)

func (status Status) String() string {
	switch status {
	case ProcessingError:
		return "apns: processing error"
	case MissingToken:
		return "apns: missing device token"
	case MissingTopic:
		return "apns: missing topic"
	case MissingPayload:
		return "apns: missing payload"
	case InvalidTokenSize:
		return "apns: invalid token size"
	case InvalidTopicSize:
		return "apns: invalid topic size"
	case InvalidPayloadSize:
		return "apns: invalid payload size"
	case InvalidToken:
		return "apns: invalid token"
	case Shutdown:
		return "apns: shutdown"
	case None:
	default:
		return "apns: none (unknown)"
	}
	return ""
}

// CredentialはAPNs接続時の資格情報をあらわす。
type Credential struct {
	// APNs証明書の秘密鍵
	KeyPEMBlock []byte
	// APNs証明書の証明書
	CertPEMBlock []byte

	// JWTのissキー
	Issuer string
	// JWTのkidキー
	KeyID string
	// JWTのPEMエンコードされたEC P-256秘密鍵
	PrivateKey []byte

	// (テスト用)接続先サーバの証明書を検査しない
	InsecureSkipVerify bool
}

// Requestはマスタからスレーブに対してリクエストするメッセージをあらわす。
type Request struct {
	// APNsホストアドレス(大抵はgateway.push.apple.com:2195)
	Addr string
	// APNs接続時の資格情報
	Credential *Credential
	// 送信メッセージ
	Messages []*Message
	// 1秒あたりの通知数(0以下なら無制限)
	BandWidth int32
}

type Message struct {
	ID    uint32
	Expir uint32
	Token []byte

	// json encoded binary
	Payload []byte

	Priority uint8

	// APNs HTTP/2 only
	Topic      string
	CollapseID string

	// for iOS 13~, watchOS 6~ (APNs HTTP/2 only
	PushType string
}

// FailedMessageは送信失敗したメッセージとその理由をあらわす。
// 必ず、ErrorString、Detailはどれか1つだけセットされる。
type FailedMessage struct {
	// APNsとは関係のない場所で発生したエラー(例えば"no such host")
	ErrorString string
	// APNsプロトコルにおけるエラーの場合にセット
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

type ProtocolError struct {
	// Binary protocol
	Cmd    uint8
	Status Status
	ID     uint32

	// HTTP/2
	StatusCode int
	Reason     string
	Time       time.Time
}

func (e *ProtocolError) Error() string {
	switch e.StatusCode {
	case 0: // Legacy
		return e.Status.String()
	default:
		return fmt.Sprintf("%d %s", e.StatusCode, e.Reason)
	}
}

func (e *ProtocolError) InvalidToken() (ret bool) {
	if e.StatusCode > 0 {
		// APNs + HTTP/2
		ret = isInvalidToken(e.StatusCode, e.Reason)
	} else {
		// Legacy
		ret = e.Status == InvalidToken
	}
	return
}

func (e *ProtocolError) Timestamp() time.Time {
	if !e.InvalidToken() {
		return time.Time{}
	}
	switch e.StatusCode {
	case 0: // Legacy
		return time.Now()
	default:
		// 404の場合はtimestampが設定されないが、
		// 有用な値を返すべきなので現在時刻を返す。
		if e.Time.IsZero() {
			return time.Now()
		}
		return e.Time
	}
}

type FBRequest struct {
	Addr       string
	Credential *Credential
}

type FBResponse struct {
	Body []*Feedback
}

type Feedback struct {
	Timestamp uint32
	Token     []byte
}

func IsLegacyAddr(addr string) bool {
	if strings.HasPrefix(addr, "https://") {
		return false
	}
	if strings.HasPrefix(addr, "http://") {
		return false
	}
	return true
}
