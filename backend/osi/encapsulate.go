package osi

import (
	"fmt"
	"strings"
)

// Request はカプセル化シミュレートへの入力。
type Request struct {
	Message string `json:"message"`
	SrcIP   string `json:"srcIp"`
	DstIP   string `json:"dstIp"`
}

// Step は 1 レイヤーでの処理結果を表す。L7 → L1 の順に積み上がる。
type Step struct {
	Level      int               `json:"level"`
	Name       string            `json:"name"`
	NameJa     string            `json:"nameJa"`
	PDU        string            `json:"pdu"`
	AddsHeader bool              `json:"addsHeader"`
	Headers    map[string]string `json:"headers"`    // このレイヤーで付与したヘッダ
	Processing []string          `json:"processing"` // ヘッダを付けないレイヤーが行う処理
	Payload    string            `json:"payload"`    // このレイヤーが受け取ったデータ（人が読める要約）
	HeaderBytes int              `json:"headerBytes"`
	TotalBytes  int              `json:"totalBytes"` // このレイヤーまでの累積バイト数
	Structure   string            `json:"structure"`  // 例: [Eth [IP [TCP [Data]]]]
	Note        string            `json:"note"`
	Bitstream   string            `json:"bitstream"` // L1 のみ: 先頭数バイトのビット表現
}

// 各レイヤーが付与するヘッダのバイト数（教育用の代表値）。
const (
	tcpHeaderBytes = 20
	ipHeaderBytes  = 20
	ethHeaderBytes = 14
	ethTrailer     = 4 // FCS
)

// Encapsulate は入力メッセージを L7 → L1 へカプセル化した各ステップを返す。
func Encapsulate(req Request) []Step {
	if req.SrcIP == "" {
		req.SrcIP = "192.168.0.10"
	}
	if req.DstIP == "" {
		req.DstIP = "93.184.216.34"
	}

	payloadBytes := len(req.Message)
	total := payloadBytes
	structure := "[Data]"

	steps := make([]Step, 0, len(Layers))

	for _, l := range Layers {
		step := Step{
			Level:      l.Level,
			Name:       l.Name,
			NameJa:     l.NameJa,
			PDU:        l.PDU,
			AddsHeader: l.AddsHeader,
			Payload:    req.Message,
			Headers:    map[string]string{},
		}

		switch l.Level {
		case 7:
			step.Headers = map[string]string{
				"protocol": "HTTP",
				"method":   "GET",
				"host":     req.DstIP,
			}
			step.Note = "アプリが生成したデータ本体。ここが最内のペイロードになる。"

		case 6:
			step.Processing = []string{"文字コード変換 (UTF-8)", "TLS による暗号化"}
			step.Note = "独立したヘッダは付与しない。TCP/IP モデルでは Application 層に含まれる。"

		case 5:
			step.Processing = []string{"セッション確立/維持", "同期ポイントの管理"}
			step.Note = "独立したヘッダは付与しない。TCP/IP モデルでは Application 層に含まれる。"

		case 4:
			step.Headers = map[string]string{
				"srcPort": "49152",
				"dstPort": "80",
				"seq":     "1",
				"flags":   "PSH,ACK",
			}
			step.HeaderBytes = tcpHeaderBytes
			total += tcpHeaderBytes
			structure = "[TCP " + structure + "]"
			step.Note = "ポート番号でアプリを識別。TCP ヘッダを付与してセグメント化する。"

		case 3:
			step.Headers = map[string]string{
				"srcIp":    req.SrcIP,
				"dstIp":    req.DstIP,
				"ttl":      "64",
				"protocol": "TCP(6)",
			}
			step.HeaderBytes = ipHeaderBytes
			total += ipHeaderBytes
			structure = "[IP " + structure + "]"
			step.Note = "IP アドレスを付与し、ネットワーク間の経路制御を可能にする。"

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
