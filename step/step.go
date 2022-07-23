package step

type Step interface {
	Do() error
	Rollback() error
}
