-- Reactor name (Reactor.wz/2502002.img.xml): 바위탁자

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 2022252, min = 1, max = 1, chance = 999999, quest = 3839 },
}

return {
	on_reactor = function(reactor)
		reactor:drop_items()
		drop.from_reactor(DROPS, reactor)
	end
}
