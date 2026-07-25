return {
	on_enter = function(me)
		local slots = me:item(4031582)
		local has = false
		for _, _ in pairs(slots) do
			has = true
			break
		end
		if not has then
			me:notice("이곳에 출입하려면 궁전 출입 자격증이 필요합니다.", Msg.PinkText)
			return
		end
		me:play_portal_sound()
		me:map(260000301, 5)
	end
}
