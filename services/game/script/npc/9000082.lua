-- NPC name (String.wz/Npc.img.xml): 퐁

local FLAME = 4001433

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "라바나를 잡으시려고?", { "라바나를 소환한다." })
		if sel ~= 1 then
			return
		end
		if item_count(me, FLAME) < 30 then
			me:dialog(npc, "태양의 불꽃 30개가 필요합니다.")
			return
		end
		if me:records():get("summon.ravana") >= 10 then
			me:dialog(npc, "오늘 이미 소환하셨네요? 다음에 시도해주세요.~")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "맵 안에 아직 몬스터 있습니다.~")
			return
		end

		local code = me:exchange({ item = { [FLAME] = 30 } }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "태양의 불꽃 30개가 필요합니다.")
			return
		end
		me:records():add("summon.ravana", 1, { daily = true })
		map:spawn_mob(9500390, 818, -513)
	end
}
