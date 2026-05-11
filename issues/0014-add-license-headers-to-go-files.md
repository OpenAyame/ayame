# 0014-add-license-headers-to-go-files

- Created: 2026-05-11
- Completed: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/add-license-headers

## 概要

すべての Go ソースファイル（18 ファイル）に Apache 2.0 ライセンスヘッダーが存在しない。

## 問題

`LICENSE` ファイルはリポジトリ直下に存在するが、各ソースファイルにライセンスヘッダーが付与されていない。Apache 2.0 ライセンスの付録（APPENDIX）では、各ソースファイルへのライセンスヘッダー付与が推奨されている。

現状、ファイルの先頭は `package ayame` から始まっている。

## 対応方針

全 18 ファイルの先頭に以下のヘッダーを追加する。`package` 宣言の前に空行を 1 行入れる:

```go
// Copyright 2019-2026, Shiguredo Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ayame
```

## 対象ファイル一覧

```bash
find . -name "*.go" -not -path "./vendor/*" | sort
```

全 18 ファイル:
- `authn_webhook.go`
- `config.go`
- `connection.go`
- `connection_log.go`
- `disconnect_webhook.go`
- `log.go`
- `metrics.go`
- `metrics_keys.go`
- `ok_handler.go`
- `ok_handler_test.go`
- `room.go`
- `server.go`
- `signaling.go`
- `signaling_handler.go`
- `signaling_handler_test.go`
- `types.go`
- `webhook.go`
- `cmd/ayame/main.go`

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
### misc
- [ADD] 全 Go ソースファイルに Apache 2.0 ライセンスヘッダーを追加する
  - @voluntas
```

## 後方互換

コメント追加のみ。コンパイル結果や実行時の挙動に影響はない。

## 関連ファイル

- 全 `.go` ファイル（18 ファイル）
