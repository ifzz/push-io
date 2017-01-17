package util

type Action interface {
    Ack() error
}
