-- NPC name (String.wz/Npc.img.xml): 미쉘론

local COST = 1800000
local DAILY_LIMIT = 1
local RECORD = "summon.spirit_of_harmony"

return {
	on_click = function(me, npc)
		local count = me:records():get(RECORD)

		if not me:dialog_yes_no(npc, "#r[주의]#b 24일 클라를 받고 소환하셔야합니다. #k 이곳엔 행운을 준다는 정령이 살고있는데 벌목을 하다보니 난폭해져버렸어. 정령을 설득시켜볼래? 그럼 " .. COST .. " 메소가 필요해.\r\n\r\n#r#e" .. count .. "/" .. DAILY_LIMIT .. "#n#k") then
			return
		end
		if count >= DAILY_LIMIT then
			me:dialog(npc, "역시 설득에 실패했구나...왜 그런 표정이지? 절대로 내가 벌목한건 아니야!")
			return
		end
		if me:meso() < COST then
			me:dialog(npc, "난폭한 정령을 부르기엔 메소가 부족 합니다.")
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			me:dialog(npc, "이미 맵 안에 무언가 소환 되어있습니다. 모든 몬스터를 처리해주세요.")
			return
		end

		local code = me:exchange({ meso = COST }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "난폭한 정령을 부르기엔 메소가 부족 합니다.")
			return
		end
		map:spawn_mob(8644011, 962, 404)
		me:records():add(RECORD, 1, { daily = true })
	end
}
