-- NPC name (String.wz/Npc.img.xml): 황금나침반

return {
	on_click = function(me, npc)
		me:dialog(npc, ".", false, false)
	end
}
