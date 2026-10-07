-- NPC name (String.wz/Npc.img.xml): 네르

local LETTER = 4031128

return {
	on_click = function(me, npc)
		if next(me:item(LETTER)) == nil then
			me:dialog(npc, "여기까지 어쩐 일로 올라오신걸까요?")
			return
		end
		if not me:dialog(npc, "이것은 형의편지! 와.. 설마 여기까지 펫과 함께 올라오신거에요? 대단하신데요~", false, true) then
			return
		end
		if me:pet_id() == 0 then
			me:dialog(npc, "어라? 펫은 어디에 있나요? 설마 펫 없이 혼자 올라오셨나요? 그러시면 안되죠~", true, false)
			return
		end
		if me:exchange({ item = { [LETTER] = 1 } }, nil) ~= ExchangeResult.OK then
			return
		end
		me:add_pet_closeness(4)
		me:dialog(npc, "어때요? 펫과 더 친해지신 느낌이 드나요? 헤헤", true, false)
	end
}
