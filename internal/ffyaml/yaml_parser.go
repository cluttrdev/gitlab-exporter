package ffyaml

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parser(r io.Reader, set func(name string, value string) error) error {
	var m map[string]interface{}
	if err := yaml.NewDecoder(r).Decode(&m); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("error parsing yaml: %w", err)
	}

	if err := traverseMap("", m, set); err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	return nil
}

const (
	delimiter string = "-"
)

func traverseMap(key string, val any, set func(name string, value string) error) error {
	key = sanitizeKey(key)

	switch v := val.(type) {
	case string:
		return set(key, v)
	case json.Number:
		return set(key, v.String())
	case uint64:
		return set(key, strconv.FormatUint(v, 10))
	case int:
		return set(key, strconv.Itoa(v))
	case int64:
		return set(key, strconv.FormatInt(v, 10))
	case float64:
		return set(key, strconv.FormatFloat(v, 'g', -1, 64))
	case bool:
		return set(key, strconv.FormatBool(v))
	case nil:
		return set(key, "")
	case []any:
		for _, v := range v {
			if err := traverseMap(key, v, set); err != nil {
				return err
			}
		}
	case map[string]any:
		for k, v := range v {
			if key != "" {
				k = key + delimiter + k
			}
			if err := traverseMap(k, v, set); err != nil {
				return err
			}
		}
	case map[any]any:
		for k, v := range v {
			ks := fmt.Sprint(k)
			if key != "" {
				ks = key + delimiter + ks
			}
			if err := traverseMap(ks, v, set); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("couldn't convert %q (type %T) to string", val, val)
	}

	return nil
}

func sanitizeKey(key string) string {
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "_", "-")
	return key
}
