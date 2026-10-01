-- Reactor name (Reactor.wz/9208002.img.xml): 고블린형 석상

local gq = require("script/lib/guild_quest")

return {
	on_reactor = function(reactor)
		gq.record_statue(reactor)
	end
}
