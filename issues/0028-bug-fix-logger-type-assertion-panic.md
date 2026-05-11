# 0028-bug-fix-logger-type-assertion-panic

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: High

## 概要

`logger.go` の型アサーションが安全でない。`i.(string)` が `string` 以外の型で呼ばれた場合に panic する。

## 問題

`logger.go:54,99,116` で安全でない型アサーションを使用している:

```go
// logger.go:54
FormatTimestamp: func(i interface{}) string {
    return strings.Join([]string{darkGray, i.(string), reset}, "")
},

// logger.go:99
switch i.(string) {

// logger.go:116
return fmt.Sprintf("[%s]", filepath.Base(i.(string)))
```

zerolog の内部実装が変更された場合や、予期しない型の値が `FormatTimestamp` / `FormatFieldValue` / `FormatCaller` に渡された場合、サーバーが panic して停止する。

## 影響

サーバー全体が panic で停止する致命的なリスク。

## 対応方針

すべての型アサーションを安全な形式 `s, ok := i.(string)` に変更する:

```go
// logger.go:54
FormatTimestamp: func(i interface{}) string {
    s, ok := i.(string)
    if !ok {
        return ""
    }
    return strings.Join([]string{darkGray, s, reset}, "")
},

// logger.go:99
s, ok := i.(string)
if !ok {
    return ""
}
switch s {

// logger.go:116
s, ok := i.(string)
if !ok {
    return ""
}
return fmt.Sprintf("[%s]", filepath.Base(s))
```

## 作業ブランチ

`feature/fix-logger-type-assertion-panic`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [FIX] logger.go の安全でない型アサーションを修正する
  - @voluntas
```

## 後方互換

型アサーション失敗時の挙動が変わる（panic → 空文字列または早期リターン）。ログ出力の欠落の可能性はあるが、panic よりは望ましい。

## 関連ファイル

- `logger.go:54,99,116`
