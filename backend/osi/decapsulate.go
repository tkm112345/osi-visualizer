package osi

// Decapsulate は受信ホスト視点で、届いたフレームを L1 → L7 へ
// デカプセル化（ヘッダを解析して外していく）各ステップを返す。
// これは Encapsulate と対称で、あくまで擬似的なシミュレーション。
func Decapsulate(req Request) []Step {
	if req.SrcIP == "" {
		req.SrcIP = "192.168.0.10"
	}
	if req.DstIP == "" {
		req.DstIP = "93.184.216.34"
	}

	payloadBytes := len(req.Message)
	// 受信時は「フル装備のフレーム」が届いた状態から始まる。
	full := payloadBytes + tcpHeaderBytes + ipHeaderBytes + ethHeaderBytes + ethTrailer
	remaining := full
	structure := "[Eth [IP [TCP [Data]]] FCS]"

	steps := make([]Step, 0, len(Layers))

	// Layers は L7→L1 順なので、受信側は逆順（L1→L7）に走査する。
	for i := len(Layers) - 1; i >= 0; i-- {
		l := Layers[i]
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
			structure = "[IP [TCP [Data]]]"
			step.Note = "宛先 MAC が自分宛かを確認し、FCS で誤りがないか検査。Ethernet ヘッダを外す。"

		case 3:
			step.Headers = map[string]string{
				"dstIp":    req.DstIP,
				"srcIp":    req.SrcIP,
				"ttl":      "63",
				"protocol": "TCP(6)",
			}
			step.HeaderBytes = ipHeaderBytes
			remaining -= ipHeaderBytes
			structure = "[TCP [Data]]"
			step.Note = "宛先 IP が自分宛かを確認し、上位プロトコルが TCP であることを判別。IP ヘッダを外す。"

		case 4:
			step.Headers = map[string]string{
				"dstPort": "80",
				"srcPort": "49152",
				"seq":     "1",
				"flags":   "PSH,ACK",
			}
			step.HeaderBytes = tcpHeaderBytes
			remaining -= tcpHeaderBytes
			structure = "[Data]"
			step.Note = "宛先ポート番号から対応するアプリへ振り分け（逆多重化）。順序を確認して TCP ヘッダを外す。"

		case 5:
			step.Processing = []string{"対応するセッションへ紐付け"}
			step.Note = "独立したヘッダは無い。TCP/IP モデルでは Application 層に含まれる。"

		case 6:
			step.Processing = []string{"TLS 復号", "文字コード復元 (UTF-8)"}
			step.Note = "独立したヘッダは無い。送信側で行った変換を元に戻す。"

		case 7:
			step.Headers = map[string]string{
				"protocol": "HTTP",
				"method":   "GET",
				"host":     req.DstIP,
			}
			step.Note = "アプリケーションが最終的にデータ本体を受信・解釈する。"
		}

		step.TotalBytes = remaining
		step.Structure = structure
		steps = append(steps, step)
	}

	return steps
}
