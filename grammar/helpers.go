package grammar

var (
	// EOFToken is the end of file token.
	EOFToken *Token
)

func init() {
	EOFToken = NewToken(-1, EtEOF, "")
}
