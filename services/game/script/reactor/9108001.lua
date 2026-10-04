-- Reactor name (Reactor.wz/9108001.img.xml): 달맞이꽃 씨앗

local henesys_pq = require("script/lib/henesys_party_quest")

return {
	on_reactor = function(reactor)
		henesys_pq.plant_seed(reactor)
	end
}
