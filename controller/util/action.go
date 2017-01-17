package util

type Action interface {
    PushAck() error
}
