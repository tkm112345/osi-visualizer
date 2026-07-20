# Contributing to OSI Model Visualizer

このプロジェクトへの貢献に興味を持っていただきありがとうございます。気軽に Issue / Pull Request を送ってください。

## 開発環境

- Go 1.24+
- Node.js 20+ / npm

## セットアップと起動

```bash
# バックエンド (:8080)
cd backend
go run .

# フロントエンド (:5173) ※別ターミナル
cd frontend
npm install
npm run dev
```

ブラウザで http://localhost:5173 を開く。

## 変更を送る前に

以下がすべて通ることを確認してください。

```bash
# バックエンド
cd backend
go vet ./...
go test ./...
go build ./...

# フロントエンド
cd frontend
npm run build   # tsc の型チェック + vite build
```

## Pull Request の流れ

1. リポジトリを fork し、`main` から作業ブランチを作成する
   （例: `feat/decapsulation-view`, `fix/l4-header-bytes`）
2. 変更はできるだけ小さく、目的が 1 つになるようにする
3. 上記のチェックがすべて通ることを確認する
4. わかりやすいタイトルと説明を付けて PR を作成する

## コミットメッセージ

- 1 行目は簡潔な要約（日本語 / 英語どちらでも可）
- 何を・なぜ変えたかが伝わるように書く

## コーディング方針

- **Go**: 標準ライブラリ中心。`gofmt` で整形する。
- **TypeScript / React**: 既存のコンポーネント構成・命名に合わせる。`strict` モードを維持する。
- 教育目的のプロジェクトなので、**正確さ**と**わかりやすさ**を優先する。
  OSI の説明を追加・修正する場合は、可能な範囲で出典（RFC / 教科書等）を PR に添えてください。

## バグ報告・提案

[Issue](../../issues) からお願いします。バグの場合は再現手順・期待する挙動・実際の挙動を書いてください。
