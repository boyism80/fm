-- NPC name (String.wz/Npc.img.xml): 낯익은 처녀

function on_click(me, npc)
	me:dialog(npc, ".", false, false)
end
