-- NPC name (String.wz/Npc.img.xml): 문월하

return {
	on_click = function(me, npc)
		me:dialog(npc, "홀홀홀", false, false)
	end
}
