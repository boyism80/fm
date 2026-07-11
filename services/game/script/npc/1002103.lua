-- NPC name (String.wz/Npc.img.xml): 리더 알

function on_click(me, npc)
	me:dialog(npc, "패밀리에 대해 궁금한게 있나?", false, false)
end
