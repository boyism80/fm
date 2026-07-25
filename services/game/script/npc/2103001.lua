-- NPC name (String.wz/Npc.img.xml): 비밀의 벽

return {
	on_click = function(me, npc)
		local q = me:quest(3927)
		if q == nil or not q:started() then
			return
		end
		if not me:dialog_yes_no(npc, "평범한 벽이지만 자세히 들여다보니 이상한 문양이 그려져 있습니다. 벽을 살피시겠습니까?") then
			return
		end
		me:dialog(npc, "벽 뒤에는 이상한 단어들이 쓰여져 있다.\r\n\r\n#b철퇴와 단검, 활과 화살만 있다면...", false, true)
		q:record("1")
		q:sync_progress()
		me:show_quest_completion(3927)
	end
}
