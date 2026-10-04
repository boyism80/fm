-- Reactor name (Reactor.wz/9202012.img.xml): 보물상자&lt;보너스>

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 2290113, min = 1, max = 1, chance = 3600 },
	{ item = 2290118, min = 1, max = 1, chance = 4400 },
	{ item = 2290120, min = 1, max = 1, chance = 5400 },
}

return {
	on_reactor = function(reactor)
		reactor:drop_items()
		drop.from_reactor(DROPS, reactor)
	end
}
