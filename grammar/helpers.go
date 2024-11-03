package grammar

var (
	EOFToken *Token
)

func init() {
	EOFToken = NewToken(-1, EtEOF, "")
}
