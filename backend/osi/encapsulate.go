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

// Encapsulate は入力メッセージを L7 → L1 へカプセル化した各ステップを返す。
func Encapsulate(req Request) []Step {
	req = normalize(req)
	p := ProtocolByKey(req.Protocol)
	isPing := p.Transport == "ICMP"

	total := len(req.Message)
	structure := "[Data]"
	steps := make([]Step, 0, len(Layers))

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
					step.Note = p.L7Name + " なので、この層で TLS 暗号化が行われる。"
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
			} else {
				step.HeaderBytes = ipHeaderBytes
				step.Note = "IP アドレスを付与し、ネットワーク間の経路制御を可能にする。"
			}
			total += ipHeaderBytes
			structure = "[IP " + structure + "]"

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

		case 1:
			step.Processing = []string{"ビット列を電気/光/電波の信号に変換"}
			step.Note = "フレームをビット列として物理媒体に送出する。"
			step.Bitstream = toBits(req.Message, 4)
		}

		step.TotalBytes = total
		step.Structure = structure
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
