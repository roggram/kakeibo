# 06. デザインシステム

> この章の根拠となるソースコード: `srv/static/index.html` の `<style>` ブロック（L14–137）。
> 数値は CSS から抽出した実値です。色やサイズで迷ったら **`index.html` の CSS を直接見てください**。

---

## カラー

### ベースカラー（ダークテーマ）

| 用途 | 値 | 備考 |
|---|---|---|
| 背景（最暗） | `#0f172a` | body 背景、進捗バーのトラック背景、入力欄背景 |
| カード／タイル背景 | `#1e293b` | tile / entry / bcard / sheet など |
| 強調カード背景 | `#334155` | アクティブ tile、`btn-secondary`、リネーム入力背景 |
| 境界線 | `#334155` | tile・入力欄のボーダー |
| ヘッダー下境界 | `#1e293b` | header / tabs の境界 |

### テキストカラー

| 用途 | 値 |
|---|---|
| プライマリ文字 | `#e2e8f0`（body デフォルト） |
| 強調白文字 | `#fff` / `#f1f5f9` |
| 補助文字（明） | `#cbd5e1` |
| 補助文字（中） | `#94a3b8` |
| 補助文字（暗） | `#64748b` |
| 非活性文字 | `#475569` |

### アクセントカラー

| 用途 | 値 |
|---|---|
| プライマリ（青） | `#3b82f6` |
| 成功（緑） | `#10b981` |
| 注意（黄） | `#fbbf24` / `#f59e0b` |
| 警告（橙） | `#f97316` |
| エラー（赤） | `#ef4444` |
| 柔らかい赤 | `#f87171`（削除ボタン文字） |
| 危険ボタン背景 | `#7f1d1d`（btn-danger） |
| 危険ボタン文字 | `#fecaca` |
| サマリーカードグラデ | `linear-gradient(135deg, #1e3a8a, #312e81)` |
| クイック記録ボタン | `linear-gradient(135deg, #3b82f6, #8b5cf6)` |
| TOP8 バー | `linear-gradient(90deg, #3b82f6, #8b5cf6)` |

### 状態別文字色（アラートメッセージ）

| クラス | 色 |
|---|---|
| `.alert-msg.ok` | `#86efac`（明るい緑）|
| `.alert-msg.warn` | `#fcd34d`（明るい黄）|
| `.alert-msg.bad` | `#fca5a5`（明るい赤）|

### Theme color（PWA）

- `<meta name="theme-color" content="#10b981">`
- `manifest.theme_color: #10b981`
- `manifest.background_color: #0f172a`

---

## タイポグラフィ

**フォントスタック**:
```
-apple-system, BlinkMacSystemFont, 'Hiragino Sans', sans-serif
```
iOS では San Francisco（システムフォント）→ Hiragino Sans へフォールバック。SwiftUI では `.body` 等のシステムフォントを使えば自動で合致します。

**サイズ**（ピクセル指定）:

| 用途 | サイズ | 太さ | letter-spacing |
|---|---|---|---|
| ヘッダー合計金額 (`.total b`) | 20px | 700 | `-0.5px` |
| ヘッダー小ラベル (`.total`) | 11px | normal | - |
| カテゴリタブ | 12px | normal | - |
| ボトムナビ文字 | 11px | normal | - |
| ボトムナビアイコン | 20px | - | - |
| タイル絵文字 | 32px | - | - |
| タイル名 | 13px | 600 | - |
| タイル価格 | 11px | - | - |
| バッジ（使用回数）| 10px | 700 | - |
| トースト | 14px | 600 | - |
| トーストボタン | 14px | 700 | - |
| エントリ名 | 14px | normal | - |
| エントリ時刻・カテゴリ | 11px | - | - |
| エントリ金額 | 14px | 600 | - |
| 削除ボタン | 12px | 600 | - |
| 月ヘッダー | 13px | 700 | - |
| 日ヘッダー | 12px | 600 | - |
| 期間タブ | 12px | 600 | - |
| 予算バー金額 (`.row b`) | 14px | - | - |
| 予算バーラベル (`.row`) | 12px | - | - |
| 予算バーメッセージ | 11px | - | - |
| ミニバー | 11–12px | - | - |
| サマリー大数値 | 26px | 700 | - |
| サマリー小ラベル | 12px | - | - |
| 統計カード見出し (`h4`) | 13px | 600 | - |
| 統計カード大数値 | 28px | 700 | - |
| バー行 (`.bar`) | 13px | - | - |
| モーダルタイトル (`h2`) | 16px | - | - |
| 入力欄 | 15px | - | - |
| 絵文字ピッカーボタン | 22px | - | - |
| ステッパー `−` / `＋` | 48px | 700 | - |
| ステッパー qty 数字 | 36px | 800 | - |
| ステッパー計算式 | 12px | - | - |
| ステッパー合計 | 18px | 700 | - |
| アクションボタン | 14px | 600 | - |
| `.field label` | 12px | - | - |

---

## スペーシング・角丸

