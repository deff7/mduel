package game

import (
	"log"
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
		return false
	}
	if s[:len(pre)] == pre {
		return true
	}
	return false
}

func needToEmpower(s string) bool {
	if len(s) < 1 {
		return false
	}
	return s[len(s)-1] == ' '
}

func (s *Spell) Empowering() bool {
	return s.powerLevel > 0
}

func (s *Spell) Cast() bool {
	return s.cast
}

func (s *Spell) Discharge() {
	s.NextWord = ""
	s.powerLevel = 0
	s.BoltSpeed = 0
	s.Distance = 0
}

func (s *Spell) castOrEmpower(word, correct string) {
	if checkPrefix(word, correct) {
		if needToEmpower(word) {
			s.Generate()
			s.powerLevel += 1
		} else {
			s.cast = true
		}
	}
}

// Check how correctly player pronounced the spell
func (s *Spell) Check(msg string) {
	s.cast = false
	if s.Empowering() {
		s.castOrEmpower(msg, s.NextWord)
	} else {
		s.castOrEmpower(msg, "hurto")
	}
	log.Printf("from check: %v", s.cast)
}

func (s *Spell) Generate() {
	rand.Seed(time.Now().Unix())
	words := []string{
		"maxima",
		"cruelo",
		"damago",
	}
	s.NextWord = words[rand.Int()%len(words)]
}
