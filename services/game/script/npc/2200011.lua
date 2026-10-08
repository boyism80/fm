-- NPC name (String.wz/Npc.img.xml): 영주의 방 결계

local COST = 10500000
local DAILY_LIMIT = 1
local RECORD = "summon.toad_lord"

return {
	on_click = function(me, npc)
		local count = me:records():get(RECORD)

		if not me:dialog_yes_no(npc, "호...이몸을 찾았단말인가? 도전을 해보겠다고? 그렇다면 재력을 증명해봐. " .. COST .. " 메소를 줘봐.\r\n\r\n#r#e" .. count .. "/" .. DAILY_LIMIT .. "#n#k") then
			return
		end
		if count >= DAILY_LIMIT then
			me:dialog(npc, "오늘은 더이상 네게 볼일이 없다.")
			return
		end
		if me:meso() < COST then
			me:dialog(npc, "나를 를 부르기엔 공물이 부족 하다.")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "이미 맵 안에 무언가 소환 되어있습니다. 모든 몬스터를 처리해주세요.")
			return
		end

		local code = me:exchange({ meso = COST }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "나를 를 부르기엔 공물이 부족 하다.")
			return
		end
		map:spawn_mob(6500011, 772, 96)
		me:records():add(RECORD, 1, { daily = true })
	end
}
