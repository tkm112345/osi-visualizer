# OSI Model Visualizer

[![CI](https://github.com/tkm112345/osi-visualizer/actions/workflows/ci.yml/badge.svg)](https://github.com/tkm112345/osi-visualizer/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

OSI 参照モデルを学ぶための Web アプリ。データが **L7 → L1 へカプセル化**されていく様子を可視化し、各レイヤーをクリックすると付与されるヘッダや処理内容が確認できる。

- **フロントエンド**: React + Vite + TypeScript
- **バックエンド**: Go（標準ライブラリのみ）

## 特徴

- L1〜L7 を下から積み上げ表示。「送信」ボタンでカプセル化アニメーションが走る。
- **L5（セッション）/ L6（プレゼンテーション）も表示**。これらは独立ヘッダを付けず「処理内容」を明示することで、
  OSI（理論）と TCP/IP（実装）のギャップまで学べる。
- 各レイヤーのヘッダ・PDU・累積バイト数・カプセル化構造 `[Eth [IP [TCP [Data]]]]` を表示。

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
| POST | `/api/encapsulate` | `{message, srcIp, dstIp}` を受け取り L7→L1 の各ステップを返す |

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
