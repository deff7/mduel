package game

type Player struct {
	ID    int
	HP    int
	Spell *Spell
}

const initialHP = 100

func newPlayer(id int) *Player {
	return &Player{ID: id, HP: initialHP, Spell: &Spell{}}
}

func (p *Player) Hurt(amount int) {
	p.HP -= amount
}
