import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    // 開発時は /api をバックエンド(:8080)へプロキシする。
    // 本番(Docker)では nginx が同じく /api をプロキシするため、
    // フロントは常に相対パス /api を叩けばよい。
    proxy: {
      "/api": "http://localhost:8080",
    },
  },
});
