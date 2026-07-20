# デモ GIF の作り方 / Recording the demo GIF

README のデモ GIF (`docs/demo.gif`) を差し替える手順です。
This describes how to produce `docs/demo.gif` used in the README.

## 1. アプリを起動 / Start the app

```bash
docker compose up --build   # または backend/frontend を個別に起動
```

ブラウザで <http://localhost:8080>（Docker）または <http://localhost:5173>（dev）を開く。

## 2. 録画 / Record

おすすめの流れ（15〜20秒程度）:

1. プロトコルで `HTTPS` を選ぶ（L6 の TLS 暗号化が見える）
2. 「擬似送信 ▶ シミュレート」を押してカプセル化→デカプセル化のアニメを見せる
3. どれかのレイヤーで「この層での実データを見る」を開く
4. 右上の言語スイッチで `English` に切り替える

録画ツールの例:

- **macOS**: [Kap](https://getkap.co/) や QuickTime → GIF 変換
- **Linux**: [Peek](https://github.com/phw/peek)（ウィンドウ選択で直接 GIF 出力）
- **CLI**: 画面を mp4 で録ってから
  ```bash
  ffmpeg -i demo.mp4 -vf "fps=12,scale=900:-1:flags=lanczos" -loop 0 docs/demo.gif
  ```

## 3. 配置 / Place

生成した GIF を `docs/demo.gif` として置き、コミットすれば README に表示されます。
（幅 800〜1000px、10MB 未満を目安に。）
