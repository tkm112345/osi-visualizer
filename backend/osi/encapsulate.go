package osi

import (
	"fmt"
	"strconv"
	"strings"
)

// Request はカプセル化シミュレートへの入力。
type Request struct {
	Message  string `json:"message"`
	SrcIP    string `json:"srcIp"`
	DstIP    string `json:"dstIp"`
	Protocol string `json:"protocol"` // Protocols の key（空なら http）
}

// FramePart は、ある層でのPDUを構成する 1 区画（ヘッダ / ペイロード / トレーラ）。
// アコーディオンで「実際どんなデータになっているか」を見せるために使う。
type FramePart struct {
	Label  string `json:"label"`  // 例: "IP ヘッダ", "ペイロード (HTML)"
	Detail string `json:"detail"` // 実際のフィールド値やペイロード内容
	Kind   string `json:"kind"`   // "header" | "payload" | "trailer"
	Bytes  int    `json:"bytes"`
}

// Step は 1 レイヤーでの処理結果を表す。L7 → L1 の順に積み上がる。
type Step struct {
	Level       int               `json:"level"`
	Name        string            `json:"name"`
	NameJa      string            `json:"nameJa"`
	PDU         string            `json:"pdu"`
	AddsHeader  bool              `json:"addsHeader"`
	Active      bool              `json:"active"` // このシナリオでこの層が使われるか
	Headers     map[string]string `json:"headers"`
	Processing  []string          `json:"processing"`
	Payload     string            `json:"payload"`
	HeaderBytes int               `json:"headerBytes"`
	TotalBytes  int               `json:"totalBytes"`
	Structure   string            `json:"structure"`
	Note        string            `json:"note"`
	Bitstream   string            `json:"bitstream"`
	Frame       []FramePart       `json:"frame"` // この層でのPDU構造（実データ）
}

func payloadLabel(p Protocol) string {
	switch p.Key {
	case "http", "https", "websocket":
		return "ペイロード (本文)"
	case "ping":
		return "ICMP データ"
	default:
		return "ペイロード (" + p.L7Name + ")"
	}
}

// tlsEncrypted は TLS 暗号化後の見かけ（教育用のダミー表現）を返す。
func tlsEncrypted(msg string) string {
	return fmt.Sprintf("🔒 Application Data (%dB, TLS暗号化されており中身は読めない)", len(msg))
}

func tcpDetail(p Protocol) string {
	return fmt.Sprintf("srcPort=49152  dstPort=%d  seq=1  ack=1  flags=PSH,ACK  win=64240", p.Port)
}
func udpDetail(p Protocol, dataLen int) string {
	return fmt.Sprintf("srcPort=49152  dstPort=%d  length=%d  checksum=0x1a2b", p.Port, udpHeaderBytes+dataLen)
}
func ipDetail(req Request, p Protocol) string {
	return fmt.Sprintf("version=4  ihl=5  ttl=64  proto=%s  src=%s  dst=%s", p.L3Protocol, req.SrcIP, req.DstIP)
}
func icmpDetail() string { return "type=8 (Echo Request)  code=0  id=0x0001  seq=1" }
func ethDetail() string {
	return "dst=11:22:33:44:55:66  src=AA:BB:CC:DD:EE:01  ethertype=0x0800"
}

func hdr(label, detail string, bytes int) FramePart {
	return FramePart{Label: label, Detail: detail, Kind: "header", Bytes: bytes}
}

// 各レイヤーが付与するヘッダのバイト数（教育用の代表値）。
const (
	tcpHeaderBytes  = 20
	udpHeaderBytes  = 8
	icmpHeaderBytes = 8
	ipHeaderBytes   = 20
	ethHeaderBytes  = 14
	ethTrailer      = 4 // FCS
)

func normalize(req Request) Request {
	if req.SrcIP == "" {
		req.SrcIP = "192.168.0.10"
	}
	if req.DstIP == "" {
		req.DstIP = "93.184.216.34"
	}
	return req
}

