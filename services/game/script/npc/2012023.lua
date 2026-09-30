-- NPC name (String.wz/Npc.img.xml): 단풍잎 구슬

return {
	on_click = function(me, npc)
		if not me:quest(6230):started() then
			me:dialog(npc, "투명한 구슬 속에 붉은 단풍잎이 있다.")
			return
		end
		if next(me:item(4031456)) ~= nil then
			return
		end
		if next(me:item(4031476)) == nil then
			return
		end

		local code = me:exchange({ item = { [4031476] = 1 } }, { item = { [4031456] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족합니다.")
			return
		end
		me:dialog(npc, "단풍잎이 빛나는 유리구슬 속으로 빨려들어갔습니다.")
	end
}
