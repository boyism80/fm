-- NPC name (String.wz/Npc.img.xml): 루퍼트

return {
	on_click = function(me, npc)
		me:dialog(npc, "3차 전직을 해주는 엔피시", false, false)
	end
}
