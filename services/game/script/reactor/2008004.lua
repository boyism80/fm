-- Reactor name (Reactor.wz/2008004.img.xml): 여신상 조각

local orbis_pq = require("script/lib/orbis_party_quest")

return {
	on_reactor = function(reactor)
		orbis_pq.piece_hit(reactor)
	end
}
