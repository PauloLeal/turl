package turl

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type Replacer func() string

var replacements = map[string]Replacer{
	"{{UUID}}": func() string { return uuid.New().String() },
	"{{XULID}}": func() string {
		hexValue := fmt.Sprintf("%x", string(ulid.Make().Bytes()))
		hexBytes, _ := hex.DecodeString(hexValue)
		u, _ := uuid.FromBytes(hexBytes)

		return u.String()
	},
	"{{ULID}}":      func() string { return ulid.Make().String() },
	"{{RANDINT-3}}": func() string { return fmt.Sprintf("%03d", randN(3)) },
	"{{RANDINT-5}}": func() string { return fmt.Sprintf("%05d", randN(5)) },
	"{{SEQINT}}":    func() string { return fmt.Sprintf("%d", nextSeq()) },
	"{{SEQINT-5}}":  func() string { return fmt.Sprintf("%05d", nextSeq()) },
}

var seq int = 0

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randN(n int) int {
	d := math.Pow(10, float64(n))
	return rand.Intn(int(d))
}

func nextSeq() int {
	seq += 1
	return seq - 1
}

func makeReplacements(s string) string {
	for k, v := range replacements {
		for {
			oldS := s
			s = strings.Replace(s, k, v(), 1)

			if oldS == s {
				break
			}
		}
	}

	return s
}

func SetReplacement(mask string, f Replacer) {
	replacements[mask] = f
}

func DelReplacement(mask string) {
	delete(replacements, mask)
}
