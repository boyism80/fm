-- NPC name (String.wz/Npc.img.xml): 황금열쇠박스

function on_click(me, npc)
	me:dialog(npc, ".", false, false)
end
