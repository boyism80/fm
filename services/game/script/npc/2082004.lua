-- NPC name (String.wz/Npc.img.xml): 앤디

return {
	on_click = function(me, npc)
		me:dialog(npc, "대사 모름.", false, false)
	end
}
