# 0002-bug-fix-ulid-entropy-source

- Created: 2026-05-11
- Priority: Medium
- Model: Qwen 3.6-plus / DeepSeek V4 Pro
- Branch: feature/fix-ulid-entropy-source

## 優先度

Medium。`math/rand` のシード衝突には同一ナノ秒での呼び出しが必要であり、通常の負荷では発生確率は極めて低い。また `ulid.Monotonic` が同一ミリ秒内のエントロピー単調増加を試みるため、完全な ULID 重複にはさらに厳しい条件が必要となる。ただし connectionId の重複はマッチングロジックの破綻に直結するため、暗号的に安全な `crypto/rand` への切り替えはセキュリティ上のベストプラクティスとして実施する。

## 概要

`getULID()` 関数で ULID 生成に予測可能なエントロピー源（ `math/rand` ）を使用している。`math/rand` は `t.UnixNano()` でシード化されており、暗号的に安全ではない。同一ナノ秒に複数コネクションが来た場合に同一シードとなる可能性があり、ULID の一意性が損なわれる。

## 再現手順

1. 高負荷状態で複数クライアントを同時接続する
2. `t.UnixNano()` が同一値となるナノ秒に `getULID()` が複数呼ばれる
3. `rand.NewSource(t.UnixNano())` が同一シードで初期化されるため、`rand.New()` が同一の乱数列を生成する
4. `ulid.Monotonic` が同一エントロピー列を生成し、タイムスタンプ部が同一の場合 ULID 全体が重複する
5. 同一 `connectionId` を持つ複数コネクションが room に登録され、マッチングロジックが破綻する

## 影響

同一 connectionId で複数コネクションが room に登録された場合、シグナリングのマッチングが正しく動作せず、誤った相手との接続やメッセージの誤配送が発生する。

## 対応方針

`math/rand` を `crypto/rand` に置き換える。

`ulid.Monotonic` は維持する。`getULID()` では呼び出しのたびに新しい `MonotonicEntropy` インスタンスが生成されるため、同一ミリ秒内の単調増加は機能しないが、`MustNew()` が要求する `io.Reader` インターフェースに適合させるために必要である。

万が一 `crypto/rand.Reader` からの読み取りに失敗した場合（ファイルディスクリプタ枯渇など）、エントロピー部分がゼロ値になるが、これは現行の `math/rand` ベースの実装でもシード衝突時に同様のリスクがあり、かつ発生確率は極めて低いため許容する。

**注意**: `ulid.MustNew(ms, nil)` とすると entropy がゼロ値になるため、`nil` を渡してはならない。ulid/v2 v2.1.1 の `New()` 実装（ulid.go:106-108）では `entropy` が `nil` の場合、`case nil: return id, err` となりエントロピーが一切埋められない。

### 修正前コード

`connection.go:443-447`:

```go
func getULID() string {
    t := time.Now()
    entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
    return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
```

### 修正後コード

```go
func getULID() string {
    t := time.Now()
    entropy := ulid.Monotonic(cryptoRand.Reader, 0)
    return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
```

### import 変更

`connection.go` の import ブロックに対する変更:

- **削除**: `"math/rand"` （7 行目）- `getULID()` 以外での使用はない
- **追加**: `cryptoRand "crypto/rand"` - `crypto/rand.Reader` 利用のため

`"time"` パッケージは `time.Now()` のために引き続き必要であるため削除しない。

## テスト戦略

### テストファイル

`connection_test.go` を新規作成し、以下を実装する（ `connection_test.go` は既存には存在しない）。

#### TestGetULID

```go
package ayame

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGetULID(t *testing.T) {
    id := getULID()
    assert.Len(t, id, 26)

    ids := make(map[string]bool)
    for i := 0; i < 10000; i++ {
        newID := getULID()
        assert.False(t, ids[newID], "duplicate ULID: %s", newID)
        ids[newID] = true
    }
}
```

#### TestGetULIDConcurrency

```go
func TestGetULIDConcurrency(t *testing.T) {
    ids := make(chan string, 10000)

    // 100 goroutine, each calling getULID() 100 times
    for i := 0; i < 100; i++ {
        go func() {
            for j := 0; j < 100; j++ {
                ids <- getULID()
            }
        }()
    }

    seen := make(map[string]bool)
    for i := 0; i < 10000; i++ {
        id := <-ids
        assert.False(t, seen[id], "duplicate ULID: %s", id)
        seen[id] = true
    }
}
```

## 変更履歴

CHANGES.md の `## develop` セクションに以下を追記する:

```
- [FIX] getULID() のエントロピー源を math/rand から crypto/rand に変更する
  - @voluntas
```

## 後方互換

ULID の文字列表現（26 文字 Crockford's Base32）は変更されないため、後方互換性は完全に維持される。生成される値のエントロピー品質のみが向上する。

## 関連ファイル

- `connection.go:443-447`（getULID 関数）
- `connection.go:7`（`"math/rand"` import 削除対象）
- `connection.go:3-13`（import ブロック）
- `signaling_handler.go:54`（getULID 呼び出し箇所）
