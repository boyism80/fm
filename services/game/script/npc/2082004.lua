-- NPC name (String.wz/Npc.img.xml): 앤디

function on_click(me, npc)
	me:dialog(npc, "대사 모름.", false, false)
end
