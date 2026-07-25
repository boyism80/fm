-- NPC name (String.wz/Npc.img.xml): 박첨지

return {
	on_click = function(me, npc)
		me:dialog(npc, ".", false, false)
	end
}
