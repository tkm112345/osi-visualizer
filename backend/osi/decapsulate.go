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

	// frame: 受信したフレーム全体（フレーミング + データ）。L2 で枠を外す。
	cur := []FramePart{
		hdr(tx(p.L7Name+" フレーミング", p.L7Name+" framing"), serialL2Detail(p), framing),
		{Label: payloadLabel(p), Detail: txSame(req.Message), Kind: "payload", Bytes: len(req.Message)},
	}
	snapshot := func() []FramePart { return append([]FramePart{}, cur...) }

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
				step.Note = tx("自分のアドレス宛か確認し、ACK を返してデータを受け取る。",
					"Checks whether the address matches, returns ACK, and receives the data.")
			} else {
				step.Note = tx("Start/Stop ビットの枠を外し、1 バイトを取り出す。",
					"Removes the start/stop bit framing and extracts each byte.")
			}
			cur = cur[1:] // フレーミングを外す
		case 3, 4, 5, 6, 7:
			step.Active = false
			step.Note = tx(
				p.L7Name+" は IP を使わないため、L3〜L7 の処理は無い。",
				p.L7Name+" does not use IP, so there is no L3–L7 processing.")
		}
		step.TotalBytes = remaining
		step.Structure = structure
		if step.Active {
			step.Frame = snapshot()
		}
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

	// frame: 受信フレーム全体を組み立て、上位層へ進むごとにヘッダを外していく。
	plDetail := txSame(req.Message)
	plLabel := payloadLabel(p)
	if p.TLS {
		plDetail = tlsEncrypted(req.Message)
		plLabel = tx("ペイロード (TLS暗号化)", "Payload (TLS-encrypted)")
	}
	cur := []FramePart{
		hdr(ethHdrLabel(), txSame(ethDetail()), ethHeaderBytes),
		hdr(ipHdrLabel(), txSame(ipDetail(req, p)), ipHeaderBytes),
	}
	if isPing {
		cur = append(cur, hdr(icmpHdrLabel(), txSame(icmpDetail()), icmpHeaderBytes))
	} else if p.Transport == "UDP" {
		cur = append(cur, hdr(l4HdrLabel(p), txSame(udpDetail(p, len(req.Message))), udpHeaderBytes))
	} else {
		cur = append(cur, hdr(l4HdrLabel(p), txSame(tcpDetail(p)), tcpHeaderBytes))
	}
	cur = append(cur, FramePart{Label: plLabel, Detail: plDetail, Kind: "payload", Bytes: len(req.Message)})
	cur = append(cur, FramePart{Label: fcsLabel(), Detail: txSame("0x1A2B3C4D"), Kind: "trailer", Bytes: ethTrailer})
	snapshot := func() []FramePart { return append([]FramePart{}, cur...) }

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
			step.Processing = []Text{tx("物理信号をビット列に復元", "Restore the physical signal to bits")}
			step.Note = tx("媒体から届いた信号をビット列に戻す。",
				"Converts the signal arriving from the medium back into bits.")
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
			step.Note = tx(
				"宛先 MAC が自分宛かを確認し、FCS で誤りがないか検査。Ethernet ヘッダを外す。",
				"Checks the destination MAC and verifies the FCS, then removes the Ethernet header.")
			cur = cur[1 : len(cur)-1] // Eth ヘッダと FCS を外す

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
				step.Note = tx(
					"宛先 IP を確認。ICMP Echo Request と判別し、Echo Reply を返す（擬似）。L4 は無い。",
					"Checks the destination IP, recognizes an ICMP Echo Request, and replies (simulated). No L4.")
				cur = cur[2:] // IP ヘッダと ICMP ヘッダを外す
			} else {
				remaining -= ipHeaderBytes
				step.HeaderBytes = ipHeaderBytes
				structure = fmt.Sprintf("[%s [Data]]", p.Transport)
				step.Note = tx(
					"宛先 IP が自分宛かを確認し、上位プロトコルを判別。IP ヘッダを外す。",
					"Checks the destination IP, determines the upper protocol, and removes the IP header.")
				cur = cur[1:] // IP ヘッダを外す
			}

		case 4:
			if isPing {
				step.Active = false
				step.Note = tx("ICMP は L4 を使わないため、この層の処理は無い。",
					"ICMP does not use L4, so there is no processing at this layer.")
			} else {
				hb := tcpHeaderBytes
				order := tx("順序を確認して", "reorders segments and ")
				if p.Transport == "UDP" {
					hb = udpHeaderBytes
					order = tx("（UDP は順序保証なし）", "(UDP has no ordering) ")
					step.Headers = map[string]string{"dstPort": port, "srcPort": "49152", "length": "8+data"}
				} else {
					step.Headers = map[string]string{"dstPort": port, "srcPort": "49152", "seq": "1", "flags": "PSH,ACK"}
				}
				step.HeaderBytes = hb
				remaining -= hb
				structure = "[Data]"
				step.Note = tx(
					fmt.Sprintf("宛先ポート %d から対応アプリへ振り分け（逆多重化）。%s%s ヘッダを外す。", p.Port, order.Ja, p.Transport),
					fmt.Sprintf("Demultiplexes to the app by destination port %d, %sremoves the %s header.", p.Port, order.En, p.Transport))
				cur = cur[1:] // トランスポートヘッダを外す
			}

		case 5:
			if isPing {
				step.Active = false
				step.Note = tx("Ping では使用しない。", "Not used for Ping.")
			} else {
				step.Processing = []Text{tx("対応するセッションへ紐付け", "Associate with the matching session")}
				step.Note = tx("独立したヘッダは無い。", "There is no separate header.")
			}

		case 6:
			if isPing {
				step.Active = false
				step.Note = tx("Ping では使用しない。", "Not used for Ping.")
			} else {
				step.Processing = []Text{tx("文字コード復元 (UTF-8)", "Restore character encoding (UTF-8)")}
				if p.TLS {
					step.Processing = append([]Text{tx("TLS 復号", "TLS decryption")}, step.Processing...)
					// 暗号文のペイロードを平文に戻す。
					if n := len(cur); n > 0 {
						cur[n-1].Label = payloadLabel(p)
						cur[n-1].Detail = txSame(req.Message)
					}
				}
				step.Note = tx("送信側で行った変換を元に戻す。",
					"Reverses the transformations applied by the sender.")
			}

		case 7:
			if isPing {
				step.Active = false
				step.Note = tx("Ping はアプリケーション層プロトコルを使わない。",
					"Ping does not use an application-layer protocol.")
			} else {
				step.Headers = l7Headers(p, req.DstIP)
				step.Note = tx(
					"アプリケーション（"+p.L7Name+"）が最終的にデータ本体を受信・解釈する。",
					"The application ("+p.L7Name+") finally receives and interprets the data body.")
			}
		}

		step.TotalBytes = remaining
		step.Structure = structure
		if step.Active {
			step.Frame = snapshot()
		}
		steps = append(steps, step)
	}

	return steps
}
