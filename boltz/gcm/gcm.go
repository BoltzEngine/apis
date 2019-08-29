// Package gcm implements communication between master and slave for FCM.
package gcm

import (
	"errors"
	"net/url"
	"strings"
)

const (
	// 1つのメッセージに含められる最大registration idの個数
	RegIDsMax = 1000
)

// APIErrorはFCM HTTP v1 APIのエラーをあらわす。
type APIError struct {
	Message *Message // オリジナルのメッセージ

	Code        int
	Status      string
	Description string
}

func (e *APIError) ErrorCode() string {
	switch e.Status {
	case "NOT_FOUND",
		"PERMISSION_DENIED",
		"RESOURCE_EXHAUSTED",
		"UNAUTHENTICATED",
		"APNS_AUTH_ERROR",
		"INTERNAL",
		"INVALID_ARGUMENT",
		"SENDER_ID_MISMATCH",
		"QUOTA_EXCEEDED",
		"UNAVAILABLE",
		"UNREGISTERED":
		return e.Status
	default:
		return "INTERNAL"
	}
}

var (
	invalidParamCodes = map[string]struct{}{
		"PERMISSION_DENIED":  struct{}{},
		"UNAUTHENTICATED":    struct{}{},
		"APNS_AUTH_ERROR":    struct{}{},
		"INVALID_ARGUMENT":   struct{}{},
		"SENDER_ID_MISMATCH": struct{}{},
	}

	invalidTokenCodes = map[string]struct{}{
		"NOT_FOUND":    struct{}{},
		"UNREGISTERED": struct{}{},
	}
)

func (e *APIError) Temporary() bool {
	_, invalidParam := invalidParamCodes[e.ErrorCode()]
	return !e.BadRegID() && !invalidParam
}

func (e *APIError) BadRegID() bool {
	_, ok := invalidTokenCodes[e.ErrorCode()]
	return ok
}

// FailureはFCMサーバから発生されるエラーをあらわす。
type Failure string

const (
	// リクエストにRegIDが含まれていることの確認が必要。
	FailureMissingRegistration = Failure("MissingRegistration")
	// FCMサーバに送信したRegIDフォーマットの確認が必要。
	FailureInvalidRegistration = Failure("InvalidRegistration")
	// RegIDが特定のセンダーのグループに縛られている。
	FailureMismatchSenderId = Failure("MismatchSenderId")
	// RegIDがなんらかの理由で無効になった。
	FailureNotRegistered = Failure("NotRegistered")
	// メッセージに含まれるペイロードデータが大きすぎる。
	FailureMessageTooBig = Failure("MessageTooBig")
	// Googleにより予約済みキー名がペイロードのキーとして使われている。
	FailureInvalidDataKey = Failure("InvalidDataKey")
	// TimeToLiveが大きすぎる、または0以下。
	FailureInvalidTtl = Failure("InvalidTtl")
	// FCMサーバでタイムアウトした。
	FailureUnavailable = Failure("Unavailable")
	// FCMサーバになんらかのエラーが発生した。
	FailureInternalServerError = Failure("InternalServerError")
	// 無効なパッケージ名。
	FailureInvalidPackageName = Failure("InvalidPackageName")
	// 特定デバイスへのメッセージレート超過。
	FailureDeviceMessageRateExceeded = Failure("DeviceMessageRateExceeded")
	// 特定トピックへのメッセージレート超過。
	FailureTopicMessageRateExceeded = Failure("TopicMessageRateExceeded")
)

// BadRegIDはfがRegIDに問題があるエラーの場合にtrueを返す。
func (f Failure) BadRegID() bool {
	return f == FailureInvalidRegistration || f == FailureNotRegistered
}

// Temporaryはfが一時的なエラーである場合にtrueを返す。
func (f Failure) Temporary() bool {
	return f == FailureUnavailable || f == FailureInternalServerError || f == FailureDeviceMessageRateExceeded || f == FailureTopicMessageRateExceeded
}

const (
	BadAck                    = "BAD_ACK"
	BadRegistration           = "BAD_REGISTRATION"
	ConnectionDraining        = "CONNECTION_DRAINING"
	DeviceMessageRateExceeded = "DEVICE_MESSAGE_RATE_EXCEEDED"
	DeviceUnregistered        = "DEVICE_UNREGISTERED"
	InternalServerError       = "INTERNAL_SERVER_ERROR"
	InvalidJSON               = "INVALID_JSON"
	ServiceUnavailable        = "SERVICE_UNAVAILABLE"
	TopicMessageRateExceeded  = "TOPICS_MESSAGE_RATE_EXCEEDED"
)

var (
	temporaryFailures = map[string]bool{
		ConnectionDraining:        true,
		DeviceMessageRateExceeded: true,
		InternalServerError:       true,
		ServiceUnavailable:        true,
		TopicMessageRateExceeded:  true,
	}
	invalidTokenFailures = map[string]bool{
		BadRegistration:    true,
		DeviceUnregistered: true,
	}
	errTimeout = errors.New("timed out waiting for ack/nack")
)

// Signal はCCSからのエラーを表す。
type Signal struct {
	Message *Message // オリジナルのメッセージ

	// ack only: canonical registration ID
	RegID string

	// nack only: some errors
	Code        string // CCSからのエラーコード
	Description string // CCSからのエラーメッセージ
}

func (p *Signal) Error() string {
	s := p.Description
	if s == "" {
		s = p.Code
	}
	return s
}

