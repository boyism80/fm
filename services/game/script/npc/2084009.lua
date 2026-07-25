-- NPC name (String.wz/Npc.img.xml): OX새

return {
	on_click = function(me, npc)
		me:dialog(npc, ".", false, false)
	end
}
