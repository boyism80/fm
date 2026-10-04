-- Reactor name (Reactor.wz/3002001.img.xml): 보라색 마력석 상자

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4001163, min = 1, max = 1, chance = 999999 },
}

return {
	on_reactor = function(reactor)
		reactor:drop_items()
		drop.from_reactor(DROPS, reactor)
	end
}
