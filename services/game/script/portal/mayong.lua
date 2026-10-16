local ROOM = 240020401

return {

	on_enter = function(me)
		local ok, count = run_on_map(ROOM, "script/lib/party_quest.lua", "character_count")
		if ok and count >= 6 then
			me:message("마뇽의 숲에 입장할 수 없습니다.")
			return
		end
		me:play_portal_sound()
		me:map(ROOM)
	end
}
