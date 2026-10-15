package stdoutcv

import (
	"bytes"
	"math/rand"
	"runtime"
	"strings"
	"testing"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var letterBytes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func BenchmarkByteSliceAppend(b *testing.B) {
	for b.Loop() {
		var str []byte

		for range 32 {
			num := rand.Intn(len(letterBytes))
			str = append(str, letterBytes[num])
		}

		runtime.KeepAlive(str)
	}
}

func BenchmarkRuneString(b *testing.B) {
	for b.Loop() {
		var str string
		for range 32 {
			num := rand.Intn(len(letterRunes))
			str += string(letterRunes[num])
		}

		runtime.KeepAlive(str)
	}
}

func BenchmarkByteString(b *testing.B) {
	for b.Loop() {
		var str string

		for range 32 {
			num := rand.Intn(len(letterBytes))
			str += string(letterBytes[num])
		}

		runtime.KeepAlive(str)
	}
}

func BenchmarkByteStringAppend(b *testing.B) {
	for b.Loop() {
		var str string
		for range 32 {
			num := rand.Intn(len(letterBytes))

			str = string(append([]byte(str), letterBytes[num]))
		}

		runtime.KeepAlive(str)
	}
}

func BenchmarkRuneSliceAppend(b *testing.B) {
	for b.Loop() {
		var str []rune
		for range 32 {
			num := rand.Intn(len(letterRunes))

			str = append(str, letterRunes[num])
		}
		runtime.KeepAlive(str)
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	for b.Loop() {
		var str strings.Builder
		for range 32 {
			num := rand.Intn(len(letterBytes))

			str.WriteByte(letterBytes[num])
		}
		runtime.KeepAlive(str)
	}
}

func BenchmarkBuffer(b *testing.B) {
	for b.Loop() {
		var str bytes.Buffer
		for range 32 {
			num := rand.Intn(len(letterBytes))

			str.WriteByte(letterBytes[num])
		}
		runtime.KeepAlive(str)
	}
}
