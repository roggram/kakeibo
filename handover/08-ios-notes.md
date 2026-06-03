# 08. iOS 実装上の示唆・推奨スタック

> この章の根拠となるソースコード: `srv/static/index.html` 全体、`srv/server.go`、`db/migrations/`。
> Web 版の挙動の正確な再現を目指す場合は **`index.html` を直接読んで実装してください**。

---

## 推奨技術スタック

| レイヤ | 推奨 | 補足 |
|---|---|---|
| UI | **SwiftUI**（iOS 17+ ターゲット推奨） | モーダル・グリッド・アニメーションが宣言的に書ける |
| 状態管理 | **`@Observable`（iOS 17+）** または `ObservableObject` | Combine は補助的に |
| 永続化 | **SwiftData** | iOS 17+。シンプルな単一ストアに最適 |
| 永続化（代替）| Core Data / GRDB / Realm | チームの慣れに応じて |
| Haptics | `UIImpactFeedbackGenerator` / `UINotificationFeedbackGenerator` | SwiftUI なら `.sensoryFeedback(...)` |
| カラー / フォント | Asset Catalog の Color Set + SwiftUI `Color(...)` | Dark 固定で 1 セットのみ定義 |

iOS 16 もターゲットにする場合: SwiftData は使えないので Core Data または独自 JSON ファイル永続化（後述）。

---

## データ層の移植方針

### オプション 1: Web 版と同じ単一 JSON ステート（推奨：初手）

Web 版は `{products, entries, budgets, notified, categories}` の 1個の JSON を 1行に保存しているだけ。これをそのまま iOS でも踏襲できます。

- `Documents/app_state.json` に 1ファイルで保存
- 起動時にロード、変更時に 300ms デバウンスで非同期書き込み
- 構造体は Web 版のフィールド名と合わせる（後述）
- 同期したいなら **iCloud Drive コンテナ** または **CloudKit private DB** に置き換え可能

メリット:
- 短時間で動く版が作れる
- Web 版とサーバー API（`/api/state`）と互換性が保てる（将来同期するなら直接 `URLSession.upload`）
- バグ時に JSON ダンプの目視デバッグが容易

デメリット:
- エントリ数が数万件オーダーになるとパフォーマンス悪化（実用上は十分大丈夫）

### オプション 2: SwiftData / Core Data に正規化

`Product` / `Entry` を Model、`budgets` / `notified` / `categories` は別途キー・バリュー（`AppSettings` 等の単一モデル）に保存。

メリット: 大量データ・高速クエリ・将来の機能拡張に強い。
デメリット: 初期コストが高い、Web 版との JSON 互換性は失われる。

**最初は オプション 1 → 後で必要に応じて 2 へ移行** が現実的です。

---

## Swift モデル定義（Web 版互換のシリアライズ）

```swift
struct AppState: Codable {
    var products: [Product] = []
    var entries: [Entry] = []
    var budgets: [String: Int] = [:]
    var notified: [String: [Int]] = [:]
    var categories: [String] = ["食品","日用品","交通","娯楽","その他"]
}

struct Product: Codable, Identifiable, Hashable {
    var id: String          // "x" + 7文字 互換 (UUID().uuidString に置き換え可)
    var emoji: String
    var name: String
    var price: Int
    var category: String
}

struct Entry: Codable, Identifiable, Hashable {
    var id: String
    var productId: String?  // null 許容（メモで記録）
    var name: String
    var emoji: String
    var price: Int
    var qty: Int
    var category: String
    var date: String        // "YYYY-MM-DD"
    var time: Int64         // ミリ秒タイムスタンプ
}
```

`notified` のキー形式は `"YYYY-MM-カテゴリ名"`（例: `"2025-06-食品"`）。詳細は `03-data-model.md` 参照。

---

## Haptics

| トリガ | パターン | Swift |
|---|---|---|
| タイル 1タップ記録 | 軽い | `UIImpactFeedbackGenerator(style: .light).impactOccurred()` |
| 長押し発動（500ms 経過時）| 中 | `UIImpactFeedbackGenerator(style: .medium).impactOccurred()` |
| ステッパー / ×N ボタン操作 | 軽い | `.light` |
| 予算超過拒否 | 警告 | `UINotificationFeedbackGenerator().notificationOccurred(.warning)` |
| 50%/80%/100% アラート | 警告 | `.warning` |
| 記録成功 | 成功 | `.success`（任意）|

SwiftUI 17+ なら `.sensoryFeedback(.impact(weight: .light), trigger: tapCount)` 形式が簡潔。

---

## PWA → iOS で失われる挙動の代替

| Web 版 | iOS 代替 |
|---|---|
| `navigator.sendBeacon` on `beforeunload` | `applicationWillResignActive` / `scenePhase == .inactive` で同期保存 |
| `manifest.webmanifest` のホーム画面追加 | Info.plist でアプリ名・アイコン定義（App Store / TestFlight 経由）|
| `theme-color` メタタグ | UIStatusBarStyle / SwiftUI `.preferredColorScheme(.dark)` |
| `confirm()` / `alert()` | `.alert(...)` / `.confirmationDialog(...)` |
| `vibrate()` | UIFeedbackGenerator（上記）|
| `setTimeout` debounce | `Task` + `Task.sleep` + 直前タスクキャンセル、または Combine `debounce` |

