package main

type Client interface {
	Status() (*Status, error)
	Suspend() bool
	Toggle()
	Stop()
}
