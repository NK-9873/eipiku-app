import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // 本番の Docker イメージ用。動かすのに必要なファイルだけを .next/standalone にまとめる（frontend/Dockerfile）
  output: "standalone",
  experimental: {
    agentFeedback: true,
    serverActions: {
      // 管理画面の単語登録で、写真とイラスト（各5MBまで）を1回の Server Action で送るため
      bodySizeLimit: "12mb",
    },
  },
  cacheComponents: true,
  partialPrefetching: true,
  turbopack: {
    rules: {
      "*.css": {
        loaders: ["@tailwindcss/turbopack"],
        as: "*.css",
      },
    },
  },
};

export default nextConfig;
