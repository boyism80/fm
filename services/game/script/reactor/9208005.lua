-- Reactor name (Reactor.wz/9208005.img.xml): 작은 방문_MapC2

return {
	on_reactor = function(reactor, item)
		if item ~= nil then
			return
		end
		reactor:drop_items()
	end
}
