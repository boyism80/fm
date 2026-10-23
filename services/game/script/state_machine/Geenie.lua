-- State machine (old/scripts/event/Geenie.js): 오르비스↔아리안트 지니

return require("script/lib/ferry_route").new({
	name = "Geenie",
	legs = {
		{ waiting = 200000152, ride = 200090400, station = 200000151, dest = 260000100 },
		{ waiting = 260000110, ride = 200090410, station = 260000100, dest = 200000100 },
	},
})
