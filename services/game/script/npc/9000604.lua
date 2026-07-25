-- NPC name (String.wz/Npc.img.xml): 베르누이

return {
	on_click = function(me, npc)
		me:dialog(npc, "제 복장이 좀 야하죠?..후방 주의 하세요.", false, false)
	end
}
