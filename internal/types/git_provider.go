package types

type IGitProvider interface {
    CreatePR() error
}
