# OSI 通信シミュレーター (OSI Model Visualizer)

[![CI](https://github.com/tkm112345/osi-visualizer/actions/workflows/ci.yml/badge.svg)](https://github.com/tkm112345/osi-visualizer/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

OSI 参照モデルを学ぶための Web アプリ。**送信ホスト**でデータが L7 → L1 へ**カプセル化**され、
**受信ホスト**で L1 → L7 へ**デカプセル化**される様子を可視化する。各レイヤーをクリックすると付与／除去されるヘッダや処理内容が確認できる。

> ⚠️ **これは擬似シミュレーションです。** 実際にネットワークへパケットを送信することは一切ありません。
> 宛先 IP は表示上のラベルで、ソケットも開かず通信は発生しません。ローカルの Go サーバが教材データを計算して返すだけです。

- **フロントエンド**: React + Vite + TypeScript
- **バックエンド**: Go（標準ライブラリのみ）

## 特徴

- **18 種類のプロトコルを選択でき**、選択に応じてスタック構成とペイロードが変わる:
  - Web: `HTTP` / `HTTPS` / `WebSocket`（HTTP/HTTPS のペイロードは **HTML テンプレート**、HTTPS は L6 で TLS）
  - ファイル/メール/リモート: `FTP` / `SMTP` / `SSH`
  - IoT/メッセージング: `MQTT` / `CoAP`
  - メディア: `RTSP` / `RTP`、インフラ: `DNS` / `DHCP` / `NTP` / `SNMP`、診断: `Ping`
  - **シリアル通信（L1-L2 のみ・IP なし）**: `UART` / `I2C` / `SPI`
- 選択で L4 が **TCP↔UDP** に変化。`Ping` は **L4・L7 を使わず** L3 に **ICMP**（`[IP [ICMP [Data]]]`）。
  `UART`/`I2C`/`SPI` は **L3〜L7 を使わず** L1/L2 だけでフレーム化・信号化する。
- 送信ホスト A（L7→L1）と受信ホスト B（L1→L7）を左右に並べ、送信→受信を連続アニメーション（途中で **停止** も可能）。
- 起動直後から両ホストのレイヤースタックを表示。
- L1〜L7 を積み上げ表示。カプセル化でヘッダが増え、デカプセル化で外れていく様子が見える。
- **L5（セッション）/ L6（プレゼンテーション）も表示**。これらは独立ヘッダを付けず「処理内容」を明示することで、
  OSI（理論）と TCP/IP（実装）のギャップまで学べる。
- 受信側は宛先 MAC/IP の確認・ポートによるアプリ振り分けなど、**受信ホスト特有の処理**も表示。
- 各レイヤーのヘッダ・PDU・累積バイト数・構造 `[Eth [IP [TCP [Data]]]]` を表示。

### カプセル化のイメージ

```
L7 Application  Data     +0   →  5B   [Data]
L6 Presentation Data     +0   →  5B   [Data]          ← ヘッダなし（処理のみ）
L5 Session      Data     +0   →  5B   [Data]          ← ヘッダなし（処理のみ）
L4 Transport    Segment +20   → 25B   [TCP [Data]]
L3 Network      Packet  +20   → 45B   [IP [TCP [Data]]]
L2 Data Link    Frame   +18   → 63B   [Eth [IP [TCP [Data]]] FCS]
L1 Physical     Bits     +0   → 63B   (ビット列に変換)
```

## 必要環境

- Go 1.24+
- Node.js 20+ / npm

## 起動方法

ターミナルを 2 つ使う。

### 1. バックエンド (:8080)
```bash
cd backend
go run .
```

### 2. フロントエンド (:5173)
```bash
cd frontend
npm install   # 初回のみ
npm run dev
```

ブラウザで <http://localhost:5173> を開き、「送信」を押す。各レイヤーをクリックすると詳細が表示される。

## テスト

```bash
cd backend && go test ./...
```

## API

| メソッド | パス | 説明 |
|---------|------|------|
| GET  | `/api/layers` | 全 7 層の静的メタ情報 |
| GET  | `/api/protocols` | 選択可能なプロトコル一覧（HTTP/HTTPS/DNS/RTSP/RTP/Ping） |
| POST | `/api/encapsulate` | `{message, srcIp, dstIp, protocol}` を受け取り送信側 L7→L1 の各ステップを返す |
| POST | `/api/decapsulate` | 同じ入力から受信側 L1→L7 の各ステップ（擬似）を返す |

## プロジェクト構成

```
osi-visualizer/
├── backend/           # Go（標準ライブラリのみ）
│   ├── main.go        # HTTP サーバ・CORS・2 エンドポイント
│   └── osi/           # レイヤー定義とカプセル化ロジック（+ テスト）
└── frontend/          # React + Vite + TypeScript
    └── src/
        ├── App.tsx
        └── components/  # LayerStack / LayerCard / PacketDetail
```

## コントリビュート

歓迎します。[CONTRIBUTING.md](./CONTRIBUTING.md) と [行動規範](./CODE_OF_CONDUCT.md) を参照してください。

## ライセンス

[MIT](./LICENSE)
