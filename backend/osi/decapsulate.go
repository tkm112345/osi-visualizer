package osi

import (
	"fmt"
	"strconv"
)

// Decapsulate は受信ホスト視点で、届いたフレームを L1 → L7 へ
// デカプセル化（ヘッダを解析して外していく）各ステップを返す。
// 選択プロトコルにより L4(TCP/UDP) の有無や ICMP かどうかが変わる。あくまで擬似。
// decapsulateSerial は受信側のシリアル通信（L1 → L2 のみ）を組み立てる。
func decapsulateSerial(req Request, p Protocol) []Step {
	framing := serialFramingBytes(p)
	remaining := len(req.Message) + framing
	structure := "[" + p.L7Name + " [Data]]"
	steps := make([]Step, 0, len(Layers))

	for i := len(Layers) - 1; i >= 0; i-- {
		l := Layers[i]
		step := Step{
			Level: l.Level, Name: l.Name, NameJa: l.NameJa, PDU: l.PDU,
			AddsHeader: l.AddsHeader, Active: true,
			Payload: req.Message, Headers: map[string]string{},
		}
		switch l.Level {
		case 1:
			proc, note := serialL1Info(p)
			step.Processing = proc
			step.Note = note
			step.Bitstream = toBits(req.Message, 4)
		case 2:
			headers, _ := serialL2Headers(p)
			step.Headers = headers
			step.HeaderBytes = framing
			remaining -= framing
			structure = "[Data]"
			if p.Key == "i2c" {
				step.Note = "自分のアドレス宛か確認し、ACK を返してデータを受け取る。"
			} else {
				step.Note = "Start/Stop ビットの枠を外し、1 バイトを取り出す。"
			}
		case 3, 4, 5, 6, 7:
			step.Active = false
			step.Note = p.L7Name + " は IP を使わないため、L3〜L7 の処理は無い。"
		}
		step.TotalBytes = remaining
		step.Structure = structure
		steps = append(steps, step)
	}
	return steps
}

func Decapsulate(req Request) []Step {
	req = normalize(req)
	p := ProtocolByKey(req.Protocol)
	if p.Family == "serial" {
		return decapsulateSerial(req, p)
	}
	isPing := p.Transport == "ICMP"

	// 送信側で積み上がった最終サイズから逆算して「フル装備」の初期値を求める。
	transportBytes := tcpHeaderBytes
	if p.Transport == "UDP" {
		transportBytes = udpHeaderBytes
	}
	l3Extra := 0 // ping の場合の ICMP ヘッダ
	if isPing {
		transportBytes = 0
		l3Extra = icmpHeaderBytes
	}
	full := len(req.Message) + transportBytes + l3Extra + ipHeaderBytes + ethHeaderBytes + ethTrailer
	remaining := full
	structure := "[Eth [IP [TCP [Data]]] FCS]"
	if p.Transport == "UDP" {
		structure = "[Eth [IP [UDP [Data]]] FCS]"
	}
	if isPing {
		structure = "[Eth [IP [ICMP [Data]]] FCS]"
	}

	steps := make([]Step, 0, len(Layers))
	port := strconv.Itoa(p.Port)

	// Layers は L7→L1 順なので、受信側は逆順（L1→L7）に走査する。
	for i := len(Layers) - 1; i >= 0; i-- {
		l := Layers[i]
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
		case 1:
			step.Processing = []string{"物理信号をビット列に復元"}
			step.Note = "媒体から届いた信号をビット列に戻す。"
			step.Bitstream = toBits(req.Message, 4)

		case 2:
			step.Headers = map[string]string{
				"dstMac":    "AA:BB:CC:DD:EE:01",
				"srcMac":    "11:22:33:44:55:66",
				"etherType": "0x0800",
				"fcs":       "OK",
			}
			step.HeaderBytes = ethHeaderBytes + ethTrailer
			remaining -= ethHeaderBytes + ethTrailer
			if isPing {
				structure = "[IP [ICMP [Data]]]"
			} else {
				structure = fmt.Sprintf("[IP [%s [Data]]]", p.Transport)
			}
			step.Note = "宛先 MAC が自分宛かを確認し、FCS で誤りがないか検査。Ethernet ヘッダを外す。"

		case 3:
			step.Headers = map[string]string{
				"dstIp":    req.DstIP,
				"srcIp":    req.SrcIP,
				"ttl":      "63",
				"protocol": p.L3Protocol,
			}
			if isPing {
				remaining -= ipHeaderBytes + icmpHeaderBytes
				step.HeaderBytes = ipHeaderBytes + icmpHeaderBytes
				step.Headers["icmpType"] = "8 (Echo Request)"
				structure = "[Data]"
				step.Note = "宛先 IP を確認。ICMP Echo Request と判別し、Echo Reply を返す（擬似）。L4 は無い。"
			} else {
				remaining -= ipHeaderBytes
				step.HeaderBytes = ipHeaderBytes
				structure = fmt.Sprintf("[%s [Data]]", p.Transport)
				step.Note = "宛先 IP が自分宛かを確認し、上位プロトコルを判別。IP ヘッダを外す。"
			}

		case 4:
			if isPing {
				step.Active = false
				step.Note = "ICMP は L4 を使わないため、この層の処理は無い。"
			} else {
				hb := tcpHeaderBytes
				order := "順序を確認して"
				if p.Transport == "UDP" {
					hb = udpHeaderBytes
					order = "（UDP は順序保証なし）"
					step.Headers = map[string]string{"dstPort": port, "srcPort": "49152", "length": "8+data"}
				} else {
					step.Headers = map[string]string{"dstPort": port, "srcPort": "49152", "seq": "1", "flags": "PSH,ACK"}
				}
				step.HeaderBytes = hb
				remaining -= hb
				structure = "[Data]"
				step.Note = fmt.Sprintf("宛先ポート %d から対応アプリへ振り分け（逆多重化）。%s%s ヘッダを外す。", p.Port, order, p.Transport)
			}

		case 5:
			if isPing {
				step.Active = false
				step.Note = "Ping では使用しない。"
			} else {
				step.Processing = []string{"対応するセッションへ紐付け"}
				step.Note = "独立したヘッダは無い。"
			}

		case 6:
			if isPing {
				step.Active = false
				step.Note = "Ping では使用しない。"
			} else {
				step.Processing = []string{"文字コード復元 (UTF-8)"}
				if p.TLS {
					step.Processing = append([]string{"TLS 復号"}, step.Processing...)
				}
				step.Note = "送信側で行った変換を元に戻す。"
			}

		case 7:
			if isPing {
				step.Active = false
				step.Note = "Ping はアプリケーション層プロトコルを使わない。"
			} else {
				step.Headers = l7Headers(p, req.DstIP)
				step.Note = "アプリケーション（" + p.L7Name + "）が最終的にデータ本体を受信・解釈する。"
			}
		}

		step.TotalBytes = remaining
		step.Structure = structure
		steps = append(steps, step)
	}

	return steps
}
