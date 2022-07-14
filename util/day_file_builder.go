package util

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

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
