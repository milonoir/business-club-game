package ui

import (
	"fmt"
	"strings"

	"github.com/milonoir/business-club-game/internal/game"
)

func cardToString(cp *CompanyProvider, c *game.Card) string {
	sb := strings.Builder{}

	for i, mod := range c.Mods {
		company := fmt.Sprintf("[fuchsia]%-12s", "???")
		if mod.Company > -1 {
			name := cp.CompanyByIndex(mod.Company)
			company = fmt.Sprintf("[%s]%-12s", cp.ColorByIndex(mod.Company), name)
		}

		var modifier string
		switch op := mod.Mod.Op(); op {
		case "+":
			modifier = fmt.Sprintf("[green]%s %-3d", op, mod.Mod.Value())
		case "-":
			modifier = fmt.Sprintf("[red]%s %-3d", op, mod.Mod.Value())
		case "*":
			modifier = fmt.Sprintf("[yellow]%s %-3d", op, mod.Mod.Value())
		case "=":
			modifier = fmt.Sprintf("[blue]%s %-3d", op, mod.Mod.Value())
		}

		sb.WriteString(fmt.Sprintf("%s %s", company, modifier))
		if i == 0 {
			sb.WriteString("   ")
		}
	}

	return sb.String()
}
