-- NPC name (String.wz/Npc.img.xml): 오르카

local COST = 600000
local DAILY_LIMIT = 15

return {
	on_click = function(me, npc)
		local d = datetime()
		local today = string.format("%04d%02d%02d", d.year, d.month, d.day)
		local count_q = me:quest(12191027)
		local day_q = me:quest(12191028)
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

		if not me:dialog_yes_no(npc, "이곳은 그분을 위한 장소...너는...누구? 혹시 날 좀 도와주겠니? 내가 기다리고 있는 사람이 끔찍한 독기때문에 올 수 없는 상황이라서...도와줄 수 있다면 " .. COST .. " 메소를 줘.\r\n\r\n#r#e" .. count .. "/" .. DAILY_LIMIT .. "#n#k") then
			return
		end
		if count >= DAILY_LIMIT then
			me:dialog(npc, "덕분에 기다리던 사람이 올 수 있을지도몰라요! 진심으로 감사를 표하는 바입니다.[내일 또 도전가능해요.]")
			return
		end
		if me:meso() < COST then
			me:dialog(npc, "카오스대왕지네를 부르기엔 메소가 부족 합니다.")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "이미 맵 안에 무언가 소환 되어있습니다. 모든 몬스터를 처리해주세요.")
			return
		end

		local code = me:exchange({ meso = COST }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "카오스대왕지네를 부르기엔 메소가 부족 합니다.")
			return
		end
		map:spawn_mob(5000006, 2331, 823)
		count_q:record(tostring(count + 1))
	end
}
