/*
libnotify is a lightweight client for libnotify https://developer.gnome.org/notification-spec/.
For now this just shells out to "notify-send".
TODO: Move this into it's own repository as time permits.
*/
package libnotify

import (
	"fmt"
	"os/exec"
	"time"
)

type Notification struct {
	Urgency string
	Expire  time.Duration
	Title   string
	Body    string
	Icon    string
}

type Client struct {
	Path string
}

func NewClient() *Client {
	return &Client{
		Path: "/bin/notify-send",
	}
}

func (c Client) Notify(n Notification) error {
	var args []string
	if n.Urgency != "" {
		args = append(args, fmt.Sprintf("--urgency=%s", n.Urgency))
	}
	if n.Icon != "" {
		args = append(args, fmt.Sprintf("--icon=%s", n.Icon))
	}
	if n.Expire > 0 {
		args = append(args, fmt.Sprintf("--expire=%s", n.Expire.Truncate(time.Millisecond)))
	}
	args = append(args, n.Title)
	args = append(args, n.Body)
	_, err := exec.Command(c.Path, args...).Output()
	return err
}
