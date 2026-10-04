-- Reactor name (Reactor.wz/2612005.img.xml): 5색 비커

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031798, min = 1, max = 1, chance = 999999, quest = 3366 },
}

return {
	on_reactor = function(reactor)
		reactor:drop_items()
		drop.from_reactor(DROPS, reactor)
	end
}
