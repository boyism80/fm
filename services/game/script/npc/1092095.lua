-- NPC name (String.wz/Npc.img.xml): 아기 젖소

return {
	on_click = function(me, npc)
		local milk_q = me:quest(126640)
		local cow_q = me:quest(126641)
		if not milk_q:started() then
			milk_q:start("0")
		end
		if not cow_q:started() then
			cow_q:start("0")
		end
		milk_q:record("0")
		cow_q:record("0")

		local count = 0
		for _, it in pairs(me:item(4031850)) do
			count = count + it:count()
		end
		if count > 0 then
			me:exchange({ item = { [4031850] = count } }, nil)
		end
		me:dialog(npc, "아기 젖소가 우유통에 든 우유를 모두 먹어버렸습니다.")
	end
}
