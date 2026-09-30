-- NPC name (String.wz/Npc.img.xml): 초보자 도우미

local COST = 6000000
local DAILY_LIMIT = 6
local ZAKUM_ARMS = {
	8800003, 8800004, 8800005, 8800006,
	8800007, 8800008, 8800009, 8800010,
}

return {
	on_click = function(me, npc)
		local d = datetime()
		local today = d.year .. "" .. d.month .. "" .. d.day
		local count_q = me:quest(12191025)
		local day_q = me:quest(12191026)
		if not count_q:started() then
			count_q:start("0")
		end
		if not day_q:started() then
			day_q:start(today)
		end
		if day_q:record() ~= today then
			count_q:record("0")
			day_q:record(today)
		end
		local count = tonumber(count_q:record()) or 0

		if not me:dialog_yes_no(npc, "자쿰을 소환 하기 위해선 " .. COST .. "가 필요합니다.\r\n\r\n#r#e" .. count .. "/" .. DAILY_LIMIT .. "#n#k") then
			return
		end
		if count >= DAILY_LIMIT then
			me:dialog(npc, "횟수를 모두 소진 했습니다.")
			return
		end
		if me:meso() < COST then
			me:dialog(npc, "소환하기엔 메소가 부족 합니다.")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "이미 맵 안에 무언가 소환 되어있습니다. 모든 몬스터를 처리해주세요.")
			return
		end

		local code = me:exchange({ meso = COST }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "소환하기엔 메소가 부족 합니다.")
			return
		end
		local zakum = map:spawn_mob(8800000, -10, -215, -2)
		if zakum ~= nil then
			zakum:fake(true)
		end
		for _, mob_id in ipairs(ZAKUM_ARMS) do
			map:spawn_mob(mob_id, -10, -215, -2)
		end
		count_q:record(tostring(count + 1))
	end
}
