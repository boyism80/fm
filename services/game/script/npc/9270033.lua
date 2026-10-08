-- NPC name (String.wz/Npc.img.xml): 밥

local ESSENCE = 4000381

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "라타니카 선장을 소환하게?", { "라타니카 선장을 소환한다." })
		if sel ~= 1 then
			return
		end
		if item_count(me, ESSENCE) < 1 then
			me:dialog(npc, "화이트 에센스가 필요합니다.")
			return
		end
		if me:records():get("summon.captain_latanica") >= 10 then
			me:dialog(npc, "오늘 이미 소환하셨네요? 다음에 시도해주세요.~")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "맵 안에 아직 몬스터 있습니다.~")
			return
		end

		local code = me:exchange({ item = { [ESSENCE] = 1 } }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "화이트 에센스가 필요합니다.")
			return
		end
		me:records():add("summon.captain_latanica", 1, { daily = true })
		map:spawn_mob(9420513, -178, -212)
	end
}
