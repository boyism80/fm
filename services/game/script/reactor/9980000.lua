-- Reactor name (Reactor.wz/9980000.img.xml): 몬스터 카니발 아티팩트

return {
	on_reactor = function(reactor)
		local name = reactor:name()
		local team_id = tonumber(name:sub(1, 1))
		local guardian = carnival.guardian(tonumber(name:sub(2)))
		local field = reactor:map()
		if guardian ~= nil and field ~= nil then
			for _, mob in pairs(field:mobs()) do
				if mob:carnival_team() == team_id then
					mob:dispel(guardian.mob_skill_id)
				end
			end
		end
		reactor:destroy()
	end
}
