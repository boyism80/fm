-- Reactor name (Reactor.wz/2112014.img.xml): 큰상자로 열쇠7개를 내려놓으면 불의원석이 나온다.

return {
	on_reactor = function(reactor)
		reactor:drop_items()
	end
}
