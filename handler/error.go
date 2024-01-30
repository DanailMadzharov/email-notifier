package handler

const (
	RECOVERABLE     = "RECOVERABLE"
	NON_RECOVERABLE = "NON_RECOVERABLE"
)

type Error struct {
	ErrorType string
}
