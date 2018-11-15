// Package boltz defines data types for communicate between master and slave.
package boltz

import (
	"time"
)

// MasterActivity represents state of BoltzEngine master service.
type MasterActivity struct {
	// スレーブに対して送信待ちしているリクエスト数。
	PendingCount int
	// アドレスをキーとする各スレーブの状態。
	SlaveActivities map[string]SlaveActivity
	// ライセンスの期限
	Expiration time.Time
}

// SlaveActivity represents state of BoltzEngine slave service.
type SlaveActivity struct {
	// マスタひとつにつき最大いくつの同時リクエストを許可するか
	MaxAgents int

	// リクエストされた数
	RequestCount int
	// 現在送信中の端末数
	DeliveringCount int
	// 送信した端末数
	DeliveredCount int
	// リクエスト開始から終了までの合計時間
	TotalExecutionTime time.Duration
	// 最新のリクエスト処理時間
	LatestExecutionTime time.Duration
	// SlaveActivityが最後に更新された時間
	LastUpdate time.Time
	// リトライした数
	RetryCount int
}
