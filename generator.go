package main

import (
	"math/rand"
	"strings"
)

// Generator ...
type Generator struct {
	Data                     map[string][]string
	recentlyUsed             []string
	characterIsGroup         bool
	characterPostDescription string
}

// NewGenerator returns a new generator
func NewGenerator() *Generator {
	data := make(map[string][]string)
	for _, name := range AssetNames() {
		bytes, _ := Asset(name)
		data[name] = strings.Split(string(bytes), "\n")
	}

	return &Generator{
		Data: data,
	}
}

// Generate ...
func (gen *Generator) Generate() string {
	gen.reset()
	return gen.fillInTemplate("template")
}

func (gen *Generator) reset() {
	gen.recentlyUsed = nil
	gen.characterIsGroup = false
	gen.characterPostDescription = ""
}

func (gen *Generator) pickRandom(category string) string {
	values, _ := gen.Data[category]
	randomIndex := rand.Int31n(int32(len(values)))

	// Avoid duplicates:
	const maxIterations = 5
	var result string
	for i := 0; i < maxIterations; i++ {
		result = resolveOptions(values[randomIndex])
		if contains(gen.recentlyUsed, result) {
			randomIndex = (randomIndex + 1) % int32(len(values))
		} else {
			gen.recentlyUsed = append(gen.recentlyUsed, result)
			break
		}
	}
	return result
}

func (gen *Generator) pickRandomOrNone(category string, probability float32) string {
	if randomChance(probability) {
		return gen.pickRandom(category)
	}
	return ""
}

func (gen *Generator) fillInTemplate(template string) string {
	return template
}
