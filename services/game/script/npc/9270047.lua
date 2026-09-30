-- NPC name (String.wz/Npc.img.xml): 알돌

local COIN = 3980018

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

		local sel = me:dialog_list(npc, "환상의 나라 최강의 무투가 타르가 에게 도전하시겠어요??", { "타르가와 붙어본다.." })
		if sel ~= 1 then
			return
		end
		if item_count(me, COIN) < 1 then
			me:dialog(npc, "타르가 코인 1개가 필요합니다.")
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

		local code = me:exchange({ item = { [COIN] = 1 } }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "타르가 코인 1개가 필요합니다.")
			return
		end
		count_q:record("1")
		map:spawn_mob(9420541, -397, 640)
	end
}
