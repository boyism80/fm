-- NPC name (String.wz/Npc.img.xml): 페리

local COST = 3500000
local DAILY_LIMIT = 1
local RECORD = "summon.stalker"

return {
	on_click = function(me, npc)
		local count = me:records():get(RECORD)

		if not me:dialog_yes_no(npc, "으...스토커는 싫어...어? 안녕하세요. 스토커에게 볼 일 이 있다고요? 그녀석은 돈을 좋아해서 " .. COST .. " 정도면 불러낼 수 있어요.\r\n\r\n#r#e" .. count .. "/" .. DAILY_LIMIT .. "#n#k") then
			return
		end
		if count >= DAILY_LIMIT then
			me:dialog(npc, "그 스토커가 또있냐고? 내일 또와봐.")
			return
		end
		if me:meso() < COST then
			me:dialog(npc, "스토커를 부르기엔 메소가 부족 합니다.")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "이미 맵 안에 무언가 소환 되어있습니다. 모든 몬스터를 처리해주세요.")
			return
		end

		local code = me:exchange({ meso = COST }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "스토커를 부르기엔 메소가 부족 합니다.")
			return
		end
		map:spawn_mob(9400120, 777, -1332)
		me:records():add(RECORD, 1, { daily = true })
	end
}
