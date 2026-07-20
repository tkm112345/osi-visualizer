package osi

// Protocol は選択可能な通信シナリオを表す。
// Family が "ip" のものは L1〜L7（またはその一部）を使う IP ネットワーク通信。
// Family が "serial" のものは IP を使わず L1〜L2 のみで動くハードウェアバス通信。
type Protocol struct {
	Key           string `json:"key"`
	L7Name        string `json:"l7Name"`
	Label         Text   `json:"label"`
	Category      Text   `json:"category"`
	Family        string `json:"family"`    // "ip" | "serial"
	Transport     string `json:"transport"` // "TCP" | "UDP" | "ICMP" | "Serial"
	Port          int    `json:"port"`
	TLS           bool   `json:"tls"`
	L3Protocol    string `json:"l3Protocol"`
	SamplePayload string `json:"samplePayload"`
	Description   Text   `json:"description"`
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

// カテゴリ名（多言語）。ドロップダウンのグループ見出しに使う。
var (
	catWeb    = tx("Web", "Web")
	catFile   = tx("ファイル/メール/リモート", "File / Mail / Remote")
	catIoT    = tx("IoT/メッセージング", "IoT / Messaging")
	catMedia  = tx("メディア", "Media")
	catInfra  = tx("インフラ", "Infrastructure")
	catDiag   = tx("診断", "Diagnostics")
	catSerial = tx("シリアル通信 (L1-L2)", "Serial (L1-L2)")
)

// Protocols は UI のドロップダウンに出す選択肢。
var Protocols = []Protocol{
	// --- Web ---
	{
		Key: "http", L7Name: "HTTP", Label: tx("HTTP — Web ページ取得", "HTTP — Fetch a web page"),
		Category: catWeb, Family: "ip", Transport: "TCP", Port: 80, L3Protocol: "TCP(6)",
		SamplePayload: htmlSample,
		Description: tx("信頼性の必要な Web 通信。TCP 上で動き、本文は HTML など。",
			"Reliable web traffic over TCP; the body is HTML, etc."),
	},
	{
		Key: "https", L7Name: "HTTPS", Label: tx("HTTPS — 暗号化された Web", "HTTPS — Encrypted web"),
		Category: catWeb, Family: "ip", Transport: "TCP", Port: 443, TLS: true, L3Protocol: "TCP(6)",
		SamplePayload: htmlSample,
		Description: tx("HTTP を TLS で暗号化。L6 で暗号化処理が入る。本文は HTML。",
			"HTTP encrypted with TLS; encryption happens at L6. The body is HTML."),
	},
	{
		Key: "websocket", L7Name: "WebSocket", Label: tx("WebSocket — 双方向通信", "WebSocket — Bidirectional"),
		Category: catWeb, Family: "ip", Transport: "TCP", Port: 80, L3Protocol: "TCP(6)",
		SamplePayload: `{"type":"message","data":"hello"}`,
		Description: tx("HTTP からアップグレードして確立する双方向のリアルタイム通信。",
			"Bidirectional real-time communication upgraded from HTTP."),
	},

	// --- ファイル/メール/リモート ---
	{
		Key: "ftp", L7Name: "FTP", Label: tx("FTP — ファイル転送", "FTP — File transfer"),
		Category: catFile, Family: "ip", Transport: "TCP", Port: 21, L3Protocol: "TCP(6)",
		SamplePayload: "RETR /pub/file.txt",
		Description:   tx("ファイル転送プロトコル。制御用に TCP を使う。", "File transfer protocol; uses TCP for control."),
	},
	{
		Key: "smtp", L7Name: "SMTP", Label: tx("SMTP — メール送信", "SMTP — Send mail"),
		Category: catFile, Family: "ip", Transport: "TCP", Port: 25, L3Protocol: "TCP(6)",
		SamplePayload: "MAIL FROM:<a@example.com>\nRCPT TO:<b@example.com>\nSubject: Hello",
		Description:   tx("メール送信プロトコル。TCP 上で動く。", "Mail-sending protocol running over TCP."),
	},
	{
		Key: "ssh", L7Name: "SSH", Label: tx("SSH — 暗号化リモート接続", "SSH — Encrypted remote access"),
		Category: catFile, Family: "ip", Transport: "TCP", Port: 22, L3Protocol: "TCP(6)",
		SamplePayload: "(暗号化された端末セッション)",
		Description:   tx("暗号化されたリモートログイン。TCP 上で動く。", "Encrypted remote login running over TCP."),
	},

	// --- IoT/メッセージング ---
	{
		Key: "mqtt", L7Name: "MQTT", Label: tx("MQTT — IoT メッセージング", "MQTT — IoT messaging"),
		Category: catIoT, Family: "ip", Transport: "TCP", Port: 1883, L3Protocol: "TCP(6)",
		SamplePayload: "topic: sensors/temperature\npayload: 23.5",
		Description: tx("IoT 向けの軽量な Pub/Sub メッセージング。TCP 上で動く。",
			"Lightweight pub/sub messaging for IoT over TCP."),
	},
	{
		Key: "coap", L7Name: "CoAP", Label: tx("CoAP — 軽量 IoT (UDP)", "CoAP — Lightweight IoT (UDP)"),
		Category: catIoT, Family: "ip", Transport: "UDP", Port: 5683, L3Protocol: "UDP(17)",
		SamplePayload: "GET /sensors/temp",
		Description: tx("制約デバイス向けの軽量プロトコル。HTTP 風だが UDP 上で動く。",
			"Lightweight protocol for constrained devices; HTTP-like but over UDP."),
	},

	// --- メディア ---
	{
		Key: "rtsp", L7Name: "RTSP", Label: tx("RTSP — ストリーミング制御", "RTSP — Streaming control"),
		Category: catMedia, Family: "ip", Transport: "TCP", Port: 554, L3Protocol: "TCP(6)",
		SamplePayload: "DESCRIBE rtsp://cam/stream RTSP/1.0",
		Description: tx("映像配信の再生・停止などの制御に使う。制御は TCP。",
			"Controls playback (play/pause) of video streams; control over TCP."),
	},
	{
		Key: "rtp", L7Name: "RTP", Label: tx("RTP — 映像/音声のリアルタイム配信", "RTP — Real-time audio/video"),
		Category: catMedia, Family: "ip", Transport: "UDP", Port: 5004, L3Protocol: "UDP(17)",
		SamplePayload: "(音声/映像フレームのバイナリ)",
		Description: tx("遅延を避けたいリアルタイム配信。多少の欠落を許容し UDP を使う。",
			"Real-time delivery that avoids delay; tolerates some loss and uses UDP."),
	},

	// --- インフラ ---
	{
		Key: "dns", L7Name: "DNS", Label: tx("DNS — 名前解決", "DNS — Name resolution"),
		Category: catInfra, Family: "ip", Transport: "UDP", Port: 53, L3Protocol: "UDP(17)",
		SamplePayload: "QNAME=example.com QTYPE=A",
		Description:   tx("小さな問い合わせを高速に行うため UDP を使う。", "Uses UDP for fast, small queries."),
	},
	{
		Key: "dhcp", L7Name: "DHCP", Label: tx("DHCP — IP アドレス自動割当", "DHCP — Automatic IP assignment"),
		Category: catInfra, Family: "ip", Transport: "UDP", Port: 67, L3Protocol: "UDP(17)",
		SamplePayload: "DHCPDISCOVER",
		Description: tx("IP アドレスを自動で配布する。UDP のブロードキャストを使う。",
			"Automatically hands out IP addresses using UDP broadcast."),
	},
	{
		Key: "ntp", L7Name: "NTP", Label: tx("NTP — 時刻同期", "NTP — Time synchronization"),
		Category: catInfra, Family: "ip", Transport: "UDP", Port: 123, L3Protocol: "UDP(17)",
		SamplePayload: "(時刻同期リクエスト)",
		Description:   tx("ネットワーク越しに時刻を同期する。UDP を使う。", "Synchronizes time across the network using UDP."),
	},
	{
		Key: "snmp", L7Name: "SNMP", Label: tx("SNMP — 機器監視", "SNMP — Device monitoring"),
		Category: catInfra, Family: "ip", Transport: "UDP", Port: 161, L3Protocol: "UDP(17)",
		SamplePayload: "GET sysUpTime.0",
		Description:   tx("ネットワーク機器の監視・管理に使う。UDP を使う。", "Monitors and manages network devices using UDP."),
	},

	// --- 診断 ---
	{
		Key: "ping", L7Name: "ICMP Echo", Label: tx("Ping — 疎通確認 (ICMP)", "Ping — Reachability check (ICMP)"),
		Category: catDiag, Family: "ip", Transport: "ICMP", Port: 0, L3Protocol: "ICMP(1)",
		SamplePayload: "abcdefghijklmnopqrstuvwabcdefghi",
		Description: tx(
			"L3 の機能。L4(TCP/UDP)・L7 を使わず、IP の直上に ICMP が乗る。ペイロードは往復測定用の 32B パディング（Windows の ping と同じ内容）。",
			"An L3 feature. It uses no L4 (TCP/UDP) or L7; ICMP rides directly on IP. The payload is 32B of round-trip padding (same bytes Windows ping sends)."),
	},

	// --- シリアル通信（L1-L2 のみ・IP を使わない） ---
	{
		Key: "uart", L7Name: "UART", Label: tx("UART — 非同期シリアル通信", "UART — Asynchronous serial"),
		Category: catSerial, Family: "serial", Transport: "Serial",
		SamplePayload: "$GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M*47",
		Description: tx(
			"IP を使わない。1 バイトを Start/Data/Parity/Stop ビットで囲んで TX/RX 線で送る。例は GPS の NMEA センテンス。",
			"No IP. Each byte is framed with start/data/parity/stop bits over TX/RX lines. The sample is a GPS NMEA sentence."),
	},
	{
		Key: "i2c", L7Name: "I2C", Label: tx("I2C — 2 線式バス通信", "I2C — Two-wire bus"),
		Category: catSerial, Family: "serial", Transport: "Serial",
		SamplePayload: "ADDR=0x76 REG=0xF7 READ 6B → 温度/気圧レジスタ",
		Description: tx(
			"IP を使わない。SDA/SCL の 2 線でデバイスアドレス指定して通信するバス。例は BME280 気圧センサの測定値読み出し。",
			"No IP. A two-wire bus (SDA/SCL) that addresses devices. The sample reads measurements from a BME280 pressure sensor."),
	},
	{
		Key: "spi", L7Name: "SPI", Label: tx("SPI — 高速同期シリアル", "SPI — High-speed synchronous serial"),
		Category: catSerial, Family: "serial", Transport: "Serial",
		SamplePayload: "MOSI: 9F 00 00 00  MISO: -- EF 40 18 (JEDEC ID)",
		Description: tx(
			"IP を使わない。MOSI/MISO/SCLK/CS の複数線でチップを選択して全二重通信する。例は SPI フラッシュの JEDEC ID 読み出し。",
			"No IP. Multiple lines (MOSI/MISO/SCLK/CS) select a chip and transfer full-duplex. The sample reads a SPI flash JEDEC ID."),
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
