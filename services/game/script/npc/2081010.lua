-- NPC name (String.wz/Npc.img.xml): 무스

local SHIELD = 1092041

return {
	on_click = function(me, npc)
		local field = me:map()
		local wz = field ~= nil and field:wz() or nil
		if wz == nil then
			return
		end
		local here = wz:id()
		if here == 924000002 then
			me:map(240010400, 0)
			return
		end
		if here ~= 924000000 then
			me:map(924000002, 0)
			return
		end

		if not me:dialog(npc, "아, 수련장으로 가기 전에 먼저 알려줄게 하나 있어. 수련장에선 반드시 내가 준 #b#t1092041##k 아이템을 착용하고 있어야 해. 그렇지 않으면 죽게 될거야. 조심해~", false, true) then
			return
		end
		local sel = me:dialog_list(npc, "절대 #r방패를 착용하는 것#k을 잊어버리면 안돼! \r\n #b", {
			" #t1092041# 받기.",
			" #m924000001# 가기.",
			" 내보내 주세요.",
		})
		if sel == nil then
			return
		end

		if sel == 2 then
			me:map(924000001, 0)
			return
		end
		if sel == 3 then
			me:map(240010400, 0)
			return
		end

		local count = 0
		for _, it in pairs(me:item(SHIELD)) do
			count = count + it:count()
		end
		if count >= 1 then
			me:dialog(npc, "이미 #t1092041##k를 갖고 있는 것 같은데?")
			return
		end
		local code = me:exchange(nil, { item = { [SHIELD] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "음... 인벤토리 공간이 부족한 건 아닌지 확인해 봐.")
			return
		end
		me:dialog(npc, "음. 그럼 너에게 #t1092041#를 줄게. 인벤토리를 확인해 봐. 반드시 그걸 착용해야 해!")
	end
}
