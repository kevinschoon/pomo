package notify

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/0xAX/notificator"

	"github.com/kevinschoon/pomo/pkg/runner"
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

func NewXNotifier(iconPath string) Notifier {
	// Write the built-in tomato icon if it
	// doesn't already exist.
	_, err := os.Stat(iconPath)
	if os.IsNotExist(err) {
		raw := MustAsset("pkg/notify/tomato-icon.png")
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

func StatusFunc(notifier Notifier) runner.StatusFunc {
	return func(s runner.Status) error {
		if s.Previous != s.State {
			switch s.State {
			// case runner.INITIALIZED:
			// 	return notifier.Notify("pomo", fmt.Sprintf("starting task %s", s.Message))
			case runner.SUSPENDED:
				return notifier.Notify("pomo", fmt.Sprintf(
					"task %s is suspended", s.Message,
				))
			case runner.RUNNING:
				return notifier.Notify("pomo", fmt.Sprintf(
					"starting pomodoro %d/%d (%s)",
					s.Count, s.NPomodoros, s.Message,
				))
			case runner.BREAKING:
				return notifier.Notify("pomo", fmt.Sprintf("it's time to take a break!"))
			case runner.COMPLETE:
				return notifier.Notify("pomo", fmt.Sprintf(
					"task %s is completed!\n",
					s.Message,
				))
			}
		}
		return nil
	}
}
