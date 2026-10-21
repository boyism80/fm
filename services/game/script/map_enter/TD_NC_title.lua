-- Map script (onUserEnter TD_NC_title): 네오시티 입장 제목 효과. 맵 번호의 백의 자리가 장소다

return {
	on_enter = function(me, map)
		local place = math.floor(map:wz():id() / 100) % 10
		if place == 0 then
			me:map_effect("temaD/enter/teraFores")
		elseif place <= 6 then
			me:map_effect("temaD/enter/neoCity" .. place)
		end
	end,
}