// l7Headers は選択プロトコルに応じた L7 ヘッダ表示を返す。
func l7Headers(p Protocol, dstIP string) map[string]string {
	port := strconv.Itoa(p.Port)
	switch p.Key {
	case "http", "https":
		return map[string]string{"protocol": p.L7Name, "method": "GET", "host": dstIP, "dstPort": port}
	case "dns":
		return map[string]string{"protocol": "DNS", "query": "example.com", "recordType": "A", "dstPort": port}
	case "rtsp":
		return map[string]string{"protocol": "RTSP", "method": "DESCRIBE", "dstPort": port}
	case "rtp":
		return map[string]string{"protocol": "RTP", "payloadType": "96", "dstPort": port}
	default:
		return map[string]string{"protocol": p.L7Name, "dstPort": port}
	}
}

// transportHeaders は TCP/UDP のヘッダ表示とバイト数を返す。
func transportHeaders(p Protocol) (map[string]string, int) {
	port := strconv.Itoa(p.Port)
	if p.Transport == "UDP" {
		return map[string]string{"srcPort": "49152", "dstPort": port, "length": "8+data", "checksum": "0x1a2b"}, udpHeaderBytes
	}
	return map[string]string{"srcPort": "49152", "dstPort": port, "seq": "1", "flags": "PSH,ACK"}, tcpHeaderBytes
}

// serialFramingBytes は L2 でのフレーミング相当のバイト数（教育用の代表値）。
func serialFramingBytes(p Protocol) int {
	if p.Key == "i2c" {
		return 1 // 7bit アドレス + R/W ≒ 1 バイト
	}
	return 0 // UART/SPI はビット単位の枠付けで、明確なヘッダバイトは持たない
}

// serialL2Headers / serialL1Headers はシリアル通信の L2/L1 表示を返す。
func serialL2Headers(p Protocol) (map[string]string, string) {
	switch p.Key {
	case "uart":
		return map[string]string{
			"frame":  "Start(1) + Data(8) + Parity + Stop(1)",
			"parity": "None",
			"flow":   "None",
		}, "1 バイトごとに Start/Stop ビットで枠付け（フレーミング）する。"
	case "i2c":
		return map[string]string{
			"address": "0x3C (7bit)",
			"rw":      "Write(0)",
			"ack":     "ACK/NACK",
		}, "先頭でスレーブアドレスと R/W を送り、各バイトで ACK を確認する。"
	case "spi":
		return map[string]string{
			"chipSelect": "CS0 (Low)",
			"mode":       "Mode 0 (CPOL=0, CPHA=0)",
		}, "CS 線で通信相手のチップを選択する。アドレスの概念は無い。"
	default:
		return map[string]string{}, ""
	}
}

func serialL1Info(p Protocol) ([]string, string) {
	switch p.Key {
	case "uart":
		return []string{"信号線: TX / RX", "ボーレート: 9600 bps", "電圧レベル: 3.3V"}, "TX/RX の 2 線で、クロックを共有せず非同期にビットを送る。"
	case "i2c":
		return []string{"信号線: SDA / SCL", "配線: オープンドレイン + プルアップ", "クロック: 400 kHz"}, "SDA(データ)/SCL(クロック)の 2 線で通信する。"
	case "spi":
		return []string{"信号線: MOSI / MISO / SCLK / CS", "クロック: SCLK を共有"}, "複数線でクロックを共有し全二重で送受信する。"
	default:
		return nil, ""
	}
}

// serialL2Detail は L2 フレーミングの実データ表現を返す。
func serialL2Detail(p Protocol) string {
	switch p.Key {
	case "uart":
		return "各バイトを Start(1) + Data(8) + Stop(1) ビットで枠付け, Parity=None"
	case "i2c":
		return "Start + Address(0x76,7bit) + R/W(0) + ACK ... 各バイト後に ACK"
	case "spi":
		return "CS=Low で選択, SCLK に同期して MOSI/MISO を全二重送受信"
	default:
		return ""
	}
}

