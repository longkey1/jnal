package journal

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// executeTemplate executes a template string with the given data
func (j *Journal) executeTemplate(tpl string, data map[string]interface{}) (string, error) {
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// getEnvMap returns a map of all environment variables
func getEnvMap() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		key, value, found := strings.Cut(env, "=")
		if found {
			envMap[key] = value
		}
	}
	return envMap
}
