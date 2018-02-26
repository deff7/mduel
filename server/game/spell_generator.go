package game

func generateSpell(t string) string {
	switch t {
	case "attack":
		return t
	case "shield":
		return t
	case "heal":
		return t
	case "stun":
		return t
	}
	return ""
}

func generateSpells() map[string]string {
	res := map[string]string{}
	for _, t := range spellTypes {
		res[t] = generateSpell(t)
	}
	return res
}