// encapsulateSerial は UART/I2C/SPI など IP を使わない L1/L2 のみの通信を組み立てる。
func encapsulateSerial(req Request, p Protocol) []Step {
	framing := serialFramingBytes(p)
	total := len(req.Message)
	structure := "[Data]"
	steps := make([]Step, 0, len(Layers))

	pl := FramePart{Label: payloadLabel(p), Detail: req.Message, Kind: "payload", Bytes: len(req.Message)}
	var hdrs []FramePart

	for _, l := range Layers {
		step := Step{
			Level: l.Level, Name: l.Name, NameJa: l.NameJa, PDU: l.PDU,
			AddsHeader: l.AddsHeader, Active: true,
			Payload: req.Message, Headers: map[string]string{},
		}
		switch l.Level {
		case 7, 6, 5, 4, 3:
			step.Active = false
			step.Note = p.L7Name + " は IP ネットワークを使わない。L3〜L7 は無く、L1/L2 だけで通信する。"
		case 2:
			headers, note := serialL2Headers(p)
			step.Headers = headers
			step.HeaderBytes = framing
			total += framing
			structure = "[" + p.L7Name + " " + structure + "]"
			step.Note = note
			hdrs = append([]FramePart{hdr(p.L7Name+" フレーミング", serialL2Detail(p), framing)}, hdrs...)
		case 1:
			proc, note := serialL1Info(p)
			step.Processing = proc
			step.Note = note
			step.Bitstream = toBits(req.Message, 4)
		}
		step.TotalBytes = total
		step.Structure = structure
		if step.Active {
			out := append([]FramePart{}, hdrs...)
			step.Frame = append(out, pl)
		}
		steps = append(steps, step)
	}
	return steps
}

