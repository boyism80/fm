-- NPC name (String.wz/Npc.img.xml): 나무책상

function on_click(me, npc)
	me:dialog(npc, "...", false, false)
end
