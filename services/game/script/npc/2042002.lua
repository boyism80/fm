-- NPC name (String.wz/Npc.img.xml): 슈피겔만 - 몬스터 카니발

local HUB_MAP = 980000000
local MIN_LEVEL = 30
local MAX_LEVEL = 50
local RANKING_QUEST = 1301
local COIN_ID = 4001129

local function reward_for_cp(total_cp, winner)
	if total_cp >= 501 then
		if winner then
			return 30000, "A"
		end
		return 25500, "A"
	end
	if total_cp >= 251 then
		if winner then
			return 25500, "B"
		end
		return 20500, "B"
	end
	if total_cp >= 50 then
		if winner then
			return 21000, "C"
		end
		return 17000, "C"
	end
	if winner then
		return 3000, "D"
	end
	return 15000, "D"
end

local function on_reward(me, npc)
	local team = me:carnival_team()
	if team == nil then
		me:map(HUB_MAP, 0)
		return
	end
	local winner = team:is_winner()
	local exp, rank = reward_for_cp(team:total_cp(), winner)
	local message = "아쉽게도 비기거나 지고 말았군. 승리를 위해 좀 더 노력해주게! \r\n#b랭크 : " .. rank
	if winner then
		message = "축하하네. 카니발에서 승리를 거두었군. 자네들의 활약은 잘 지켜보았다네. 다음 번에도 기대하겠네! \r\n#b랭크 : " .. rank
	end
	if not me:dialog(npc, message) then
		return
	end
	me:end_party_quest(RANKING_QUEST)
	me:exchange({}, { exp = exp })
	team:remove_member(me)
	me:map(HUB_MAP, 0)
end

return {
	on_click = function(me, npc)
		local map_id = me:map():wz():id()
		if map_id == 980000010 then
			me:map(HUB_MAP, 0)
			return
		end
		if map_id >= 980000000 and (map_id % 10 == 3 or map_id % 10 == 4) then
			on_reward(me, npc)
			return
		end
		local sel = me:dialog_list(npc, "무얼 하겠는가?", {
			"몬스터 카니발 필드로 이동한다.",
			"몬스터 카니발에 대한 설명을 듣는다.",
			"메이플 코인을 교환한다.",
		})
		if sel == nil then
			return
		end
		if sel == 1 then
			if me:level() < MIN_LEVEL or me:level() > MAX_LEVEL then
				me:dialog(npc, "레벨 30 이상, 50 이하의 캐릭터만 몬스터 카니발을 즐길 수 있다네.")
				return
			end
			me:save_location("MONSTERCARNIVAL")
			me:map(HUB_MAP, 0)
		elseif sel == 2 then
			me:dialog(npc, "직접 그 전율을 느껴보기 전에는 이것이 무엇인지 알 수 없지.")
		elseif sel == 3 then
			local code = me:exchange({ item = { [COIN_ID] = 50 } }, { item = { [1122007] = 1 } })
			if code == ExchangeResult.LackCapacity then
				me:dialog(npc, "장비창에 빈 칸이 없는 것 같군. 확인해 보게나.")
				return
			end
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "#b#t" .. COIN_ID .. "##k이 부족한 것 같군.")
				return
			end
			me:dialog(npc, "여기 #b#t1122007##k일세. 잘 사용하게나.")
		end
	end
}
