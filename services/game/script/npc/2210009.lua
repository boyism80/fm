-- NPC name (String.wz/Npc.img.xml): ???

return {
	on_click = function(me, npc)
		if not me:dialog(npc, "여신의 반려동물이 이토록 사악하고 포악한 몬스터로 탄생하다니...?", false, true) then
			return
		end
		if not me:dialog(npc, "후후후.... 이것은 시작 일 뿐, 앞으로의 일을 더욱 기대해도 좋을 것이다.", true, true) then
			return
		end
		me:dialog(npc, "그럼 이만.", false, false)
	end
}
