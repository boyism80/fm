-- Reactor name (Reactor.wz/2519003.img.xml): 데비존의 문

local pirate_pq = require("script/lib/pirate_party_quest")

return {
	on_reactor = function(reactor)
		pirate_pq.lock_door(reactor, 9300126)
	end
}
