package bitkub

import "strings"

func SwapSymbol(symbol string) string {
	parts := strings.SplitN(symbol, "_", 2)
	return parts[1] + "_" + parts[0]
}
