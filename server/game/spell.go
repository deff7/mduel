package game

import (
	//"log"
	"math/rand"
	"time"
)

type Spell struct {
	Distance   int
	BoltSpeed  int
	NextWord   string
	powerLevel int
	cast       bool
}

const initialDistance = 1000

func checkPrefix(s, pre string) bool {
	if len(pre) > len(s) {
		return true
	}
	if s[:len(pre)] == pre {
		return true
	}
	return false
}

func needToEmpower(s string) bool {
	return s[len(s)-1] == ' '
}

func (s *Spell) Empowering() bool {
	return s.powerLevel > 0
}

func (s *Spell) Cast() bool {
	return s.cast
}

// Check how correctly player pronounced the spell
func (s *Spell) Check(msg string) {
	s.cast = false
	if s.Empowering() {
		if checkPrefix(msg, s.NextWord) {
			if needToEmpower(msg) {
				s.Generate()
				s.powerLevel += 1
			} else {
				s.cast = true
			}
		}
	} else {
		if checkPrefix(msg, "hit") {
			if needToEmpower(msg) {
				s.Generate()
				s.powerLevel += 1
			} else {
				s.cast = true
			}
		}
	}
}

func (s *Spell) Generate() {
	rand.Seed(time.Now().Unix())
	words := []string{
		"foo",
		"bar",
		"yoo",
	}
	s.NextWord = words[rand.Int()%len(words)]
}
