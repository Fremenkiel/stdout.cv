package dotenv

import (
	"bytes"
	"io"
	"os"
	"strings"
)

func Parse(r io.Reader) (map[string]string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	return UnmarshalBytes(buf.Bytes())
}

func Load() error {
	var filenames []string
	switch os.Getenv("APP_ENV") {
	case "production":
		return nil
	case "test":
		filenames = []string{".env", ".env.test"}
	case "development":
		filenames = []string{".env", ".env.development"}
	default:
		filenames = []string{".env"}
	}

	for _, filename := range filenames {
		if err := loadFile(filename, true); err != nil {
			return err
		}
	}

	return nil
}

func Unmarshal(str string) (envMap map[string]string, err error) {
	return UnmarshalBytes([]byte(str))
}

func UnmarshalBytes(src []byte) (map[string]string, error) {
	out := make(map[string]string)
	err := parseBytes(src, out)

	return out, err
}

func loadFile(filename string, overload bool) error {
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key, _, _ := strings.Cut(rawEnvLine, "=")
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	return Parse(file)
}
