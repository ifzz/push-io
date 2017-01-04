package util

type Action interface {
    Notify() error
}
