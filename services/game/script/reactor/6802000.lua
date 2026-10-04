-- Reactor name (Reactor.wz/6802000.img.xml): 메이플아일랜드 범용리엑터

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 10, max = 49, chance = 999999 },
	{ meso = true, min = 10, max = 49, chance = 999999 },
	{ meso = true, min = 10, max = 49, chance = 999999 },
	{ meso = true, min = 10, max = 49, chance = 999999 },
	{ item = 2050004, min = 1, max = 1, chance = 250000 },
}

return {
	on_reactor = function(reactor)
		reactor:drop_items()
		drop.from_reactor(DROPS, reactor)
	end
}
