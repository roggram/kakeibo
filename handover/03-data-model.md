# 03. データモデル・スキーマ・状態遷移

> この章の根拠となるソースコード:
> - `srv/static/index.html`（状態変数・`loadState`・`save`・`uid` 関数、L160–235）
> - `srv/server.go`（`handleGetState`・`handlePutState`）
> - `db/migrations/002-state.sql`
>
> 仕様書に書かれていない振る舞いは **`index.html` を直接読んで実装してください**。

---

## DB スキーマ

```sql
-- db/migrations/002-state.sql
CREATE TABLE IF NOT EXISTS app_state (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    data TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

- テーブルには **常に 1行だけ**（`id = 1` の CHECK 制約）。
- `data` 列に JSON テキスト全体を格納。
- `PUT /api/state` は UPSERT（`INSERT ... ON CONFLICT DO UPDATE`）。
- `GET /api/state` は行なし（初回）の場合 `"null"` を返す。

---

## JSON ステートの全体形

```json
{
  "products": [ ...Product[] ],
  "entries":  [ ...Entry[] ],
  "budgets":  { "カテゴリ名": 金額(number), ... },
  "notified": { "YYYY-MM-カテゴリ名": [閾値, ...], ... },
  "categories": [ "食品", "日用品", ... ]
}
```

---

## エンティティ詳細

### Product（商品）

| フィールド | 型 | 説明 |
|---|---|---|
| `id` | `string` | `uid()` で生成。形式: `"x"` + 36進数 7文字（例: `"x4a2b8c1"`）|
| `emoji` | `string` | タイルに表示する絵文字。EMOJIS 定数（L155）のいずれかが推奨、ユーザーが選択 |
| `name` | `string` | 商品名。必須（空文字は保存不可） |
| `price` | `number` | 単価（円）。0 以上の整数を想定（小数は非推奨だが制約なし） |
| `category` | `string` | `CATEGORIES` 配列のいずれか。削除されたカテゴリでも値は残る |

**デフォルト商品**: サーバーから `null` が返った場合（初回アクセス）、`index.html` L160–172 で定義された 8件のサンプル商品が初期値として使われる。サーバーに保存されるのは `loadState` 後に `save()` が呼ばれたときのみ。

---

### Entry（記録エントリ）

| フィールド | 型 | 説明 |
|---|---|---|
| `id` | `string` | `uid()` で生成 |
| `productId` | `string \| null` | 紐づく商品の `id`。メモで記録した場合は `null` |
| `name` | `string` | 購入時の商品名スナップショット（商品を後から編集しても履歴は変わらない） |
| `emoji` | `string` | 購入時の絵文字スナップショット。メモで記録の場合は `"🛍️"` |
| `price` | `number` | 購入時の単価スナップショット |
| `qty` | `number` | 数量。1タップ記録は常に `1`。長押しメニューで 1〜任意 |
| `category` | `string` | 購入時のカテゴリスナップショット |
| `date` | `string` | `"YYYY-MM-DD"` 形式（`today()` 関数で生成）|
| `time` | `number` | `Date.now()` によるミリ秒タイムスタンプ（履歴のソート・期間フィルタに使用）|

---

### budgets（予算）

```json
{
  "食品": 30000,
  "日用品": 5000,
  "交通": 10000,
  "娯楽": 8000,
  "その他": 5000
}
```

- キーはカテゴリ名（文字列）、値は月次上限額（数値、円）。
- `0` または存在しないキーは「予算なし」として扱われ、バーが非表示になる。
- カテゴリ削除時は `delete budgets[cat]`（L758）で削除。
- カテゴリリネーム時は `budgets[newName] = budgets[oldName]; delete budgets[oldName]`（L740–741）で移行。

---

### notified（アラート通知済みフラグ）

```json
{
  "2025-06-食品": [50, 80],
  "2025-06-娯楽": [100],
  "2025-07-交通": [50]
}
```

- キー形式: `"YYYY-MM"` + `"-"` + カテゴリ名（例: `"2025-06-食品"`）。
- 値: そのカテゴリ・その月に既に発火済みの閾値の配列（`50`、`80`、`100` のいずれかの組み合わせ）。
- `triggerAlerts`（L375）が閾値を超えるたびにキーに閾値を push して永続化。
- 予算編集時は `delete notified[ym+'-'+cat]`（L782）でリセットし、再アラートを可能にする。

---

### categories（カテゴリ一覧）

```json
["食品", "日用品", "交通", "娯楽", "その他"]
```

- 順序があり、UI タブや `<select>` の選択肢の並び順に対応する。
- `CATEGORIES` という `let` 変数に格納（L175）。サーバー送信時は `categories` キーとして保存。
- デフォルト値: `['食品','日用品','交通','娯楽','その他']`（L175）。

---

## 状態の初回ロード・マイグレーション

`loadState()`（L220）の挙動:

```
GET /api/state
  → "null" の場合: 何もしない（variables はデフォルト値のまま）
  → JSON オブジェクトの場合:
      products    が Array なら上書き
      entries     が Array なら上書き
      budgets     が object なら上書き
      notified    が object なら上書き
      categories  が非空 Array なら CATEGORIES を上書き
```

**スキーマ移植**: 追加フィールドがサーバー側 JSON にない場合はデフォルト値が使われる（部分マッチ方式）。iOS 実装でも同様のマージ戦略を推奨。

---

## ID 生成戦略

`uid()`（L233）:

```js
function uid() {
  return 'x' + Math.random().toString(36).slice(2, 9);
}
```

- 先頭に `'x'` を付けることで CSS クラス名・HTML 属性値として安全に使用できる。
- 7文字の 36進数（0–9, a–z）で約 78 億通り。シングルユーザーアプリで重複は実用上無視できる。
- iOS ではより安全な `UUID` を使用することを推奨（Swift の `UUID().uuidString`）。

---

## API エンドポイント

| メソッド | パス | 説明 |
|---|---|---|
| `GET` | `/api/state` | DB から JSON を返す。行なしなら `"null"` を返す。`Cache-Control: no-store` |
| `PUT` | `/api/state` | リクエストボディの JSON を `app_state` に UPSERT。空ボディは 400 エラー |

---

## データフロー概略

```
[ユーザー操作]
    ↓
[JS 状態変数更新]
    ↓
save() → _savePending=true
    ↓ (300ms デバウンス)
flushSave() → PUT /api/state
    ↓
SQLite app_state.data 更新

[ページロード時]
GET /api/state → loadState() → JS 状態変数をサーバー値で初期化
```
