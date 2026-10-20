# 0029-enhance-validate-url-scheme

- Created: 2026-05-11
- Priority: High
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/add-validate-url-scheme

## 優先度

Medium。`url.ParseRequestURI` は `file://` 等の不正スキームを許可するが、設定ファイル経由での注入であるため攻撃難易度は高い（設定ファイルへの書き込み権限が必要）。ただし `http`/`https` 以外のスキームが設定された場合の予期しない動作を防ぐ防御的バリデーションとして実施する。設定ファイルのバリデーションはサーバー起動時の 1 回のみであり、パフォーマンス影響はない。

## 概要

webhook URL のバリデーションが `http` / `https` 以外のスキーム（`file://`、`javascript:` 等）を許可している。

## 問題

`config.go:107-117`:

```go
if config.AuthnWebhookURL != "" {
    if _, err := url.ParseRequestURI(config.AuthnWebhookURL); err != nil {
        return nil, err
    }
}
```

`url.ParseRequestURI` は相対 URL を拒否するが、スキームは検証しない。`file:///etc/passwd` や `javascript:alert(1)` のような不正なスキームも有効な URL として受け入れられる。`DisconnectWebhookURL` のバリデーションも同様。

## 影響

不正なスキームの URL が設定ファイル経由で注入された場合、予期しない動作やセキュリティ上の問題を引き起こす可能性がある。

## 対応方針

URL のスキームが `http` または `https` であることを明示的に検証するバリデーション関数を追加する:

```go
func validateWebhookURL(rawURL string) error {
    if rawURL == "" {
        return nil
    }
    u, err := url.Parse(rawURL)
    if err != nil {
        return err
    }
    if u.Scheme != "http" && u.Scheme != "https" {
        return fmt.Errorf("webhook URL must use http or https scheme, got %q", u.Scheme)
    }
    return nil
}
```

各 webhook URL のバリデーションで `url.ParseRequestURI` の代わりにこの関数を使用する。

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [FIX] webhook URL のスキームを http/https のみに制限する
  - @voluntas
```

## 後方互換

`http` / `https` 以外のスキームを設定していた場合、サーバー起動時にエラーとなる。通常の利用では影響なし。

## 関連ファイル

- `config.go:107-117`（AuthnWebhookURL バリデーション）
- `config.go`（DisconnectWebhookURL バリデーション）
