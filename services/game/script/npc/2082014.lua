-- NPC name (String.wz/Npc.img.xml): 아아시아

return {
	on_click = function(me, npc)
		me:dialog(npc, "대사 모름2", false, false)
	end
}
