package print

import "io"

type Func func(w io.Writer, runes []rune)
