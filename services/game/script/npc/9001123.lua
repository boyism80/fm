-- NPC name (String.wz/Npc.img.xml): 골드리치

return {
	on_click = function(me, npc)
		local code = me:exchange(nil, { item = { [5030008] = 1, [5470000] = 1 } })
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "인벤토리 공간이 부족합니다.")
			return
		end
		me:dialog(npc, "자유시장 발전에 기여해주시게나.")
	end
}
