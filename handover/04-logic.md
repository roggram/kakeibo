# 04. コアロジックと対応ソース関数

> この章の根拠となるソースコード: `srv/static/index.html`（全関数）。
> 各セクションに対応する関数名と行番号目安を示します。
> 仕様の細部は **`index.html` を直接読んで実装してください**。

---

## 1. 当月支出集計（カテゴリ別 / 全体）

| 項目 | 内容 |
|---|---|
| **目的** | カテゴリ別または全体の当月支出を計算する |
| **入力** | `entries[]`, カテゴリ名（省略可）, 年月 `"YYYY-MM"`（省略時は当月） |
| **出力** | 数値（円） |
| **対応関数** | `categorySpent(cat, ym)`（`index.html` L364）／全体は `entries.filter(e=>e.date.startsWith(ym)).reduce(...)` のインライン式（L254, L262, L621 等） |
| **アルゴリズム** | `entries` をフィルタ（`e.date.startsWith(ym)` + カテゴリ一致）→ `reduce((s,e) => s + e.price * e.qty, 0)` |

```js
function categorySpent(cat, ym) {
  ym = ym || today().slice(0, 7);
  return entries
    .filter(e => e.category === cat && e.date.startsWith(ym))
    .reduce((s, e) => s + e.price * e.qty, 0);
}
```

---

## 2. 予算進捗計算と色判定の閾値

| 項目 | 内容 |
|---|---|
| **目的** | 進捗バーの幅・色を計算する |
| **入力** | `spent`, `limit` |
| **出力** | `pct`（0〜100 にクリップ）、`cls`（`''`, `'mid'`, `'warn'`, `'danger'`）|
| **対応関数** | `budgetBarHTML(opts)`（L259）内のインラインロジック、および `renderBudget()`（L616）内の `barColor`/`cls` 判定 |
| **閾値（色マッピング）** | <br>・`pct < 50%` → 青 `#3b82f6`（クラス無印）<br>・`pct >= 50%` → 黄 `#f59e0b`（クラス `mid`）<br>・`pct >= 80%` → 橙 `#f97316`（クラス `warn`）<br>・`spent >= limit` → 赤 `#ef4444`（クラス `danger` / `over`）|

実装上の注意:
- `pct` は `Math.min(100, spent/limit*100)` で 100 にクリップして CSS 幅に使う。
- 色判定は **クリップしていない実 pct** で行う。`spent >= limit` のチェックは `pct` ではなく金額で直接比較している（L269）。

---

## 3. 50/80/100% アラート発火条件と `notified` への書き込み

| 項目 | 内容 |
|---|---|
| **目的** | 予算消化の節目で「一回だけ」通知する |
| **入力** | カテゴリ `cat`, 記録前合計 `before`, 記録後合計 `after`, 予算 `limit` |
| **対応関数** | `triggerAlerts(cat, before, after, limit)`（L375） |
| **アルゴリズム** | 閾値 `[50, 80, 100]` をループし、`tv = limit * t / 100` の境界を **またいだ** 場合（`before < tv && after >= tv`）かつ `notified[ym+'-'+cat]` に当該閾値が **記録されていない** 場合のみ `alert` を 200ms 遅延で表示。発火した閾値は配列に push して `save()` |

メッセージ文言:
- 100%: `⚠️ ${cat}が予算上限に達しました！`
- それ以外: `🔔 ${cat}が予算の${t}%に到達 (${yen(after)}/${yen(limit)})`

---

## 4. 100% 超購入時の確認ダイアログ

| 項目 | 内容 |
|---|---|
| **目的** | 予算超過する記録は明示的な同意なしには発生させない |
| **対応関数** | `checkBudget(p, addAmount)`（L365）|
| **アルゴリズム** | <br>1. `limit = budgets[p.category]`、未設定（0）なら `{allow:true}` を返して終了。<br>2. `before = categorySpent(cat, ym)`、`after = before + addAmount`。<br>3. `before >= limit` → `{allow:false, reason:'もうこれ以上買えません', cat, limit, spent:before}`（呼び出し元が `alert`「⛔ ...」表示・記録拒否）。<br>4. `after > limit` → `{allow:'confirm', reason:'予算を…超過します。記録しますか？', cat, limit, spent:before, after}`（呼び出し元が `confirm`、キャンセルで記録中止）。<br>5. それ以外 → `{allow:true, before, after, limit}` |

メモで記録（`openQuickEntry`、L462）でも同じロジックを `pseudo = {category, price}` で呼ぶ。

---

## 5. 商品の使用頻度ソート

| 項目 | 内容 |
|---|---|
| **目的** | よく買う商品をタイル先頭に並べる |
| **対応関数** | `usageCount(pid)`（L238）, `tilesHTML()`（L311） |
| **アルゴリズム** | `entries.filter(e => e.productId === pid).length` をその商品の使用回数として算出し、降順ソート。同数の場合は `Array.prototype.sort` の安定性に依存（元順序維持）|

```js
list = [...list].sort((a, b) => usageCount(b.id) - usageCount(a.id));
```

`productId === null` のメモ記録は集計対象外（productId が一致しないため自然に除外される）。

---

## 6. 履歴のグループ化（期間別に月/日）

