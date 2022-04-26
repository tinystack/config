package codec

import (
	"bytes"

	"github.com/subosito/gotenv"
)

type EnvCodec struct{}

func (EnvCodec) Decode(b []byte, v map[string]interface{}) error {
	var buf bytes.Buffer

	_, err := buf.Write(b)
	if err != nil {
		return err
	}

	env, err := gotenv.StrictParse(&buf)
	if err != nil {
		return err
	}

	for key, value := range env {
		v[key] = value
	}

	return nil
}
