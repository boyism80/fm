local GROUP_NAME = "bishop_resurrection"

return {
	on_enter = function(me)
		local quest = me:quest(6132)
		if not quest:started() then
			me:notice("알 수 없는 힘으로 봉인되어 지나갈 수 없습니다.", Msg.PinkText)
			return
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:notice("지금은 차가운 동굴에 들어갈 수 없습니다.", Msg.PinkText)
			return
		end
		if group:get_property("started") == "true" then
			me:notice("이미 누군가가 퀘스트에 도전 중입니다.", Msg.PinkText)
			return
		end
		local sm, err = group:start_solo(me)
		if sm == nil then
			me:notice("지금은 차가운 동굴에 들어갈 수 없습니다.", Msg.PinkText)
			if err ~= nil then
				log("bishop_resurrection start_solo:", err)
			end
		end
	end
}
