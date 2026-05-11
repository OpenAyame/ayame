# 0013-fix-pin-ghr-version-in-release-yml

Created: 2026-05-11
Model: Qwen 3.6-plus
Priority: Medium

## 概要

`.github/workflows/release.yml` で `ghr` を `@latest` でインストールしており、ビルドの再現性がない。

## 問題

`release.yml:21`:

```yaml
- run: go install github.com/tcnksm/ghr@latest
```

`@latest` は実行時の最新バージョンを解決するため、リリースごとに異なるバージョンの `ghr` が使用される可能性がある。`ghr` の破壊的変更が入った場合、リリースワークフローが突然失敗する。

## 影響

リリースプロセスが不安定になり、予期しないタイミングでリリースバイナリの生成に失敗する可能性がある。

## 対応方針

特定のバージョンにピン留めする。2026年5月時点の最新安定版 `v0.16.2` を使用する:

```yaml
- run: go install github.com/tcnksm/ghr@v0.16.2
```

## 作業ブランチ

`feature/fix-pin-ghr-version`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] release.yml の ghr を @latest から @v0.16.2 にピン留めする
  - @voluntas
```

## 後方互換

CI 設定の変更のみ。リリースバイナリ生成の挙動は変わらない。

## 関連ファイル

- `.github/workflows/release.yml:21`
