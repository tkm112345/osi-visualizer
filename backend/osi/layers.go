package osi

// Layer は OSI 参照モデルの 1 レイヤーの静的なメタ情報を表す。
type Layer struct {
	Level       int      `json:"level"`
	Name        string   `json:"name"`
	NameJa      string   `json:"nameJa"`
	PDU         string   `json:"pdu"`
	Protocols   []string `json:"protocols"`
	Description string   `json:"description"`
	// AddsHeader が false のレイヤーは、実際のパケットに独立したヘッダを付与しない
	// （L5/L6 は TCP/IP モデルでは Application 層に統合されるため）。
	AddsHeader bool `json:"addsHeader"`
}

// Layers は L7 → L1 の順（上位層が先頭）で定義した OSI 全 7 層。
var Layers = []Layer{
	{
		Level: 7, Name: "Application", NameJa: "アプリケーション層", PDU: "Data",
		Protocols:   []string{"HTTP", "DNS", "FTP", "SMTP"},
		Description: "アプリケーションがネットワークを利用するためのインターフェース。ユーザーが直接触れるサービスを提供する。",
		AddsHeader:  true,
	},
	{
		Level: 6, Name: "Presentation", NameJa: "プレゼンテーション層", PDU: "Data",
		Protocols:   []string{"TLS/SSL", "JPEG", "ASCII", "UTF-8"},
		Description: "データの表現形式を変換する。暗号化・圧縮・文字コード変換など。TCP/IP モデルでは Application 層に統合される。",
		AddsHeader:  false,
	},
	{
		Level: 5, Name: "Session", NameJa: "セッション層", PDU: "Data",
		Protocols:   []string{"RPC", "NetBIOS", "SMB"},
		Description: "通信の開始から終了までの対話（セッション）を管理する。TCP/IP モデルでは Application 層に統合される。",
		AddsHeader:  false,
	},
	{
		Level: 4, Name: "Transport", NameJa: "トランスポート層", PDU: "Segment",
		Protocols:   []string{"TCP", "UDP"},
		Description: "エンドツーエンドの通信を担う。ポート番号でアプリを識別し、TCP なら順序制御・再送で信頼性を確保する。",
		AddsHeader:  true,
	},
	{
		Level: 3, Name: "Network", NameJa: "ネットワーク層", PDU: "Packet",
		Protocols:   []string{"IP", "ICMP", "ARP"},
		Description: "IP アドレスを使って、異なるネットワーク間でエンドツーエンドの経路制御（ルーティング）を行う。",
		AddsHeader:  true,
	},
	{
		Level: 2, Name: "Data Link", NameJa: "データリンク層", PDU: "Frame",
		Protocols:   []string{"Ethernet", "PPP", "MAC"},
		Description: "同一ネットワーク内の隣接ノード間で、MAC アドレスを使ってフレームを転送する。誤り検出も行う。",
		AddsHeader:  true,
	},
	{
		Level: 1, Name: "Physical", NameJa: "物理層", PDU: "Bits",
		Protocols:   []string{"Ethernet(PHY)", "Wi-Fi", "光ファイバ"},
		Description: "ビット列を電気信号・光信号・電波などの物理的な信号に変換して伝送する。",
		AddsHeader:  false,
	},
}
