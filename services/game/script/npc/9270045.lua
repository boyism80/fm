-- NPC name (String.wz/Npc.img.xml): 사령관 짐

local HAMMER = 4031942

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local d = datetime()
		local today = string.format("%04d%02d%02d", d.year, d.month, d.day)
		local count_q = me:quest(9999998)
		local day_q = me:quest(999998)
		if not count_q:started() then
			count_q:start("0")
		end
		if not day_q:started() then
			day_q:start("")
		end
		if day_q:record() ~= today then
			count_q:record("0")
			day_q:record(today)
		end

		local sel = me:dialog_list(npc, "크렉셀을 잡으시려고?", { "크렉셀 을 소환한다." })
		if sel ~= 1 then
			return
		end
		if item_count(me, HAMMER) < 1 then
			me:dialog(npc, "망치 1개가 필요합니다.")
			return
		end
		if (tonumber(count_q:record()) or 0) >= 10 then
			me:dialog(npc, "오늘 이미 소환하셨네요? 다음에 시도해주세요.~")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "맵 안에 아직 몬스터 있습니다.~")
			return
		end

		local code = me:exchange({ item = { [HAMMER] = 1 } }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "망치 1개가 필요합니다.")
			return
		end
		count_q:record("1")
		map:spawn_mob(9420520, -178, -212)
	end
}
