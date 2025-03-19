package envutils

import (
	"encoding"
	"encoding/hex"

	"github.com/pkg/errors"
)

var _ encoding.TextUnmarshaler = (*HexBytes)(nil)

type HexBytes []byte

func (b *HexBytes) UnmarshalText(text []byte) error {
	decoded, err := hex.DecodeString(string(text))
	if err != nil {
		return errors.WithStack(err)
	}

	*b = decoded

	return nil
}
