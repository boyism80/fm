-- NPC name (String.wz/Npc.img.xml): 프로드

local LETTER = 4031035

return {
	on_click = function(me, npc)
		if next(me:item(LETTER)) == nil then
			me:dialog(npc, "펫을 산책시키러 오신건가요? #b#t4031035##k가 있다면 펫과 좀 더 친해질 수 있을텐데..~")
			return
		end
		if not me:dialog(npc, "오호~ #b#t4031035##k를 가져오셨군요! 좋아요! 펫과 좀 더 친해질 수 있도록 마법을 드리겠어요! 이야압~", false, true) then
			return
		end
		if me:pet_id() == 0 then
			me:dialog(npc, "어라라? 펫은 어디에 있나요? 설마.. 혼자 올라오셨다거나...", true, false)
			return
		end
		if me:exchange({ item = { [LETTER] = 1 } }, nil) ~= ExchangeResult.OK then
			return
		end
		me:add_pet_closeness(2)
		me:dialog(npc, "펫과 좀 더 친해지신게 느껴지시나요? 열심히 펫과 친밀도를 올려보세요~", true, false)
	end
}
