package types

type IBaseCommand interface {
	Execute(bool, string, string) error
}
