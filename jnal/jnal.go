package jnal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"
)

func CreateFile(config Config, targetDay time.Time) {
	dayFile := BuildTargetDayFileName(config.BaseDirectory, config.FileNameFormat, targetDay)
	dayDir := filepath.Dir(dayFile)
	dayDate := targetDay.Format(config.DateFormat)
	if _, err := os.Stat(dayDir); os.IsNotExist(err) {
		if err = os.MkdirAll(dayDir, 0755); err != nil {
			log.Fatalf("Unable to make directory, %v", err)
		}
	}
	if _, err := os.Stat(dayFile); os.IsNotExist(err) {
		file, err := os.OpenFile(dayFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalf("Unable to open file, %v", err)
		}
		_, err = fmt.Fprintln(file, BuildTargetDayFileContent(config.FileTemplate, dayDate))
		if err != nil {
			log.Fatalf("Unable to build file content, %v", err)
		}
		err = file.Close()
		if err != nil {
			log.Fatalf("Unable to close file, %v", err)
		}
	}
}

func BuildCommand(tpl string, dir string, date string, file string, pattern string) (*exec.Cmd, error) {
	t := template.Must(template.New("").Parse(tpl))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, map[string]interface{}{
		"BaseDir": dir,
		"Date":    date,
		"File":    file,
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

func BuildTargetDayFileContent(tpl string, date string) string {
	t := template.Must(template.New("").Parse(tpl))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, map[string]interface{}{
		"Date": date,
	})
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func BuildTargetDayFileName(dir string, file string, day time.Time) string {
	return fmt.Sprintf("%s/%s", dir, day.Format(file))
}
