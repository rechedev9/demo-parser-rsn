package mocking

import (
	demoinfocs "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs"
	events "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/events"
)

func collectKills(parser demoinfocs.Parser) (kills []events.Kill, err error) {
	parser.RegisterEventHandler(func(kill events.Kill) {
		kills = append(kills, kill)
	})
	err = parser.ParseToEnd()
	return
}
