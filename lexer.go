package main

func matchPrefix(d deterministicFiniteAutomata, s string) string {
	matchingPrefixIndex := -1
	for i, character := range s {
		d.move(transitionLabel(character))
		if d.accepted {
			matchingPrefixIndex = i
		}
		if d.dead {
			break
		}
	}

	if matchingPrefixIndex == -1 {
		return ""
	}
	return s[:matchingPrefixIndex+1]
}
