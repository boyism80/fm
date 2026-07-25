-- NPC name (String.wz/Npc.img.xml): 미로의 멜로디꽃1

return {
	on_click = function(me, npc)
		me:dialog(npc, ".", false, false)
	end
}
