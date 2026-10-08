-- NPC name (String.wz/Npc.img.xml): 예쁜 할로윈 마녀

local COST = 850000
local DAILY_LIMIT = 1
local RECORD = "summon.gleam_eyes"

return {
	on_click = function(me, npc)
		local count = me:records():get(RECORD)

		if not me:dialog_yes_no(npc, "어서와. 이곳엔 무서운 몬스터가 있어. 뭐? 너는 그 몬스터를 찾고있다고? 공짜로 가르쳐 줄 순없고 " .. COST .. " 메소를 줘.\r\n\r\n#r#e" .. count .. "/" .. DAILY_LIMIT .. "#n#k") then
			return
		end
		if count >= DAILY_LIMIT then
			me:dialog(npc, "그 몬스터가 또있냐고? 몰라. 내일 또와봐.")
			return
		end
		if me:meso() < COST then
			me:dialog(npc, "더 글림아이즈를 부르기엔 메소가 부족 합니다.")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "이미 맵 안에 무언가 소환 되어있습니다. 모든 몬스터를 처리해주세요.")
			return
		end

		local code = me:exchange({ meso = COST }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "더 글림아이즈를 부르기엔 메소가 부족 합니다.")
			return
		end
		map:spawn_mob(9401011, -615, 336)
		me:records():add(RECORD, 1, { daily = true })
	end
}
