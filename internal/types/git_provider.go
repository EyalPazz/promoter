package types

type IGitProvider interface {
    CreatePR(head, body, title string) error
}