func (p *Signal) BadRegID() bool {
	return invalidTokenFailures[p.Code]
}

func (p *Signal) Temporary() bool {
	return temporaryFailures[p.Code]
}

func (p *Signal) CausedMessage() *Message {
	return p.Message
}

// Message represents request message to send to FCM
type Message struct {
	ID               string            `json:"message_id,omitempty"`
	RegIDs           []string          `json:"registration_ids,omitempty"`
	To               string            `json:"to,omitempty"`
	Data             map[string]string `json:"data,omitempty"`
	CollapseKey      string            `json:"collapse_key,omitempty"`
	DelayWhileIdle   bool              `json:"delay_while_idle,omitempty"`
	TimeToLive       int               `json:"time_to_live,omitempty"`
	Priority         string            `json:"priority,omitempty"`
	ContentAvailable bool              `json:"content_available,omitempty"`
	MutableContent   bool              `json:"mutable_content,omitempty"`
	DryRun           bool              `json:"dry_run,omitempty"`
	Notification     map[string]string `json:"notification,omitempty"`
}

// CredentialはFCMサーバにアクセスする資格情報をあらわす。
type Credential struct {
	// APIが有効なservice-account.jsonの内容
	ServiceAccount string

	// FCMサーバキー
	ServerKey string
	// センダーID
	SenderID string

	// (テスト用)サーバ証明書の正当性を確認をしない
	InsecureSkipVerify bool
}

// Requestはマスタからスレーブに対してリクエストするメッセージをあらわす。
type Request struct {
	// FCMサーバのURL
	URL string
	// FCMサーバへアクセスするための資格情報
	Credential *Credential
	// FCMサーバへ送信するメッセージ
	Messages []*Message
	// 1秒あたりの通知数(0以下なら無制限)
	BandWidth int32
}

// FailedMessageは送信失敗したメッセージとその理由をあらわす。
// 必ず、ErrorString、APIError、Detail、Signalはどれか1つだけセットされる。
// なのでDetailまたはSignalを判定し、nilならErrorStringをエラーの理由として扱うこと。
type FailedMessage struct {
	// FCMとは関係のない場所で発生したエラー(例えば"no such host")
	ErrorString string
	// FCM HTTP v1 APIプロトコルにおけるエラーの場合にセット
	Error *APIError
	// FCM-HTTPプロトコルにおけるエラーの場合にセット
	Detail *ResponseBody
	// FCM-XMPPプロトコルにおけるエラーまたは登録ID更新の場合にセット
	Signal *Signal
	// 失敗を引き起こしたメッセージ
	Message *Message
}

// ResponseBody represents response message from FCM
type ResponseBody struct {
	ID           int64    `json:"multicast_id"`
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIDs int      `json:"canonical_ids"`
	Results      []Result `json:"results"`
	RetryCount   int
}

// ResultはFCMサーバからのRegIDひとつに対する結果をあらわす。
type Result struct {
	// 送信ID。送信成功の場合に値が入る。
	ID string `json:"message_id,omitempty"`
	// 送信成功したが該当するRegIDが更新されていた場合に新しいRegIDが入る。
	RegID string `json:"registration_id,omitempty"`
	// 送信失敗の場合に、エラーの理由をセットする。
	Error Failure `json:"error,omitempty"`
}

func (r Result) IsSuccess() bool {
	return r.ID != ""
}

// IsGarbageRegIDは今後送信すべきではないRegIDの場合にtrueを返す。
func (r Result) IsGarbageRegID() bool {
	return !r.IsSuccess() && r.Error.BadRegID()
}

// IsCanonicalIDは該当するRegIDに更新が発生している場合にtrueを返す。
func (r Result) IsCanonicalID() bool {
	return r.RegID != ""
}

// Responseはリクエストに対するスレーブからの応答をあらわす。
type Response struct {
	// 送信失敗したメッセージと理由。
	// すべて成功した場合は空の配列。
	FailedMessages []*FailedMessage
}

const (
	GcmEndpoint = iota - 2
	FcmEndpointError
	FcmXmppApi
	FcmLegacyHttpApi
	FcmHttpV1Api
)

// ProtocolVersion determines the protocol to be used by the argument addr.
//  - HTTP v1 API protocol or higher. : FcmHttpV1Api
//  - Legacy HTTP protocol.           : FcmLegacyHttpApi
//  - Legacy XMPP protocol.           : FcmXmppApi
//  - Using GCM endpoint.             : FcmEndpointError [error]
//  - Incorrect endpoint.             : GcmEndpoint      [error]
func ProtocolVersion(addr string) int {
	if u, err := url.Parse(addr); err != nil {
		return FcmEndpointError
	} else {
		switch {
		case len(u.Host) == 0 && strings.Contains(u.Scheme, "fcm"),
			(u.Scheme == "http" || u.Scheme == "https") && strings.Contains(u.Host, "fcm"):
			switch {
			case strings.Contains(u.Path, "/fcm/") && len(u.Opaque) == 0:
				return FcmLegacyHttpApi
			case strings.Contains(u.Path, "/v1/") && len(u.Opaque) == 0:
				return FcmHttpV1Api
			case strings.Contains(u.Scheme, "xmpp") && len(u.Opaque) > 0:
				return FcmXmppApi
			}
		case strings.Contains(u.Scheme, "gcm"), strings.Contains(u.Host, "gcm"):
			return GcmEndpoint
		}
	}
	return FcmEndpointError
}
