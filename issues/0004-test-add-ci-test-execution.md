# 0004-test-add-ci-test-execution

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: High

## 概要

CI パイプラインで `go test` が実行されていない。静的解析（staticcheck）のみが行われており、動的テストによるリグレッション検知の仕組みがない。

`ok_handler_test.go` と `signaling_handler_test.go` にテストが存在するが、CI で実行されないため結果が確認できない。

## 問題

`.github/workflows/ci.yml` にテスト実行ステップが存在しない。`go build` まででパイプラインが終了している。

```yaml
# ci.yml の現在のステップ
- run: go version
- run: go fmt .
- run: go build -v .
```

Makefile には `check` ターゲット（`go test ./...`）と `test` ターゲット（`go test -race -v`）が存在するが、CI から呼び出されていない。また `test` ターゲットには `./...` が欠けており再帰的にテストが走らない。

## 影響

- リグレッション（特に connection、room、webhook のロジック変更によるもの）が PR マージ前に検知されない
- 既存テストが動作する状態を維持する強制力がない

## 対応方針

### CI へのテスト追加

`ci.yml` の `go build -v .` の後にテストステップを追加する:

```yaml
- run: go test -race -v ./...
```

追加後の ci.yml ステップ:

```yaml
- uses: actions/checkout@v6
- name: setup go
  uses: actions/setup-go@v6
  with:
    go-version-file: ./go.mod
    cache: true
    cache-dependency-path: ./go.sum
- uses: dominikh/staticcheck-action@v1
  with:
    version: "2026.1"
    install-go: false
- run: go version
- run: go fmt .
- run: go build -v .
- run: go test -race -v ./...
```

`-race` はデータレース検出に必要。`./...` で全パッケージを再帰的にテストする。`go test` は内部的に `go vet` を実行するため、`go vet ./...` の個別追加は不要。

### Makefile の test ターゲット修正

`Makefile:36` の `go test -race -v` に `./...` が欠けている問題を同時に修正する:

```makefile
test:
	go test -race -v ./...
```

## 作業ブランチ

`feature/add-ci-test`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [ADD] CI に go test -race -v ./... を追加する
  - @voluntas
```

## 後方互換

CI 設定の変更のみであり、アプリケーションコードへの影響はない。

## 関連ファイル

- `.github/workflows/ci.yml`（テストステップ追加）
- `Makefile:36`（test ターゲットの `./...` 追加）
