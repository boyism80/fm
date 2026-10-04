-- Reactor name (Reactor.wz/3002000.img.xml): 물 웅덩이

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4001162, min = 1, max = 1, chance = 999999 },
}

return {
	on_reactor = function(reactor)
		reactor:drop_items()
		drop.from_reactor(DROPS, reactor)
	end
}
