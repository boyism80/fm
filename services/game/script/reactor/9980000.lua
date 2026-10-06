-- Reactor name (Reactor.wz/9980000.img.xml): 몬스터 카니발 아티팩트

local cpq = require("script/lib/carnival")

return {
	on_reactor = function(reactor)
		cpq.destroy_guardian(reactor)
	end
}
