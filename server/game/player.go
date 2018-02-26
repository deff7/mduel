package game

type Player struct {
	ID          int
	HP          int
	Spell       *Spell
	Suggestions map[string]string
	enemy       *Player
}

const initialHP = 100

func newPlayer(id int) *Player {
	p := Player{ID: id, HP: initialHP, Spell: &Spell{}}
	p.updateSuggestions()
	return &p
}

func (p *Player) Hurt(amount int) {
	p.HP -= amount
}

func (p *Player) updateSuggestions() {
	p.Suggestions = generateSpells()
}
