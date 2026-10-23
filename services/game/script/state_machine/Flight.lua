-- State machine (old/scripts/event/Flight.js): 오르비스↔리프레 비행선

return require("script/lib/ferry_route").new({
	name = "Flight",
	legs = {
		{ waiting = 200000132, ride = 200090200, station = 200000131, dest = 240000100 },
		{ waiting = 240000111, ride = 200090210, station = 240000110, dest = 200000100 },
	},
})
