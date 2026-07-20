package osi

// Protocol は選択可能な通信シナリオを表す。
// Family が "ip" のものは L1〜L7（またはその一部）を使う IP ネットワーク通信。
// Family が "serial" のものは IP を使わず L1〜L2 のみで動くハードウェアバス通信。
type Protocol struct {
	Key           string `json:"key"`
	L7Name        string `json:"l7Name"`
	Label         string `json:"label"`
	Category      string `json:"category"`
	Family        string `json:"family"`    // "ip" | "serial"
	Transport     string `json:"transport"` // "TCP" | "UDP" | "ICMP" | "Serial"
	Port          int    `json:"port"`
	TLS           bool   `json:"tls"`
	L3Protocol    string `json:"l3Protocol"`
	SamplePayload string `json:"samplePayload"`
	Description   string `json:"description"`
}

const htmlSample = `<!DOCTYPE html>
<html>
  <head>
    <title>Sample</title>
  </head>
  <body>
    <h1>Hello</h1>
  </body>
</html>`

// Protocols は UI のドロップダウンに出す選択肢。
var Protocols = []Protocol{
	// --- Web ---
	{
		Key: "http", L7Name: "HTTP", Label: "HTTP — Web ページ取得",
		Category: "Web", Family: "ip", Transport: "TCP", Port: 80, L3Protocol: "TCP(6)",
		SamplePayload: htmlSample,
		Description:   "信頼性の必要な Web 通信。TCP 上で動き、本文は HTML など。",
	},
	{
		Key: "https", L7Name: "HTTPS", Label: "HTTPS — 暗号化された Web",
		Category: "Web", Family: "ip", Transport: "TCP", Port: 443, TLS: true, L3Protocol: "TCP(6)",
		SamplePayload: htmlSample,
		Description:   "HTTP を TLS で暗号化。L6 で暗号化処理が入る。本文は HTML。",
	},
	{
		Key: "websocket", L7Name: "WebSocket", Label: "WebSocket — 双方向通信",
		Category: "Web", Family: "ip", Transport: "TCP", Port: 80, L3Protocol: "TCP(6)",
		SamplePayload: `{"type":"message","data":"hello"}`,
		Description:   "HTTP からアップグレードして確立する双方向のリアルタイム通信。",
	},

	// --- ファイル/メール/リモート ---
	{
		Key: "ftp", L7Name: "FTP", Label: "FTP — ファイル転送",
		Category: "ファイル/メール/リモート", Family: "ip", Transport: "TCP", Port: 21, L3Protocol: "TCP(6)",
		SamplePayload: "RETR /pub/file.txt",
		Description:   "ファイル転送プロトコル。制御用に TCP を使う。",
	},
	{
		Key: "smtp", L7Name: "SMTP", Label: "SMTP — メール送信",
		Category: "ファイル/メール/リモート", Family: "ip", Transport: "TCP", Port: 25, L3Protocol: "TCP(6)",
		SamplePayload: "MAIL FROM:<a@example.com>\nRCPT TO:<b@example.com>\nSubject: Hello",
		Description:   "メール送信プロトコル。TCP 上で動く。",
	},
	{
		Key: "ssh", L7Name: "SSH", Label: "SSH — 暗号化リモート接続",
		Category: "ファイル/メール/リモート", Family: "ip", Transport: "TCP", Port: 22, L3Protocol: "TCP(6)",
		SamplePayload: "(暗号化された端末セッション)",
		Description:   "暗号化されたリモートログイン。TCP 上で動く。",
	},

	// --- IoT/メッセージング ---
	{
		Key: "mqtt", L7Name: "MQTT", Label: "MQTT — IoT メッセージング",
		Category: "IoT/メッセージング", Family: "ip", Transport: "TCP", Port: 1883, L3Protocol: "TCP(6)",
		SamplePayload: "topic: sensors/temperature\npayload: 23.5",
		Description:   "IoT 向けの軽量な Pub/Sub メッセージング。TCP 上で動く。",
	},
	{
		Key: "coap", L7Name: "CoAP", Label: "CoAP — 軽量 IoT (UDP)",
		Category: "IoT/メッセージング", Family: "ip", Transport: "UDP", Port: 5683, L3Protocol: "UDP(17)",
		SamplePayload: "GET /sensors/temp",
		Description:   "制約デバイス向けの軽量プロトコル。HTTP 風だが UDP 上で動く。",
	},

	// --- メディア ---
	{
		Key: "rtsp", L7Name: "RTSP", Label: "RTSP — ストリーミング制御",
		Category: "メディア", Family: "ip", Transport: "TCP", Port: 554, L3Protocol: "TCP(6)",
		SamplePayload: "DESCRIBE rtsp://cam/stream RTSP/1.0",
		Description:   "映像配信の再生・停止などの制御に使う。制御は TCP。",
	},
	{
		Key: "rtp", L7Name: "RTP", Label: "RTP — 映像/音声のリアルタイム配信",
		Category: "メディア", Family: "ip", Transport: "UDP", Port: 5004, L3Protocol: "UDP(17)",
		SamplePayload: "(音声/映像フレームのバイナリ)",
		Description:   "遅延を避けたいリアルタイム配信。多少の欠落を許容し UDP を使う。",
	},

	// --- インフラ ---
	{
		Key: "dns", L7Name: "DNS", Label: "DNS — 名前解決",
		Category: "インフラ", Family: "ip", Transport: "UDP", Port: 53, L3Protocol: "UDP(17)",
		SamplePayload: "QNAME=example.com QTYPE=A",
		Description:   "小さな問い合わせを高速に行うため UDP を使う。",
	},
	{
		Key: "dhcp", L7Name: "DHCP", Label: "DHCP — IP アドレス自動割当",
		Category: "インフラ", Family: "ip", Transport: "UDP", Port: 67, L3Protocol: "UDP(17)",
		SamplePayload: "DHCPDISCOVER",
		Description:   "IP アドレスを自動で配布する。UDP のブロードキャストを使う。",
	},
	{
		Key: "ntp", L7Name: "NTP", Label: "NTP — 時刻同期",
		Category: "インフラ", Family: "ip", Transport: "UDP", Port: 123, L3Protocol: "UDP(17)",
		SamplePayload: "(時刻同期リクエスト)",
		Description:   "ネットワーク越しに時刻を同期する。UDP を使う。",
	},
	{
		Key: "snmp", L7Name: "SNMP", Label: "SNMP — 機器監視",
		Category: "インフラ", Family: "ip", Transport: "UDP", Port: 161, L3Protocol: "UDP(17)",
		SamplePayload: "GET sysUpTime.0",
		Description:   "ネットワーク機器の監視・管理に使う。UDP を使う。",
	},

	// --- 診断 ---
	{
		Key: "ping", L7Name: "ICMP Echo", Label: "Ping — 疎通確認 (ICMP)",
		Category: "診断", Family: "ip", Transport: "ICMP", Port: 0, L3Protocol: "ICMP(1)",
		SamplePayload: "abcdefghijklmnop",
		Description:   "L3 の機能。L4(TCP/UDP)・L7 を使わず、IP の直上に ICMP が乗る。",
	},

	// --- シリアル通信（L1-L2 のみ・IP を使わない） ---
	{
		Key: "uart", L7Name: "UART", Label: "UART — 非同期シリアル通信",
		Category: "シリアル通信 (L1-L2)", Family: "serial", Transport: "Serial",
		SamplePayload: "Hi",
		Description:   "IP を使わない。1 バイトを Start/Data/Parity/Stop ビットで囲んで TX/RX 線で送る。",
	},
	{
		Key: "i2c", L7Name: "I2C", Label: "I2C — 2 線式バス通信",
		Category: "シリアル通信 (L1-L2)", Family: "serial", Transport: "Serial",
		SamplePayload: "reg=0x00 data=0x1F",
		Description:   "IP を使わない。SDA/SCL の 2 線でデバイスアドレス指定して通信するバス。",
	},
	{
		Key: "spi", L7Name: "SPI", Label: "SPI — 高速同期シリアル",
		Category: "シリアル通信 (L1-L2)", Family: "serial", Transport: "Serial",
		SamplePayload: "0x9F (Read ID)",
		Description:   "IP を使わない。MOSI/MISO/SCLK/CS の複数線でチップを選択して全二重通信する。",
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
