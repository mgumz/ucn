package print

import (
	"io"

	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

type Func func(w io.Writer, entries []uc.Entry)
