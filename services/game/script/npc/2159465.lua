-- NPC name (String.wz/Npc.img.xml): 라니아

local COST = 10000000
local DAILY_LIMIT = 1

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

		if not me:dialog_yes_no(npc, "이곳은 어떤 소녀의 불안정한 꿈속이에요. 그아이를 구해주시겠어요? 그럴려면 " .. COST .. " 메소 가 필요합니다.\r\n\r\n#r#e" .. count .. "/" .. DAILY_LIMIT .. "#n#k") then
			return
		end
		if count >= DAILY_LIMIT then
			me:dialog(npc, "당신은 그녀를 구했습니다. 진심으로 감사를 표하는 바입니다.[내일 또 도전가능해요.]")
			return
		end
		if me:meso() < COST then
			me:dialog(npc, "루시드를 부르기엔 메소가 부족 합니다.")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "이미 맵 안에 무언가 소환 되어있습니다. 모든 몬스터를 처리해주세요.")
			return
		end

		local code = me:exchange({ meso = COST }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "루시드를 부르기엔 메소가 부족 합니다.")
			return
		end
		map:spawn_mob(9880140, 1032, 48)
		count_q:record(tostring(count + 1))
	end
}
