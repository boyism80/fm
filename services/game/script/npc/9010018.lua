-- NPC name (String.wz/Npc.img.xml): 크리샤

return {
	on_click = function(me, npc)
		me:dialog(npc, "메이플 TCG 카드가 출시된 것! 알고 계시나요?", false, false)
	end
}
