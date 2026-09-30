-- NPC name (String.wz/Npc.img.xml): 고스트헌터 밥

return {
	on_click = function(me, npc)
		local q = me:quest(3250)
		local count = 0
		for _, it in pairs(me:item(4220046)) do
			count = count + it:count()
		end
		if q == nil or not q:started() or count < 1 then
			me:dialog(npc, "냠냠.. 햄버거가 제일 맛있어.")
			return
		end
		if not me:dialog_yes_no(npc, "정말 타이머 키우기를 포기할거야?") then
			return
		end
		me:exchange({ item = { [4220046] = 1 } }, nil)
		q:forfeit()
		me:dialog(npc, "수고했어. 다시 키우고 싶어지면 언제든지 나를 찾아와.")
	end
}
