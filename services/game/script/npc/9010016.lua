-- NPC name (String.wz/Npc.img.xml): 아루

local TRADES = {
	{ cost = 30, reward = 4001168, count = 1, text = "#v4001168# #z4001168# 1개 획득!" },
	{ cost = 1000, reward = 2049100, count = 3, text = "#v2049100# #z2049100# 3장 획득!" },
	{ cost = 20, reward = 1082433, count = 1, text = "#v1082433# #b프리미엄 메이플 이올렛 타투#k 획득!" },
}

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "안녕하세요! #h #님,\r\n혹시 #r단풍잎#k을 가지고 계시지 않으신가요?", {
			"#r#v4001126# #z4001126# 30개 = #v4001168# #z4001168# 1개",
			"#v4001126# #z4001126# 1000개 = #v2049100# #z2049100 3개#",
			"준비중...",
			"준비중...",
			"준비중...",
			"준비중...",
			"준비중...",
		})
		if sel == nil then
			return
		end
		local trade = TRADES[sel]
		if trade == nil then
			return
		end

		local code = me:exchange({ item = { [4001126] = trade.cost } }, { item = { [trade.reward] = trade.count } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "단풍잎이 있는지 확인해주세요")
			return
		end
		me:dialog(npc, trade.text)
	end
}
