package osi

// Protocol は選択可能な通信シナリオ（主に L7 アプリケーションプロトコル）を表す。
// 選択に応じて L4（TCP/UDP）や L3（IP/ICMP）の構成が変わる。
type Protocol struct {
	Key         string `json:"key"`
	L7Name      string `json:"l7Name"`
	Label       string `json:"label"`
	Transport   string `json:"transport"` // "TCP" | "UDP" | "ICMP"
	Port        int    `json:"port"`
	TLS         bool   `json:"tls"`
	L3Protocol  string `json:"l3Protocol"` // IP ヘッダの protocol フィールド表示
	Description string `json:"description"`
}

// Protocols は UI のドロップダウンに出す選択肢。
var Protocols = []Protocol{
	{
		Key: "http", L7Name: "HTTP", Label: "HTTP — Web ページ取得",
		Transport: "TCP", Port: 80, TLS: false, L3Protocol: "TCP(6)",
		Description: "信頼性の必要な Web 通信。TCP 上で動く。",
	},
	{
		Key: "https", L7Name: "HTTPS", Label: "HTTPS — 暗号化された Web",
		Transport: "TCP", Port: 443, TLS: true, L3Protocol: "TCP(6)",
		Description: "HTTP を TLS で暗号化。L6 で暗号化処理が入る。",
	},
	{
		Key: "dns", L7Name: "DNS", Label: "DNS — 名前解決",
		Transport: "UDP", Port: 53, TLS: false, L3Protocol: "UDP(17)",
		Description: "小さな問い合わせを高速に行うため UDP を使う。",
	},
	{
		Key: "rtsp", L7Name: "RTSP", Label: "RTSP — ストリーミング制御",
		Transport: "TCP", Port: 554, TLS: false, L3Protocol: "TCP(6)",
		Description: "映像配信の再生・停止などの制御に使う。制御は TCP。",
	},
	{
		Key: "rtp", L7Name: "RTP", Label: "RTP — 映像/音声のリアルタイム配信",
		Transport: "UDP", Port: 5004, TLS: false, L3Protocol: "UDP(17)",
		Description: "遅延を避けたいリアルタイム配信。多少の欠落を許容し UDP を使う。",
	},
	{
		Key: "ping", L7Name: "ICMP Echo", Label: "Ping — 疎通確認 (ICMP)",
		Transport: "ICMP", Port: 0, TLS: false, L3Protocol: "ICMP(1)",
		Description: "L3 の機能。L4(TCP/UDP)・L7 を使わず、IP の直上に ICMP が乗る。",
	},
}

// ProtocolByKey は key に対応する Protocol を返す。未知の場合は http。
func ProtocolByKey(key string) Protocol {
	for _, p := range Protocols {
		if p.Key == key {
			return p
		}
	}
	return Protocols[0]
}
