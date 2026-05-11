# 0036-refactor-remove-stale-todo-comments

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Low

## 概要

放置された TODO コメントが複数存在する。不要なものを削除し、必要なものを issue 化する。

## 問題

以下の TODO コメントがコード中に放置されている:

| ファイル | 行 | 内容 | 対応 |
|---------|-----|------|------|
| `logger.go` | 96 | 各色を定数に置き換える | 削除（優先度低、現状で問題なし） |
| `logger.go` | 118 | Caller をファイル名と行番号だけの表示で出力する | 削除（優先度低） |
| `logger.go` | 121 | name=value が無い場合に `\|` を消す方法がわからなかった | 削除（zerolog の制限） |
| `logger.go` | 132 | カンマ区切りをどう実現するかわからなかった | 削除（zerolog の制限、タイポ「同」あり） |
| `authn_webhook.go` | 68 | ヘッダーのサイズも計測する | 削除（TODO としては残す価値なし） |
| `disconnect_webhook.go` | 47 | ヘッダーのサイズも計測する | 削除（同上） |
| `connection.go` | 73 | reason の長さが不十分な場合に TextMessage を使用 | issue 化（バグの可能性） |
| `connection.go` | 418 | standalone で type: connected 受信時のエラー化検討 | 削除（設計判断済み、standalone モードでは切断する） |

## 対応方針

- ほとんどの TODO は実現可能性が低いか低優先度のため削除する
- `connection.go:73` の `reason` 長さ対応のみ、別途 issue を立てる価値がある

## 作業ブランチ

`feature/refactor-remove-stale-todo-comments`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] 放置された TODO コメントを削除する
  - @voluntas
```

## 関連ファイル

- `logger.go:96,118,121,132`
- `authn_webhook.go:68`
- `disconnect_webhook.go:47`
- `connection.go:418`
