// localStorage の読み書き。API ができたら、この関数を使っている lib/ の処理ごと差し替える。

const PREFIX = "eipiku:";

export function readStore<T>(key: string, fallback: T): T {
  try {
    const raw = window.localStorage.getItem(PREFIX + key);
    return raw ? (JSON.parse(raw) as T) : fallback;
  } catch {
    return fallback;
  }
}

export function writeStore<T>(key: string, value: T) {
  try {
    window.localStorage.setItem(PREFIX + key, JSON.stringify(value));
  } catch {
    // 保存できない環境（プライベートモードなど）では画面上の状態だけで続ける
  }
}

// DB の DATETIME と同じ "YYYY-MM-DD HH:MM:SS" 形式
export function now() {
  const d = new Date();
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}
