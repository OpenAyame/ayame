# 0034-refactor-simplify-setdefaultsconfig

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Low

## 概要

`setDefaultsConfig` で同一パターンの if 文が 15 回繰り返されており、保守性が低い。

## 問題

`config.go:124-189`:

```go
func setDefaultsConfig(config *Config) {
    if config.LogDir == "" {
        config.LogDir = defaultLogDir
    }
    if config.LogName == "" {
        config.LogName = defaultLogName
    }
    if config.LogRotateMaxSize == 0 {
        config.LogRotateMaxSize = defaultLogRotateMaxSize
    }
    // ... 15個の同一パターン
}
```

同一パターンの if 文が連続しており、新しい設定を追加するたびに同じコードを追加する必要がある。

フィールドの型は `string`、`int32`、`[]string` の 3 種が混在する。`string` と `int32` はゼロ値比較（`== ""` / `== 0`）で判定可能だが、`[]string` は `comparable` 制約を満たさないため Go ジェネリクスの直接適用に制約がある。

## 対応方針

### 方法1: 非ゼロ値フィールドのみにジェネリックヘルパーを使用

```go
func setDefaultIfZero[T ~string | ~int32](field *T, defaultVal T) {
    var zero T
    if *field == zero {
        *field = defaultVal
    }
}
```

`SignalingLogFilters`（`[]string`）は 0009 で修正済みの `len(...) == 0` パターンをそのまま使う（`== nil` チェックは `len` に統合済みの想定）。

```go
func setDefaultsConfig(config *Config) {
    setDefaultIfZero(&config.LogDir, defaultLogDir)
    setDefaultIfZero(&config.LogName, defaultLogName)
    setDefaultIfZero(&config.LogRotateMaxSize, defaultLogRotateMaxSize)
    // ... 14項目
    if len(config.SignalingLogFilters) == 0 {
        config.SignalingLogFilters = defaultSignalingLogFilters
    }
}
```

### 方法2: 現状維持

if 文の羅列は冗長だが、コードの意図は明確である。また `gopkg.in/ini.v1` の struct tag に `Default` 的な機能がないため、ライブラリ任せにもできない。

### 推奨

方法1 を推奨する。`~string | ~int32` の型制約により、将来フィールドが `int64` などに変更された場合でもコンパイルエラーで検知可能。

## 作業ブランチ

`feature/refactor-simplify-setdefaultsconfig`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [UPDATE] setDefaultsConfig の重複 if 文をジェネリックヘルパーで簡素化する
  - @voluntas
```

## 関連ファイル

- `config.go:124-189`
