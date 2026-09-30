-- NPC name (String.wz/Npc.img.xml): 포포

local COST = 1500000
local DAILY_LIMIT = 1

return {
	on_click = function(me, npc)
		local d = datetime()
		local today = d.year .. "" .. d.month .. "" .. d.day
		local count_q = me:quest(12191053)
		local day_q = me:quest(12191054)
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

		if not me:dialog_yes_no(npc, "으...요즘 뒷골목 깡패들이 너무 어슬렁대... 어?! 너는 누구니? 혹시 나대신 돈좀 내주겠어? 나는 " .. COST .. " 메소가 필요해.\r\n\r\n#r#e" .. count .. "/" .. DAILY_LIMIT .. "#n#k") then
			return
		end
		if count >= DAILY_LIMIT then
			me:dialog(npc, "그 깡패가 또있냐고? 몰라. 내일 또와봐.")
			return
		end
		if me:meso() < COST then
			me:dialog(npc, "부하를 부르기엔 메소가 부족 합니다.")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "이미 맵 안에 무언가 소환 되어있습니다. 모든 몬스터를 처리해주세요.")
			return
		end

		local code = me:exchange({ meso = COST }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "부하를 부르기엔 메소가 부족 합니다.")
			return
		end
		map:spawn_mob(9400103, 777, -1332)
		count_q:record(tostring(count + 1))
	end
}
