-- Reactor name (Reactor.wz/2618001.img.xml): 메이플아일랜드 범용리엑터

function rnj32_out_hit(map, reactor_name)
	local reactor = map:reactor_by_name(reactor_name)
	if reactor == nil then
		return false
	end
	reactor:hit(1)
	sleep(2000)
	if reactor:state() == 1 then
		reactor:hit(0)
	end
	return true
end

function on_reactor_2618001(reactor, item)
	if item ~= nil then
		return
	end
	local map = reactor:map()
	if map == nil then
		return
	end
	local wz = map:wz()
	if wz == nil then
		return
	end
	local reactor_name = "jnr32_out"
	if wz.id == 926100200 then
		reactor_name = "rnj32_out"
	end
	run_on_map(wz.id + 2, "script/reactor/2618001.lua", "rnj32_out_hit", reactor_name)
end
