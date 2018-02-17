package game

import (
	//"log"
	"math/rand"
	"time"
)

type Player struct {
	ID         int
	HP         int
	NextWord   string
	powerLevel int
	cast       bool
}

const initialHP = 100

func newPlayer(id int) *Player {
	return &Player{ID: id, HP: initialHP}
}

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

func (p *Player) Empowering() bool {
	return p.powerLevel > 0
}

func (p *Player) Cast() bool {
	return p.cast
}

// Check how correctly player pronounced the spell
func (p *Player) CheckSpelling(msg string) {
	p.cast = false
	if p.Empowering() {
		if checkPrefix(msg, p.NextWord) {
			if needToEmpower(msg) {
				p.GenerateSpell()
				p.powerLevel += 1
			} else {
				p.cast = true
			}
		}
	} else {
		if checkPrefix(msg, "hit") {
			if needToEmpower(msg) {
				p.GenerateSpell()
				p.powerLevel += 1
			} else {
				p.cast = true
			}
		}
	}
}

func (p *Player) GenerateSpell() {
	rand.Seed(time.Now().Unix())
	words := []string{
		"foo",
		"bar",
		"yoo",
	}
	p.NextWord = words[rand.Int()%len(words)]
}
