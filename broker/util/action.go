package util

type Action interface {
    Save() error
    Update() error
    Notify() error
    GetStatus()
}
