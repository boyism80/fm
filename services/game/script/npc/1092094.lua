-- NPC name (String.wz/Npc.img.xml): 아기 젖소

return {
	on_click = function(me, npc)
		me:records():remove("cow.milk")
		me:records():remove("cow.last")

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
