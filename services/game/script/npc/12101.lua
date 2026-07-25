-- NPC name (String.wz/Npc.img.xml): 레인

return {
	on_click = function(me, npc)
		me:dialog(npc, "이곳 메이플 아일랜드는 위험한 몬스터들이 없어서 안심이에요!", false, false)
	end
}
