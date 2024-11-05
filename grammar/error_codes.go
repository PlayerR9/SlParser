package grammar

type FaultCode int

const (
	NotAsExpected FaultCode = iota
)

func (code FaultCode) String() string {
	return [...]string{
		"Not As Expected",
	}[code]
}
