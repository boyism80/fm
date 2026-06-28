function on_reactor_2618002(reactor, item)
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
	local reactor_name = "jnr31_out"
	if wz.id == 926100200 then
		reactor_name = "rnj31_out"
	end
	run_on_map(wz.id + 1, "script/reactor/2618001.lua", "rnj32_out_hit", reactor_name)
end