| 用途 | 値 |
|---|---|
| body padding（main） | `12px`（下は 80px、ボトムナビ分） |
| header padding | `12px 16px` |
| カテゴリタブ padding | `8px 10px`（tab 内側 `6px 4px`） |
| ボトムナビ padding | `8px 0 + safe-area-inset-bottom` |
| タイルグリッド gap | `10px` |
| タイル padding | `8px` |
| エントリ padding | `10px 12px` |
| 予算バー padding | `10px 12px` |
| カテゴリカード (`.bcard`) padding | `14px` |
| サマリーカード padding | `16px` |
| 統計カード padding | `16px` |
| モーダルシート padding | `20px` |

**角丸（border-radius）**:

| 用途 | 値 |
|---|---|
| タイル | `16px` |
| エントリ・予算バー・期間タブ親 | `10–12px` |
| ミニバー（`.budget-mini .m`） | `10px` |
| カテゴリカード (`.bcard`) | `14px` |
| サマリーカード | `16px` |
| 統計カード | `14px` |
| モーダルシート上端 | `20px 20px 0 0` |
| タブピル（`.tab`） | `999px`（完全な丸） |
| ボタン一般 | `8–10px` |
| トースト | `14px` |
| 進捗バー | `2–5px`（高さに応じて） |

**ボーダー**:

| 用途 | 値 |
|---|---|
| タイル | `1px solid #334155` |
| 予算バー左 | `4px solid {color}` |
| ミニバー左 | `3px solid {color}` |
| 入力欄 | `1px solid #334155` |
| ヘッダー下 | `1px solid #1e293b` |
| ボトムナビ上 | `1px solid #334155` |

---

## タイルグリッド

- `display: grid; grid-template-columns: repeat(3, minmax(0, 1fr))`
- `gap: 10px`
- 各タイル: `aspect-ratio: 1`（正方形）
- タイル内部は flex 縦並びで中央寄せ
- 「＋ 追加」タイル: 背景 transparent、`border: 2px dashed #334155`

---

## 予算バーの色閾値

進捗率による色マッピング:

| 進捗率 | クラス | バー色 | 左ボーダー色 |
|---|---|---|---|
| `pct < 50%` | （無印）| `#3b82f6` | `#3b82f6` |
| `50% <= pct < 80%` | `.mid` | `#f59e0b` | `#f59e0b` |
| `80% <= pct < 100%` | `.warn` | `#f97316` | `#f97316` |
| `pct >= 100%` (= `spent >= limit`) | `.danger` | `#ef4444` | `#ef4444` |

予算カード (`.bcard`) では `spent >= limit` で `.over` クラス（赤ボーダー枠）、`80% <=` で `.warn` クラス（黄ボーダー枠）。

---

## アニメーション

| 要素 | アニメーション |
|---|---|
| タイル active | `transform: scale(0.94); background: #334155;`（押下時）|
| タイル flash | `.flash` クラスで背景 `#10b981`、`transition: background .15s`、250ms 後にクラス除去 |
| タイル一般 | `transition: transform .08s, background .15s` |
| 予算バー幅 (`.pb div`) | `transition: width .3s, background .3s` |
| 予算カード進捗バー幅 (`.bcard .pb div`) | `transition: width .4s` |
| トースト | `transition: opacity .2s, transform .2s`、`transform: translateY(8px) → 0` でスライドアップ |
| モーダルシート | `@keyframes slideUp` で 250ms、`translateY(100%) → 0` |
| ステッパーボタン active | `transition: background .1s` |

---

## レイアウト・セーフエリア

- `body { max-width: 480px; margin: 0 auto; }`（PC でも携帯幅で中央表示）
- 高さ: `100vh` フォールバック + `100dvh`（dynamic viewport、iOS Safari の URL バー対応）
- `overscroll-behavior: none`（バウンス抑止）
- ボトムナビ: `padding-bottom: calc(8px + env(safe-area-inset-bottom))`（ホームバー領域回避）
- `touch-action: pan-y`（横スワイプ防止）

---

## モード

- **ダークモード固定**（ライトモード未実装）。`prefers-color-scheme` には反応しない。
- iOS 実装では当面 `.preferredColorScheme(.dark)` で固定推奨。

---

## タップ領域

| 要素 | 最低サイズ |
|---|---|
| 履歴の削除ボタン (`.entry .del`) | `min-height: 40px; min-width: 48px`（44pt 推奨だが現状 40px）|
| トーストの取り消しボタン | `min-height: 44px; padding: 12px 20px` |
| タイル | 約 `120 × 120px`（3列、480px 幅基準）|
| ステッパー `−`/`＋` | 高さ 120px × 半分（横幅約 240px の左右半分）|

iOS の HIG では最低 44pt 推奨。削除ボタンは `min-height: 44pt` に上げることを推奨。

---

## 長押し時のテキスト選択抑止

タイルとステッパーには:
```
user-select: none;
-webkit-user-select: none;
-webkit-touch-callout: none;
-webkit-tap-highlight-color: transparent;
```

`contextmenu` イベントも `preventDefault()` で抑止（iOS Safari の長押し選択メニューを出さないため）。

iOS ネイティブでは `UILongPressGestureRecognizer` を使えばこの問題は発生しない。
