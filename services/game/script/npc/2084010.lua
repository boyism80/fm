-- NPC name (String.wz/Npc.img.xml): 황금열쇠박스

return {
	on_click = function(me, npc)
		me:dialog(npc, ".", false, false)
	end
}
