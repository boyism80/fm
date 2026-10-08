-- NPC name (String.wz/Npc.img.xml): 카오

local PETS = { 5000072, 5000073 }

return {
	on_click = function(me, npc)
		if me:records():get("gift.mystic_pets") > 0 then
			me:dialog(npc, "이미 펫들을 데려갔는걸?")
			return
		end
		local reward = {}
		for _, pet in ipairs(PETS) do
			reward[pet] = 1
		end
		if me:exchange(nil, { item = reward }) ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족한건 아닌지 확인해보게.")
			return
		end
		me:records():set("gift.mystic_pets", 1)
		me:dialog(npc, "당신 에게 신비로운 펫들을 증정하죠.")
	end
}
