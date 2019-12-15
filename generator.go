package main

import (
	"math/rand"
	"strings"

	"github.com/razzie/ideagen/internal"
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
	for _, name := range internal.AssetNames() {
		bytes, _ := internal.Asset(name)
		data[name] = strings.Split(string(bytes), "\n")
	}

	return &Generator{
		Data: data,
	}
}

// Generate ...
func (gen *Generator) Generate() string {
	gen.reset()
	template := gen.pickRandom("template")
	result := gen.fillInTemplate(template)
	return formatOutput(result)
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
	// @ symbol represents a generator. So '@character@', for example,
	// should be replaced with a call to the generateCharacter function
	if strings.Contains(template, "@") {
		command := getTextBetweenTags(template, "@", "@")
		generator := strings.SplitN(command, ":", 2)[0]
		var replacement string
		var parameters []string
		if strings.Contains(command, ":") {
			parameters = strings.Split(strings.Split(command, ":")[1], ",")
		}
		switch generator {
		case "character":
			replacement = gen.generateCharacter(parameters)
		case "goal":
			replacement = gen.generateGoal()
		case "genre":
			replacement = gen.generateGenre(parameters)
		case "wildcard":
			replacement = gen.generateWildcard(parameters)
		case "mood":
			replacement = gen.pickRandomOrNone(generator, 0.3)
		case "setting_description":
			replacement = gen.pickRandomOrNone(generator, 0.7)
		default:
			// theme, setting, character_description
			replacement = gen.pickRandom(generator)
		}

		template = replaceTextBetweenTags(template, replacement, "@", "@")
		// recursively fill in all generators
		return gen.fillInTemplate(template)
	}

	// replace <a> with appropriate indefinite article based on context
	if strings.Contains(template, "<") {
		firstWord := template[strings.Index(template, ">")+2:]
		replacement := indefiniteArticle(firstWord)
		template = replaceTextBetweenTags(template, replacement, "<", ">")
		// recursively fill in all generators
		return gen.fillInTemplate(template)
	}

	// pick conjugation of verb based on character being singular or multiple. E.g. (is,are)
	if strings.Contains(template, "(") {
		optionsList := strings.Split(getTextBetweenTags(template, "(", ")"), ",")
		option := func() string {
			if gen.characterIsGroup {
				return optionsList[1]
			}
			return optionsList[0]
		}()

		template = replaceTextBetweenTags(template, option, "(", ")")
		// recursively fill in all generators
		return gen.fillInTemplate(template)
	}

	return template
}

func (gen *Generator) generateCharacter(parameters []string) string {
	// params:
	allowPostDesc := !contains(parameters, "nopost")
	isPlayer := !contains(parameters, "npc")

	makeGroup := randomChance(0.2)
	preDesc := gen.pickRandomOrNone("character_description", 0.6)

	postDescChance := conditional(allowPostDesc, conditional(len(preDesc) > 0, 0.25, 0.8), 0)
	gen.characterPostDescription = gen.pickRandomOrNone("character_description_post", postDescChance)

	if isPlayer {
		gen.characterIsGroup = makeGroup
	}

	if makeGroup {
		return "<a> " + gen.pickRandom("group_name") + " of " + preDesc + " " +
			pluralize(gen.pickRandom("character")) + " " + gen.characterPostDescription + " "
	}

	character := gen.pickRandom("character")
	return "<a> " + preDesc + " " + character + " " + gen.characterPostDescription
}

func (gen *Generator) generateGoal() string {
	prefix := gen.pickRandom("goal_prefix")

	// Avoid awkward phrasing like: you play as a zombie who is addicted to brains who wants to leave the planet
	// Instead change to: you play as a zombie who is addicted to brains and wants to leave the planet
	if strings.Contains(gen.characterPostDescription, "who") || strings.Contains(gen.characterPostDescription, "that") {
		prefix = strings.ReplaceAll(prefix, "who", "and")
	}

	return prefix + " " + gen.pickRandom("goal")
}

func (gen *Generator) generateGenre(parameters []string) string {
	useModifiers := contains(parameters, "nomods")
	if useModifiers {
		perspective := gen.pickRandomOrNone("perspective", 0.2)
		genreDetail := gen.pickRandomOrNone("genre_modifier", 0.25)
		genre := gen.pickRandomOrNone("genre", conditional(len(perspective) > 0, 0.25, 0.85))
		result := perspective + " " + genreDetail + " " + genre
		return result
	}

	return gen.pickRandom("genre")
}

func (gen *Generator) generateWildcard(parameters []string) string {
	alwaysInclude := contains(parameters, "always")
	wildcard := gen.pickRandomOrNone("wildcard", conditional(alwaysInclude, 1.0, 0.4))
	if len(wildcard) > 0 {
		wildcard = strings.TrimSpace(wildcard) + ","
	}
	return wildcard
}
