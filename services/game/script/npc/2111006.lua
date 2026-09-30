-- NPC name (String.wz/Npc.img.xml): 파웬

return {
	on_click = function(me, npc)
		local q3320 = me:quest(3320)
		local q3321 = me:quest(3321)
		local q3353 = me:quest(3353)
		local q3354 = me:quest(3354)
		local q3321_none = not q3321:started() and not q3321:completed()
		local q3354_none = not q3354:started() and not q3354:completed()
		if not (q3320:started() or (q3320:completed() and q3321_none) or q3353:started() or (q3353:completed() and q3354_none)) then
			me:dialog(npc, "으으.. 이놈의 사이티들, 너무 무서워~...")
			return
		end
		if not me:dialog_accept(npc, "지금 당장 연금술사를 보러 가보겠나?") then
			me:dialog(npc, "엥? 싫은가? 자네가 싫다면 하는 수 없지만... 그럼 여기서 연구하던 연금술사는 알려줄 수가 없는데?")
			return
		end
		me:map(926120200)
	end
}
