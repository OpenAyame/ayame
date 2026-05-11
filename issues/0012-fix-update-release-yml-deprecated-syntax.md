# 0012-fix-update-release-yml-deprecated-syntax

- Created: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-release-yml-set-output

## 概要

`.github/workflows/release.yml` で deprecated な `::set-output` 構文を使用している。

## 問題

`release.yml:19`:

```yaml
echo ::set-output name=version::$VERSION
```

`::set-output` は GitHub Actions で 2023 年 5 月 31 日以降非推奨となっている。現在の推奨構文は `$GITHUB_OUTPUT` 環境変数への追記である。

## 影響

将来の GitHub Actions ランナーのアップデートでワークフローが動作しなくなる可能性がある。

## 対応方針

`$GITHUB_OUTPUT` を使う構文に更新する:

修正後:
```yaml
echo "version=$VERSION" >> $GITHUB_OUTPUT
```

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [FIX] release.yml の deprecated な ::set-output 構文を $GITHUB_OUTPUT に更新する
  - @voluntas
```

## 後方互換

CI 設定の変更のみ。出力変数の値は変わらない。

## 関連ファイル

- `.github/workflows/release.yml:19`
