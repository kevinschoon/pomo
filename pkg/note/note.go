package note

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/kevinschoon/pomo/pkg/internal/random"
)

// Note is a string that can be modified
// in a configured editor.
type Note string

// Edit creates a temporary file populated with
// the content of the Note. Once the file is closed
// it is read again and saves as a string
func (n *Note) Edit(editor string) error {
	path := fmt.Sprintf("/tmp/%s-pomo.tmp", random.NewString(10))
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(string(*n))
	_, err = io.Copy(f, buf)
	if err != nil {
		return err
	}
	f.Close()

	cmd := exec.Command(editor, path)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	f, err = os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	buf = bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f)
	if err != nil {
		return err
	}
	*n = Note(buf.String())
	return nil
}

func (n *Note) MarshalText() ([]byte, error) {
	return []byte(*n), nil
}

func (n *Note) UnmarshalText(raw []byte) error {
	*n = Note(raw)
	return nil
}
