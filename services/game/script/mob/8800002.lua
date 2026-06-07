function on_mob_die_8800002(mob, attacker, map)
	if map == nil then
		return
	end
	if attacker ~= nil then
		map:message(attacker:name() .. '님이 자쿰을 처치했습니다.')
	else
		map:message('자쿰이 처치되었습니다.')
	end
end
