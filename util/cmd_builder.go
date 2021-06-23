package util

import (
	"bytes"
	"os"
	"os/exec"
	"text/template"
)

func BuildCommand(tpl string, dir string, date string, file string, pattern string) (*exec.Cmd, error) {
	t := template.Must(template.New("").Parse(tpl))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, map[string]interface{}{
		"BaseDir": dir,
		"Date": date,
		"File": file,
		"Pattern": pattern,
	})
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("sh", "-c", buf.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}
