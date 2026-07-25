-- NPC name (String.wz/Npc.img.xml): 치노

return {
	on_click = function(me, npc)
		me:dialog(npc, "안녕하세요! 달맞이언덕에 달이 떴을때의 정취를 아시나요? 모르신다구요? 에...그럼 저도몰라요!", false, false)
	end
}
