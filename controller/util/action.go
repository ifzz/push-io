package util

type Action interface {
    Acknowledge() error
}
