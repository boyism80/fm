-- NPC name (String.wz/Npc.img.xml): 바로크

function on_click(me, npc)
	me:dialog(npc, "나는 다 보았다...네가 무엇을 잘못했는지는 잘 알고있겠지? 이번엔 그냥 넘어가지만 다음에 또 걸리면 처벌을 면치 못할것이니 주의하라고!", false, false)
end
