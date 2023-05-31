package ayame

const (
	// 登録成功
	// 部屋が作成された相手を待つ
	one int = iota
	// すでに部屋はあり相手が待ってる
	// offer を出させる
	// 登録成功
	two
	// 満員だったので Reject か Error を返す
	// 登録失敗
	full
)

type register struct {
	connection    *connection
	resultChannel chan int
}

// rawMessage には JOSN パース前の offer / answer / candidate が入る
type forward struct {
	connection *connection
	rawMessage []byte
}

type unregister struct {
	connection *connection
}
