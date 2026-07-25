-- NPC name (String.wz/Npc.img.xml): 세자르

return {
	on_click = function(me, npc)
		me:dialog(npc, "도전하라!! 라고 하고싶지만... 아쉽지만 지금은 정비중이네. 다음에 다시 오게.", false, false)
	end
}