| 項目 | 内容 |
|---|---|
| **目的** | 月→日のヒエラルキーで履歴を見せる |
| **対応関数** | `renderHistory()`（L537） |
| **アルゴリズム** | <br>1. 期間フィルタ: `cutoffs[historyPeriod]` を `now` から逆算（`month`=1ヶ月前, `3m`=3ヶ月前, `year`=1年前, `all`=0）。`e.time >= cutoff` で絞り込み。<br>2. `time` 降順ソート。<br>3. 月ごとに `byMonth[e.date.slice(0,7)] = []` でグループ化。<br>4. 月内をさらに日ごとに `byDate[e.date] = []` でグループ化。<br>5. 月合計・日合計をそれぞれ `reduce` で算出して表示 |

期間ラベル: `{month:'1ヶ月', '3m':'3ヶ月', year:'1年', all:'全期間'}`（L555）

---

## 7. 集計タブのトップ8計算

| 項目 | 内容 |
|---|---|
| **目的** | 当月の支出が大きい商品ベスト 8 を可視化 |
| **対応関数** | `renderStats()`（L591） |
| **アルゴリズム** | <br>1. 当月エントリのみ抽出。<br>2. `byProd[e.emoji + ' ' + e.name] += e.price * e.qty` で集計（**絵文字と名前の組み合わせがキー**）。<br>3. `Object.entries(byProd).sort((a,b) => b[1] - a[1]).slice(0, 8)`。<br>4. 最大値で正規化して横バー幅を算出 |

注意: `productId` ではなく `"emoji name"` を集計キーにしているので、商品を編集して絵文字や名前を変えると過去エントリと別商品扱いになる。これは仕様。

カテゴリ別バーも同様: `byCat[e.category] += e.price * e.qty`（L602）→ 金額降順ソート → 最大値で正規化。

---

## 8. サーバー同期（debounced PUT・beforeunload の sendBeacon）

| 項目 | 内容 |
|---|---|
| **目的** | 全状態を確実にサーバーへ永続化する |
| **対応関数** | `save()`（L201）, `flushSave()`（L206）, `loadState()`（L220）, `beforeunload` リスナ（L805）|
| **アルゴリズム** | <br>**save()**: `_savePending = true`。タイマー未セットなら 300ms 後に `flushSave` を予約。<br>**flushSave()**: `_saveInflight` が true なら何もしない。フラグを立てて `PUT /api/state` を fetch。失敗時は `_savePending = true` に戻して 500ms 後にリトライ。成功時は `_savePending` が再度立っていればもう一度フラッシュ予約。<br>**beforeunload**: `_savePending` が true の場合のみ `navigator.sendBeacon('/api/state', JSON)` で同期送信 |

送信ペイロード:
```js
JSON.stringify({ products, entries, budgets, notified, categories: CATEGORIES })
```

---

## 9. カテゴリ削除時の商品移行ロジック

| 項目 | 内容 |
|---|---|
| **目的** | カテゴリ削除時に紐づく商品の扱いをユーザーに選ばせる |
| **対応関数** | `openCategoryManager()`（L674）の保存ハンドラ（L726–752） |
| **アルゴリズム** | <br>1. 削除カテゴリのリスト `delList` を作る。<br>2. `delList` に紐づく商品 `affected` を集計。<br>3. `affected.length > 0` の場合 `confirm(`削除するカテゴリに ${count}件 の商品が紐づいています。商品も一緒に削除しますか？\n\nキャンセル=商品は「${fallback}」に移動`)` を表示。<br>4. **OK = 商品も削除**: `products = products.filter(p => !delList.includes(p.category))`<br>5. **キャンセル = 移動**: `fallback = newCats[0] \|\| 'その他'`。`products.forEach(p => { if (delList.includes(p.category)) p.category = fallback; })`<br>6. リネームは `products`, `entries`, `budgets` の全てに波及（L740–741）。<br>7. 削除カテゴリの予算は `delete budgets[cat]`（L744） |

注意: 履歴 `entries` の `category` は **削除時には変更されない**（移動 or 商品削除を選んだ場合のみ。リネームは波及）。これは履歴の歴史性を保つため。

---

## 10. ペース分析（予算タブ：1日あたり、月末予測）

| 項目 | 内容 |
|---|---|
| **目的** | このペースで使い続けたら月末いくらか・1日あたりいくらまで使えるかを表示 |
| **対応関数** | `renderBudget()`（L616, 特に L624–648） |
| **アルゴリズム** | <br>- `daysInMonth = new Date(year, month+1, 0).getDate()`（その月の日数）<br>- `day = 今日の日付`<br>- `dayPct = day / daysInMonth * 100`（月の経過率）<br>- `dailyAllow = round(remain / (daysInMonth - day + 1))`（残額 ÷ 残日数）<br>- `projection = round(totalSpent / day * daysInMonth)`（このペースでの月末予測） |

**全体サマリーの判定**（L633–637）:
- `totalPct > dayPct + 15` → ⚠️ ペースが速い（projection を表示）
- `totalPct < dayPct - 15` → ✨ 順調（節約額 = `totalBudget - projection`）
- それ以外 → 👌 計画通り
- 全パターン共通で `1日あたり ¥XXX まで使えます` を併記

**カテゴリカードの判定**（L654–663）:
- `pct >= 100` → ⛔ 上限到達（超過額を表示）
- `pct >= 80` → ⚠️ 80%超（残額を表示）
- `pct >= 50` → 🔔 半分使用
- 加えて `pct > dayPct + 20 && pct < 100` → ⚠️ ペース速すぎ（projection 表示）
- 加えて `pct < dayPct - 20 && spent > 0` → ✨ 順調

これらは表示のみで、記録の許可・拒否は `checkBudget` 側が担当する。