**保存タイミングの推奨実装**:
```swift
// 状態変更ごとに save() を呼ぶ
// save() は 300ms デバウンス
// scenePhase が .background / .inactive に遷移したら即時 flush
```

---

## Dynamic Type 対応

Web 版はピクセル固定（タイル名 13px、ヘッダー合計 20px 等）。iOS では **Dynamic Type 対応の有無を要検討**:

| 方針 | メリット | デメリット |
|---|---|---|
| 完全に Dynamic Type 対応（`.font(.body)` 等）| アクセシビリティ◎ | グリッドが崩れる、3列タイルが破綻する可能性 |
| 部分対応（ヘッダー・サマリーは Dynamic、タイル名は固定）| バランス | 設計の手間 |
| 完全固定（`.font(.system(size: 13))`）| Web 版と完全一致 | アクセシビリティ△ |

推奨: **タイル名・絵文字は固定、それ以外（履歴・予算カード・サマリー）は Dynamic Type 対応**。最初は固定で出して、後でユーザー要望に応じて拡張する手もあり。

---

## Safe Area・レイアウト

- **max-width 480px** の中央寄せは iPhone では不要（画面幅いっぱい使う）。iPad はオプションで考慮。
- ボトムナビ: SwiftUI の `TabView`（標準）を使えば Safe Area 自動処理。Web 版のような自前 fixed バーを再現するなら `safeAreaInset(edge: .bottom)`。
- iPhone のホームバー領域は SwiftUI が自動で扱う。

---

## Dark Mode 固定（一旦）

- ルートに `.preferredColorScheme(.dark)` を付ける
- Asset Catalog の Color Set は Dark のみ定義（または Any Appearance = Dark）
- 後でライトモード対応するときに Color Set に Light バリアントを追加

---

## 横画面対応

Web 版は縦画面前提（max-width: 480px のせいで横でも縦表示）。

- iPhone は **Portrait のみサポート**（Info.plist `UISupportedInterfaceOrientations` で `Portrait` のみ）が初期推奨
- iPad は将来検討（タイルを 5–6 列にする、サイドバーで履歴と新規を並べる 等）
- **ユーザーに横画面サポートの要否を確認すること**

---

## Live Activities / Widget（将来拡張）

Web 版にはない要素ですが、サクサク家計簿との相性は良いです:

- **Widget（ホーム画面ウィジェット）**: 今月の予算進捗バーを表示。タップで新規タブを開く
- **Live Activity（Dynamic Island）**: 「今日の合計 ¥XXX、予算残り ¥XXX」を表示
- **Lock Screen Widget**: 今月の予算ゲージ

ただし優先度は **低**。コア機能を完成させてからユーザーと相談してください。

---

## iCloud 同期

Web 版はサーバー保存（SQLite）でデバイス間自動同期されている状態。iOS 版でも何らかの同期が欲しい:

| 方式 | 推奨度 | 補足 |
|---|---|---|
| **iCloud Drive のドキュメント JSON** | ◎ | オプション 1 のデータ層と相性◎、ユーザー認証不要、最小実装で実現 |
| **CloudKit private DB** | ○ | きめ細かい同期可能だが実装重 |
| **Web 版のサーバー (`/api/state`) を使う** | △ | 認証なし・シングルユーザー前提で同じ DB を共有することは可能だが、競合解決ロジックが必要 |
| **同期なし（デバイスローカル）** | ○ | 初手としてはアリ |

**推奨**: 最初はデバイスローカル → ユーザーから要望が出たら iCloud Drive JSON 同期を追加。

---

## アクセシビリティ

- VoiceOver ラベル: 各タイルに「{name}、{price}円、これまで {count}回購入、ダブルタップで記録」
- 動的フォントサイズ: 上記 Dynamic Type 方針に従う
- カラーコントラスト: 予算バー色は AA 基準で見直し（特に黄 `#f59e0b` は背景 `#1e293b` でギリギリ）
- Reduce Motion 対応: `accessibilityReduceMotion` が true ならタイル flash・モーダル slide を無効化

---

## まず作る順序の提案

1. **データモデル + 永続化（オプション 1: JSON ファイル）** を先に作ってユニットテストで `categorySpent` / `checkBudget` / `triggerAlerts` を緑にする
2. **新規タブ** だけ動かす（タイルグリッド + 1タップ記録 + Undo トースト）
3. **長押しメニュー** を追加
4. **履歴タブ**
5. **予算タブ + 予算編集モーダル**
6. **状況タブ + 集計タブ**
7. **カテゴリ管理モーダル**
8. **アラート（50/80/100%）** を仕上げ
9. （任意）iCloud 同期 / Widget

各画面の実装着手前に **ユーザーへスクリーンショットを要求すること**（`00-INDEX.md` の方針）。

---

## デバッグ・検証用

Web 版が動いている URL: `https://shop-log.exe.xyz:8000/`

`GET /api/state` でサーバー上の本物のデータが見えます。iOS 版の出力 JSON と比較してフィールド・形式が一致しているかチェックしてください。
