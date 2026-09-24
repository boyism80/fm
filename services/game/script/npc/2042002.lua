-- NPC name (String.wz/Npc.img.xml): 슈피겔만 - 몬스터 카니발

local pq = require("script/lib/party_quest")
local HUB_MAP = 980000000
local MIN_LEVEL = 30
local MAX_LEVEL = 50
local RANKING_QUEST = 1301
local COIN_ID = 4001129

local function exp_for_cp(total_cp, winner)
	if total_cp >= 501 then
		if winner then
			return 30000
		end
		return 25500
	end
	if total_cp >= 251 then
		if winner then
			return 25500
		end
		return 20500
	end
	if total_cp >= 50 then
		if winner then
			return 21000
		end
		return 17000
	end
	if winner then
		return 3000
	end
	return 15000
end

local function on_reward(me, npc)
	local team = me:carnival_team()
	if team == nil then
		me:map(HUB_MAP, 0)
		return
	end
	local total_cp = team:total_cp()
	local winner = team:is_winner()
	me:end_party_quest(RANKING_QUEST)
	me:exchange({}, { exp = exp_for_cp(total_cp, winner) })
	team:clear()
	me:map(HUB_MAP, 0)
end

return {
	on_click = function(me, npc)
		local map_id = me:map_id()
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
			if not pq.has_item(me, COIN_ID, 50) then
				me:dialog(npc, "#b#t" .. COIN_ID .. "##k이 부족하거나, 장비창에 빈 칸이 없는건 아닌가?")
				return
			end
			me:exchange({ items = { { id = COIN_ID, count = 50 } } }, { items = { { id = 1122007, count = 1 } } })
		end
	end
}
