-- NPC name (String.wz/Npc.img.xml): 조시

return {
	on_click = function(me, npc)
		me:dialog(npc, "히이잉..", false, false)
	end
}
