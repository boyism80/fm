-- NPC name (String.wz/Npc.img.xml): 토리

local pq = require("script/lib/party_quest")

local GROUP_NAME = "henesys_party_quest"
local BONUS_GROUP = "henesys_party_quest_bonus"
local TOWN_MAP = 100000200
local SHORTCUT_A = 910010100
local SHORTCUT_B = 910010400
local MIN_LEVEL = 10
local MAX_LEVEL = 250
local MIN_PARTY = 3
local MAX_PARTY = 6
local RICE_CAKE = 4001101
local HAT_REWARD = 1002798
local EXIT_TOKEN = 3994012

local CLEAR_ITEMS = {
	4001101,
	4001100,
	4001099,
	4001098,
	4001097,
	4001096,
	4001095,
}

local function remove_clear_items_party(me)
	for _, id in ipairs(CLEAR_ITEMS) do
		pq.remove_all_party(id, me)
	end
end

local function rice_cake_reward(me, npc)
	local quest = me:quest(1200)
	local have = nil
	if quest ~= nil then
		have = quest:record_ex("have")
	end
	local code = me:exchange(
		{ item = { [RICE_CAKE] = 20 } },
		{ item = { [HAT_REWARD] = 1 } }
	)
	if code == ExchangeResult.LackCapacity or code == ExchangeResult.LackCost then
		me:dialog(npc, "떡은 제대로 갖고 계신지, 혹은 인벤토리 공간이 부족하신 건 아닌지 확인해주세요.")
		return
	end
	if code ~= ExchangeResult.OK then
		me:dialog(npc, "떡은 제대로 갖고 계신지, 혹은 인벤토리 공간이 부족하신 건 아닌지 확인해주세요.")
		return
	end
	if have == nil and quest ~= nil then
		if not quest:started() then
			quest:start(0, true)
			quest = me:quest(1200)
		end
		if quest ~= nil then
			quest:record_ex("have", "1")
		end
	end
	me:dialog(npc, "떡 20개를 모아오셨군요! 선물로 #b머리위에 떡 하나#k 를 드리도록 할게요!")
end

local function party_ready(me, party, map_id)
	local members = party:members()
	local map = me:map()
	if map == nil then
		return false
	end
	local count = 0
	local in_map = 0
	for _, mem in ipairs(members) do
		if mem ~= nil then
			count = count + 1
			local lv = mem:level()
			if lv < MIN_LEVEL or lv > MAX_LEVEL then
				return false
			end
			if mem:map_id() == map_id and mem:channel_index() ~= nil then
				local ch = map:characters()[mem:id()]
				if ch ~= nil then
					if pq.is_gm(ch) then
						in_map = in_map + 6
					else
						in_map = in_map + 1
					end
				end
			end
		end
	end
	if count > MAX_PARTY then
		return false
	end
	return in_map >= MIN_PARTY
end

local function group_free(name)
	local group = state_machine(name)
	if group == nil then
		return false
	end
	local state = group:get_property("state")
	return state == nil or state == "" or state == "0"
end

function on_click(me, npc)
	local map = me:map()
	if map == nil then
		return
	end
	local wz = map:wz()
	if wz == nil then
		return
	end
	local map_id = wz.id

	if map_id == SHORTCUT_A or map_id == SHORTCUT_B then
		local selected = me:dialog_list(npc, "무엇을 도와 드릴까요?", {
			"이곳에서 나가고 싶어요",
		})
		if selected == nil then
			return
		end
		me:map(TOWN_MAP)
		pq.gain_item(me, EXIT_TOKEN, 1)
		remove_clear_items_party(me)
		return
	end

	local party = me:party()
	if party == nil then
		local selected = me:dialog_list(npc, "안녕하세요? 저는 토리라고 합니다. 이 안은 달맞이꽃이 피어나는 아름다운 언덕이에요. 그런데 그 곳에 살고 있는 어흥이라는 호랑이가 몹시 배가 고파 먹을 것을 찾고 있다고 하네요.", {
			"떡 20개를 가져 왔어요.",
		})
		if selected == 0 then
			rice_cake_reward(me, npc)
		end
		return
	end

	if party:leader_id() ~= me:id() then
		local selected = me:dialog_list(npc, "퀘스트에 도전해 보고 싶다면 #b파티장#k에게 제게 말을 걸어달라고 해주세요.", {
			"떡 20개를 가져 왔어요.",
		})
		if selected == 0 then
			rice_cake_reward(me, npc)
		end
		return
	end

	if not party_ready(me, party, map_id) then
		local selected = me:dialog_list(npc, "퀘스트에 도전하려면 다음과 같은 조건을 만족시켜야 합니다\r\n\r\n#r필요조건: 최소 "
			.. MIN_PARTY
			.. " 명의 파티, 레벨제한 : "
			.. MIN_LEVEL
			.. " ~ "
			.. MAX_LEVEL, {
			"떡 20개를 가져 왔어요.",
		})
		if selected == 0 then
			rice_cake_reward(me, npc)
		end
		return
	end

	local main = state_machine(GROUP_NAME)
	local bonus = state_machine(BONUS_GROUP)
	if main == nil or bonus == nil then
		me:dialog(npc, "퀘스트에 현재 오류가 있습니다.")
		return
	end
	if not group_free(GROUP_NAME) or not group_free(BONUS_GROUP) then
		local selected = me:dialog_list(npc, "이미 다른 파티가 이 안에 들어가서 퀘스트에 도전중입니다. 잠시 후 다시 시도해 주세요.", {
			"떡 20개를 가져 왔어요.",
		})
		if selected == 0 then
			rice_cake_reward(me, npc)
		end
		return
	end

	local sm, err = main:start_party(me, party, 200)
	if sm == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		if err ~= nil then
			log("henesys pq start_party:", err)
		end
		return
	end
	remove_clear_items_party(me)
end
