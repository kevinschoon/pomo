package main

import (
	"io/ioutil"
	"os"

	"github.com/0xAX/notificator"
)

// Notifier sends a system notification
type Notifier interface {
	Notify(string, string) error
}

// NoopNotifier does nothing
type NoopNotifier struct{}

func (n NoopNotifier) Notify(string, string) error { return nil }

// Xnotifier can push notifications to mac, linux and windows.
type Xnotifier struct {
	*notificator.Notificator
	iconPath string
}

func NewXnotifier(iconPath string) Notifier {
	// Write the built-in tomato icon if it
	// doesn't already exist.
	_, err := os.Stat(iconPath)
	if os.IsNotExist(err) {
		raw := MustAsset("tomato-icon.png")
		_ = ioutil.WriteFile(iconPath, raw, 0644)
	}
	return Xnotifier{
		Notificator: notificator.New(notificator.Options{}),
		iconPath:    iconPath,
	}
}

// Notify sends a notification to the OS.
func (n Xnotifier) Notify(title, body string) error {
	return n.Push(title, body, n.iconPath, notificator.UR_NORMAL)
}
