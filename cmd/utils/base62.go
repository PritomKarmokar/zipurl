package utils

import "math/big"

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeString(s string) string {
	n := new(big.Int).SetBytes([]byte(s))

	if n.Sign() == 0 {
		return "0"
	}

	base := big.NewInt(62)
	zero := big.NewInt(0)

	var result []byte

	for n.Cmp(zero) > 0 {
		mod := new(big.Int)
		n.DivMod(n, base, mod)

		result = append(result, alphabet[mod.Int64()])
	}

	// reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

func DecodeString(encoded string) string {
	lookup := make(map[rune]int)

	for i, c := range alphabet {
		lookup[c] = i
	}

	result := big.NewInt(0)
	base := big.NewInt(62)

	for _, c := range encoded {
		result.Mul(result, base)
		result.Add(result, big.NewInt(int64(lookup[c])))
	}

	return string(result.Bytes())
}
