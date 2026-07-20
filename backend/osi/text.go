package osi

// Text は日本語・英語の両方を保持する多言語文字列。
// フロントエンドは表示時に言語を選ぶため、API は常に両方を返す。
type Text struct {
	Ja string `json:"ja"`
	En string `json:"en"`
}

func tx(ja, en string) Text { return Text{Ja: ja, En: en} }

// txSame は技術的な値（IPアドレスやフィールドダンプ等）で言語差が無い場合に使う。
func txSame(s string) Text { return Text{Ja: s, En: s} }
