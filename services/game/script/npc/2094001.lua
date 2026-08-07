-- NPC name (String.wz/Npc.img.xml): 우양

local pq = require("script/lib/party_quest")

local RANKING_QUEST = 1204
local EXIT_MAP = 925100700
local REWARD_MAP = 925100600
local BOSS_MAP = 925100500
local KEY_ID = 4001117
local SEAL_A = 4001120
local SEAL_B = 4001121
local SEAL_C = 4001122

local function strip_pq_items(me)
	pq.remove_all(KEY_ID, me)
	pq.remove_all(SEAL_A, me)
	pq.remove_all(SEAL_B, me)
	pq.remove_all(SEAL_C, me)
end

local function clear_exp(sm)
	local players = sm:players()
	local n = #players
	if n <= 0 then
		return 80000
	end
	local over120 = 0
	local total = 0
	for _, p in ipairs(players) do
		if p ~= nil then
			local lv = p:level()
			total = total + lv
			if lv >= 120 then
				over120 = over120 + 1
			end
		end
	end
	local avg = total / n
	if over120 == 0 then
		return 80000
	end
	if avg <= 120 then
		return 90000
	end
	if avg <= 200 then
		return 90200
	end
	return 100000
end

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		if wz == nil then
			return
		end
		local map_id = wz.id
		if map_id == BOSS_MAP then
			local sm = me:state_machine()
			if sm == nil then
				return
			end
			if not pq.is_leader(me) then
				me:dialog(npc, "당신들의 선장에게 말을 걸라고 하시게.")
				return
			end
			if sm:get_property("clearstage") ~= "1" then
				local exp = clear_exp(sm)
				if exp > 80000 then
					sm:notice("파티원 중 120 레벨을 넘는 플레이어가 있어 보상 경험치량이 증가하였습니다.")
				end
				pq.party_exp(sm, exp)
				for _, p in ipairs(sm:players()) do
					if p ~= nil then
						p:end_party_quest(RANKING_QUEST)
					end
				end
				sm:set_property("clearstage", "1")
				me:dialog(npc, "구해줘서 정말 고맙네. 이제까지 겪어보지 못한 위험이었지만, 도라지들을 구하고 해적의 손길에서 풀어줘서 진심으로 고맙네. 나에게 다시 한번 말을 걸면 내보내 주겠네.")
				return
			end
			pq.party_warp(sm, REWARD_MAP)
			return
		end
		local sel = me:dialog_list(npc, "#b해적왕#k을 물리치고 도라지들과 저를 구해주셔서 정말 감사합니다. 무엇을 도와드릴까요?", {
			"해적왕을 물리친 횟수 확인",
			"이곳에서 내보내 주세요.",
		})
		if sel == 1 then
			local q = me:quest(RANKING_QUEST)
			local cmp = 0
			if q ~= nil then
				cmp = tonumber(q:record_ex("cmp")) or 0
			end
			me:dialog(npc, "#b#h0##k 님은 지금까지 해적왕을 #b"
				.. tostring(cmp)
				.. "#k번 물리쳤습니다. 수고하셨습니다.")
		elseif sel == 2 then
			strip_pq_items(me)
			me:map(EXIT_MAP)
		end
	end
}
