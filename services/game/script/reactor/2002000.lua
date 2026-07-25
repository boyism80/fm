-- Reactor name (Reactor.wz/2002000.img.xml): 오르비스용 리엑터.. 빈병이 수거된다.

return {
	on_reactor = function(reactor)
		reactor:drop_items()
	end
}
