-- NPC name (String.wz/Npc.img.xml): 토미

local pq = require("script/lib/party_quest")

local GROUP_NAME = "henesys_party_quest_bonus"
local BONUS_MAP = 910010200
local SHORTCUT_MAP = 910010100
local EXIT_PATH_MAP = 910010300
local TOWN_MAP = 100000200
local SHORTCUT_OUT = 910010400

local CLEAR_ITEMS = {
	4001101,
	4001100,
	4001099,
	4001098,
	4001097,
	4001096,
	4001095,
}

local function remove_clear_items(me)
	for _, id in ipairs(CLEAR_ITEMS) do
		pq.remove_all(id, me)
	end
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

		if map_id == SHORTCUT_MAP then
			me:dialog(npc, "안녕하세요? 저는 토미입니다. 이곳 근처에는 돼지마을이 있습니다. 그 곳의 돼지들은 성격이 포악하고 욕심이 유달리 많아 여행자들이 가지고 다니던 각종 무기들을 빼앗아 마을에서 추방되어 이곳에 숨어 지내고 있어요.")
			local selected = me:dialog_list(npc, "파티원들과 함께 그 곳으로 여행을 떠나 못된 돼지들을 혼내 주시는 것은 어떨까요?", {
				"네, 그곳으로 보내주세요.",
			})
			if selected == nil then
				return
			end
			local party = me:party()
			if party == nil or party:leader_id() ~= me:id() then
				me:dialog(npc, "돼지 마을은 파티장이 입장을 신청할 수 있어요.")
				return
			end
			local group = state_machine(GROUP_NAME)
			if group == nil then
				me:dialog(npc, "This PQ is not currently available.")
				return
			end
			local state = group:get_property("state")
			if state ~= nil and state ~= "" and state ~= "0" then
				me:dialog(npc, "으음..이미 이 안에 다른 파티가 들어간 것 같은데요?")
				return
			end
			local prev = me:state_machine()
			if prev ~= nil then
				prev:finish(0)
			end
			local sm, err = group:start_party(me, party, 200)
			if sm == nil then
				me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
				if err ~= nil then
					log("henesys bonus start_party:", err)
				end
			end
			return
		end

		if map_id == BONUS_MAP then
			local selected = me:dialog_list(npc, "무엇을 도와 드릴까요?", {
				"이곳에서 나가고 싶어요",
			})
			if selected == nil then
				return
			end
			local spawn = 0
			local group = state_machine(GROUP_NAME)
			if group ~= nil then
				local dest = group:map(SHORTCUT_OUT)
				if dest ~= nil then
					local portal = dest:portal("st00")
					if portal ~= nil then
						spawn = portal:id()
					end
					me:map(dest, spawn)
					return
				end
			end
			me:map(SHORTCUT_OUT)
			return
		end

		if map_id == EXIT_PATH_MAP then
			local selected = me:dialog_list(npc, "무엇을 도와 드릴까요?", {
				"이곳에서 나가고 싶어요",
			})
			if selected == nil then
				return
			end
			remove_clear_items(me)
			me:map(TOWN_MAP)
		end
	end
}
