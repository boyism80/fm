-- NPC name (String.wz/Npc.img.xml): 나무책상

return {
	on_click = function(me, npc)
		me:dialog(npc, "...", false, false)
	end
}
