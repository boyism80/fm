-- NPC name (String.wz/Npc.img.xml): 휴지통

local LETTER = 4031839

return {
	on_click = function(me, npc)
		if not me:dialog(npc, "쓰다가 만 편지가 보인다.. 중요해 보이는 것 같다. 종이를 주워본다.", false, true) then
			return
		end

		local count = 0
		for _, it in pairs(me:item(LETTER)) do
			count = count + it:count()
		end
		if count >= 1 then
			me:dialog(npc, "이미 편지를 주웠다. 또 다시 주울 필요는 없다.")
			return
		end

		local code = me:exchange(nil, { item = { [LETTER] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족합니다.")
			return
		end
		me:dialog(npc, "뭐라 적혀있는지는 모르겠다..")
	end
}
