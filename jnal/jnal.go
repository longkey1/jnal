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

type Jnal struct {
	cnf Config
}

func NewJnal(cnf Config) Jnal {
	return Jnal{cnf}
}

func (j Jnal) CreateDayFile(day time.Time) string {
	dayFile := j.GetDayFilePath(day)
	if _, err := os.Stat(dayFile); err == nil {
		return dayFile
	}
	if _, err := os.Stat(filepath.Dir(dayFile)); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(dayFile), 0755); err != nil {
			log.Fatalf("Unable to make directory, %v", err)
		}
	}
	file, err := os.OpenFile(dayFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Unable to open file, %v", err)
	}
	_, err = fmt.Fprintln(file, j.buildTargetDayFileContent(day))
	if err != nil {
		log.Fatalf("Unable to build file content, %v", err)
	}
	err = file.Close()
	if err != nil {
		log.Fatalf("Unable to close file, %v", err)
	}

	return dayFile
}

func (j Jnal) BuildOpenCommand(day time.Time) (*exec.Cmd, error) {
	date := day.Format(j.cnf.DateFormat)
	file := j.GetDayFilePath(day)
	return j.buildCommand(j.cnf.OpenCommand, j.cnf.BaseDirectory, date, file)
}

func (j Jnal) BuildListCommand() (*exec.Cmd, error) {
	return j.buildCommand(j.cnf.ListCommand, j.cnf.BaseDirectory, "", "")
}

func (j Jnal) GetBaseDirPath() string {
	return j.cnf.BaseDirectory
}

func (j Jnal) GetDayFilePath(day time.Time) string {
	return fmt.Sprintf("%s/%s", j.cnf.BaseDirectory, day.Format(j.cnf.FileNameFormat))
}

func (j Jnal) buildCommand(tpl string, dir string, date string, file string) (*exec.Cmd, error) {
	t := template.Must(template.New("").Parse(tpl))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, map[string]interface{}{
		"BaseDir": dir,
		"Date":    date,
		"File":    file,
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

func (j Jnal) buildTargetDayFileContent(day time.Time) string {
	t := template.Must(template.New("").Parse(j.cnf.FileTemplate))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, map[string]interface{}{
		"Date": day.Format(j.cnf.DateFormat),
	})
	if err != nil {
		panic(err)
	}
	return buf.String()
}
