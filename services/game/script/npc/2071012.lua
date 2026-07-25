-- NPC name (String.wz/Npc.img.xml): 낯익은 처녀

return {
	on_click = function(me, npc)
		me:dialog(npc, ".", false, false)
	end
}
