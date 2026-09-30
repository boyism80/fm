-- NPC name (String.wz/Npc.img.xml): 생명의 샘

return {
	on_click = function(me, npc)
		local q = me:quest(6280)
		if q == nil or not q:started() then
			return
		end
		local empty = 0
		for _, it in pairs(me:item(4031454)) do
			empty = empty + it:count()
		end
		if empty < 1 then
			return
		end
		local filled = 0
		for _, it in pairs(me:item(4031455)) do
			filled = filled + it:count()
		end
		if filled >= 1 then
			me:dialog(npc, "이미 #b#t4031455##k를 갖고 있습니다.")
			return
		end
		local code = me:exchange({ item = { [4031454] = 1 } }, { item = { [4031455] = 1 } })
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "인벤토리 공간이 부족하여 성수를 채집할 수 없습니다.")
		end
	end
}
