-- NPC name (String.wz/Npc.img.xml): 피터

function on_click(me, npc)
	me:dialog(npc, "어디로 가야할 지 모르겠다면 내가 알려 줄게.", false, false)
end
