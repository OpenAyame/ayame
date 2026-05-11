# 0031-refactor-rename-register-result-constants

- Created: 2026-05-11
- Priority: Low
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/refactor-rename-register-result-constants

## 優先度

Low。定数名 `one`/`two`/`full` は数値そのものを表現しており、登録結果のセマンティクス（部屋作成/相手待機中/満員）が伝わらない。可読性の問題であり機能影響はない。広範囲のリネームが必要だが検索置換で安全に行える。他のリファクタリングと同時に対応すればよい。

## 概要

登録結果定数 `one` / `two` / `full` の名前が意図を伝えず、可読性が低い。

## 問題

`messages.go:3-14`:

```go
const (
    one int = iota  // 部屋が作成された
    two             // すでに部屋はあって相手が待ってる
    full            // 満員だった
)
```

`one` / `two` / `full` という名前では登録結果としての意味が伝わらない。

使用箇所:
- `room.go:35,38,49`：`rch <- two` / `rch <- full` / `rch <- one`
- `connection.go:374,383`

## 対応方針

セマンティクスを反映した名前にリネームする:

```go
const (
    registerResultCreated int = iota  // 部屋が作成された
    registerResultPaired              // 部屋に相手が待っている
    registerResultFull                // 部屋が満員
)
```

使用箇所をすべてリネーム:
- `one` → `registerResultCreated`
- `two` → `registerResultPaired`
- `full` → `registerResultFull`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] 登録結果定数を one/two/full から registerResultCreated/Paired/Full にリネームする
  - @voluntas
```

## 関連ファイル

- `messages.go:3-14`（定数定義）
- `room.go:35,38,49`（使用箇所）
- `connection.go:374,383`（使用箇所）
