-- NPC name (String.wz/Npc.img.xml): 고양이

return {
	on_click = function(me, npc)
		if me:records():get("gift.fishing_cat") > 0 then
			me:dialog(npc, "이미 고양이를 가져갔는걸?")
			return
		end
		local code = me:exchange(nil, { item = { [3010132] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족합니다.")
			return
		end
		me:records():set("gift.fishing_cat", 1)
		me:dialog(npc, "역시 낚시에 가장 필요한건 고양이겠지?")
	end
}
