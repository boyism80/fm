-- Reactor name (Reactor.wz/1202002.img.xml): 노틸러스호 동력실 조개

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031843, min = 1, max = 1, chance = 999999, quest = 2169 },
}

return {
	on_reactor = function(reactor)
		reactor:drop_items()
		drop.from_reactor(DROPS, reactor)
	end
}