// Encapsulate は入力メッセージを L7 → L1 へカプセル化した各ステップを返す。
func Encapsulate(req Request) []Step {
	req = normalize(req)
	p := ProtocolByKey(req.Protocol)
	if p.Family == "serial" {
		return encapsulateSerial(req, p)
	}
	isPing := p.Transport == "ICMP"

	total := len(req.Message)
	structure := "[Data]"
	steps := make([]Step, 0, len(Layers))

	// frame は「その層での実データ構造」を組み立てるための状態。
	pl := FramePart{Label: payloadLabel(p), Detail: req.Message, Kind: "payload", Bytes: len(req.Message)}
	var hdrs []FramePart
	var trailer *FramePart
	buildFrame := func() []FramePart {
		out := append([]FramePart{}, hdrs...)
		out = append(out, pl)
		if trailer != nil {
			out = append(out, *trailer)
		}
		return out
	}

	for _, l := range Layers {
		step := Step{
			Level:      l.Level,
			Name:       l.Name,
			NameJa:     l.NameJa,
			PDU:        l.PDU,
			AddsHeader: l.AddsHeader,
			Active:     true,
			Payload:    req.Message,
			Headers:    map[string]string{},
		}

		switch l.Level {
		case 7:
			if isPing {
				step.Active = false
				step.Note = "Ping(ICMP) はアプリケーション層プロトコルを使わない。"
			} else {
				step.Headers = l7Headers(p, req.DstIP)
				step.Note = "アプリ（" + p.L7Name + "）が生成したデータ本体。これが最内のペイロードになる。"
			}

		case 6:
			if isPing {
				step.Active = false
				step.Note = "Ping では使用しない。"
			} else {
				step.Processing = []string{"文字コード変換 (UTF-8)"}
				if p.TLS {
					step.Processing = append([]string{"TLS による暗号化"}, step.Processing...)
					step.Note = p.L7Name + " なので、この層で TLS 暗号化が行われる。以降ペイロードは暗号文になる。"
					pl.Detail = tlsEncrypted(req.Message)
					pl.Label = "ペイロード (TLS暗号化)"
				} else {
					step.Note = "独立したヘッダは付与しない。TCP/IP モデルでは Application 層に含まれる。"
				}
			}

		case 5:
			if isPing {
				step.Active = false
				step.Note = "Ping では使用しない。"
			} else {
				step.Processing = []string{"セッション確立/維持", "同期ポイントの管理"}
				step.Note = "独立したヘッダは付与しない。TCP/IP モデルでは Application 層に含まれる。"
			}

		case 4:
			if isPing {
				step.Active = false
				step.Note = "ICMP は L4(TCP/UDP) を使わず、IP の直上で動作する。"
			} else {
				headers, hb := transportHeaders(p)
				step.Headers = headers
				step.HeaderBytes = hb
				total += hb
				structure = "[" + p.Transport + " " + structure + "]"
				step.Note = fmt.Sprintf("ポート %d でアプリを識別。%s ヘッダ(%dB)を付与する。", p.Port, p.Transport, hb)
				detail := tcpDetail(p)
				if p.Transport == "UDP" {
					detail = udpDetail(p, len(req.Message))
				}
				hdrs = append([]FramePart{hdr(p.Transport+" ヘッダ", detail, hb)}, hdrs...)
			}

		case 3:
			step.Headers = map[string]string{
				"srcIp":    req.SrcIP,
				"dstIp":    req.DstIP,
				"ttl":      "64",
				"protocol": p.L3Protocol,
			}
			if isPing {
				// ICMP ヘッダは IP の直上（L3 の中）に置かれる。
				total += icmpHeaderBytes
				structure = "[ICMP " + structure + "]"
				step.Headers["icmpType"] = "8 (Echo Request)"
				step.Headers["icmpCode"] = "0"
				step.HeaderBytes = icmpHeaderBytes + ipHeaderBytes
				step.Note = "ICMP Echo Request を作り、IP ヘッダを付与する。TCP/UDP は挟まらない。"
				hdrs = append([]FramePart{hdr("ICMP ヘッダ", icmpDetail(), icmpHeaderBytes)}, hdrs...)
			} else {
				step.HeaderBytes = ipHeaderBytes
				step.Note = "IP アドレスを付与し、ネットワーク間の経路制御を可能にする。"
			}
			total += ipHeaderBytes
			structure = "[IP " + structure + "]"
			hdrs = append([]FramePart{hdr("IP ヘッダ", ipDetail(req, p), ipHeaderBytes)}, hdrs...)

		case 2:
			step.Headers = map[string]string{
				"srcMac":    "AA:BB:CC:DD:EE:01",
				"dstMac":    "11:22:33:44:55:66",
				"etherType": "0x0800",
				"fcs":       "(4B trailer)",
			}
			step.HeaderBytes = ethHeaderBytes + ethTrailer
			total += ethHeaderBytes + ethTrailer
			structure = "[Eth " + structure + " FCS]"
			step.Note = "MAC アドレスを付与してフレーム化。末尾に FCS（誤り検出）も付く。"
			hdrs = append([]FramePart{hdr("Ethernet ヘッダ", ethDetail(), ethHeaderBytes)}, hdrs...)
			trailer = &FramePart{Label: "FCS (CRC32)", Detail: "0x1A2B3C4D", Kind: "trailer", Bytes: ethTrailer}

		case 1:
			step.Processing = []string{"ビット列を電気/光/電波の信号に変換"}
			step.Note = "フレームをビット列として物理媒体に送出する。"
			step.Bitstream = toBits(req.Message, 4)
		}

		step.TotalBytes = total
		step.Structure = structure
		if step.Active {
			step.Frame = buildFrame()
		}
		steps = append(steps, step)
	}

	return steps
}

// toBits は文字列の先頭 n バイトを 8 桁ビット表現に変換する（L1 表示用）。
func toBits(s string, n int) string {
	b := []byte(s)
	if len(b) > n {
		b = b[:n]
	}
	parts := make([]string, 0, len(b))
	for _, c := range b {
		parts = append(parts, fmt.Sprintf("%08b", c))
	}
	suffix := ""
	if len(s) > n {
		suffix = " ..."
	}
	return strings.Join(parts, " ") + suffix
}
