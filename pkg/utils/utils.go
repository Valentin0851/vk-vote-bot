package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type IDGenerator interface {
	Generate() string
}

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) Generate() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

type TimePrefixGenerator struct {
	prefix string
}

func NewTimePrefixGenerator(prefix string) *TimePrefixGenerator {
	return &TimePrefixGenerator{prefix: prefix}
}

func (g *TimePrefixGenerator) Generate() string {
	timestamp := time.Now().Unix()
	randomPart, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("%s_%d_%04d", g.prefix, timestamp, randomPart)
}

type NanoIDGenerator struct {
	alphabet string
	length   int
}

func NewNanoIDGenerator() *NanoIDGenerator {
	return &NanoIDGenerator{
		alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz",
		length:   12,
	}
}

func (g *NanoIDGenerator) Generate() string {
	b := make([]byte, g.length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(g.alphabet))))
		b[i] = g.alphabet[n.Int64()]
	}
	return string(b)
}

func NewIDGenerator(generatorType string) IDGenerator {
	switch generatorType {
	case "uuid":
		return NewUUIDGenerator()
	case "time":
		return NewTimePrefixGenerator("poll")
	default: // "nanoid"
		return NewNanoIDGenerator()
	}
}
