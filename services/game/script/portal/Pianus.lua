local ROOM = 230040420

return {

	on_enter = function(me)
		local ok, count = run_on_map(ROOM, "script/lib/party_quest.lua", "character_count")
		if ok and count > 10 then
			me:message("이 방은 이미 피아누스와의 전투를 위한 최대 인원수 만큼 가득 찼습니다.")
			return
		end
		me:play_portal_sound()
		me:map(ROOM)
	end
}
