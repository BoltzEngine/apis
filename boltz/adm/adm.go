// Package adm implements communication between master and slave for ADM.
package adm

// Credential はADMサーバにリクエストするための視覚情報をあらわす。
type Credential struct {
	TokenURL           string // OAuth2エンドポイント
	ClientID           string
	ClientSecret       string
	InsecureSkipVerify bool // (テスト用)サーバ証明書の正当性を確認をしない
}

// Message represents request message to send to ADM.
type Message struct {
	RegID string `json:"-"`

	Data             map[string]string `json:"data,omitempty"`
	ConsolidationKey string            `json:"consolidationKey,omitempty"`
	ExpiresAfter     int               `json:"expiresAfter,omitempty"` // second
	MD5              string            `json:"md5,omitempty"`          // checksum of Data
}

type Request struct {
	OriginURL  string // メッセージリクエストURL(scheme://host[:port])
	Credential *Credential
	Messages   []*Message
	BandWidth  int32 // 1秒あたりの通知数(0以下なら無制限)
}

type ProtocolError string

const (
	InvalidRegistrationIdError   ProtocolError = "InvalidRegistrationId"
	InvalidDataError                           = "InvalidData"
	InvalidConsolidationKeyError               = "InvalidConsolidationKey"
	InvalidExpirationError                     = "InvalidExpiration"
	InvalidChecksumError                       = "InvalidChecksum"
	InvalidTypeError                           = "InvalidType"
	UnregisteredError                          = "Unregistered"
	AccessTokenExpiredError                    = "AccessTokenExpired"
	MessageTooLargeError                       = "MessageTooLarge"
	MaxRateExceededError                       = "MaxRateExceeded"
)

var (
	invalidRequestFailures = map[ProtocolError]struct{}{
		InvalidDataError:             struct{}{},
		InvalidConsolidationKeyError: struct{}{},
		InvalidExpirationError:       struct{}{},
		InvalidChecksumError:         struct{}{},
		InvalidTypeError:             struct{}{},
		MessageTooLargeError:         struct{}{},
	}
	invalidTokenFailures = map[ProtocolError]struct{}{
		InvalidRegistrationIdError: struct{}{},
		UnregisteredError:          struct{}{},
	}
)

func (e ProtocolError) Error() string {
	return string(e)
}

func (e ProtocolError) Temporary() bool {
	if e.InvalidToken() {
		return false
	}
	_, ok := invalidRequestFailures[e]
	return !ok
}

func (e ProtocolError) InvalidToken() bool {
	_, ok := invalidTokenFailures[e]
	return ok
}

// FailedMessage は送信失敗したメッセージと失敗理由をあらわす。
// 必ず、ErrorString、Detail, RegIDはどれか1つだけセットされる。
// エラー判定を行ってからRegIDの確認をすることを推奨する。
type FailedMessage struct {
	// WebPushとは関係のない場所で発生したエラー(例えば"no such host")
	ErrorString string
	// ADMプロトコルにおけるエラーの場合にセット
	Detail *ProtocolError
	// 成功だけどトークン更新があった場合にセット
	RegID string
	// リクエストしたメッセージ
	Message *Message
}

func (r *FailedMessage) IsSuccess() bool {
	return r.ErrorString == "" && r.Detail == nil
}

func (r *FailedMessage) CanRetry() bool {
	if r.ErrorString != "" {
		return true
	}
	if r.Detail != nil {
		return r.Detail.Temporary()
	}
	return false
}

func (r *FailedMessage) UpdatedToken() (string, bool) {
	if !r.IsSuccess() {
		return "", false
	}
	if r.RegID == "" || r.RegID == r.Message.RegID {
		return "", false
	}
	return r.RegID, true
}

// Response はリクエストに対するスレーブからの応答をあらわす。
type Response struct {
	// 送信失敗またはトークン更新したメッセージと理由。
	// すべて成功した場合は空の配列。
	FailedMessages []*FailedMessage
}
