-- NPC name (String.wz/Npc.img.xml): 혼테일의 이정표

local pq = require("script/lib/party_quest")

local GROUP_NAME = "horntail_party_quest"
local ENTRY_MAP = 240050000
local MAZE_MAP = 240050100
local LIGHT_CAVE = 240050300
local DARK_CAVE = 240050310
local NEXT_STAGE = 240050200
local PARTY_SIZE = 6

local KEY_ITEMS = {
	4001087,
	4001088,
	4001089,
	4001090,
	4001091,
	4001092,
	4001093,
}

local function remove_keys(me)
	for _, id in ipairs(KEY_ITEMS) do
		pq.remove_all(id, me)
	end
end

local function party_ready(me, party, map_id)
	local members = party:members()
	local count = 0
	local in_map = 0
	for _, mem in ipairs(members) do
		if mem ~= nil then
			count = count + 1
			if mem:map_id() == map_id and mem:channel_index() ~= nil then
				in_map = in_map + 1
			end
		end
	end
	if in_map ~= count then
		return false
	end
	if pq.is_gm(me) then
		return true
	end
	return count == PARTY_SIZE
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
		local map_id = wz:id()

		if map_id == ENTRY_MAP then
			local party = me:party()
			if party == nil then
				me:dialog(npc, "만용을 부리는군. 어리석은 자들이여.. 강한자들과 함께 도전하라.")
				return
			end
			if party:leader_id() ~= me:id() then
				me:dialog(npc, "만용을 부리는군. 어리석은 자여.. 강한자들과 함께 도전하라.")
				return
			end
			if not party_ready(me, party, map_id) then
				me:dialog(npc, "만용을 부리는군. 어리석은 자들이여.. 강한자들과 함께 도전하라.")
				return
			end
			local selected = me:dialog_list(npc, "겁없이 생명의 동굴로 발을 내딛은 어리석은 자들이여... 숨겨진 열쇠를 찾은 자만이 나에게 다가올 수 있을것이다. 무모한 게임에 도전하겠는가?", {
				"도전한다.",
			})
			if selected == nil then
				return
			end
			local group = state_machine(GROUP_NAME)
			if group == nil then
				me:dialog(npc, "미안하지만 파티퀘스트 시스템에 현재 문제가 생겼다. 지금은 입장할 수 없다..")
				return
			end
			local state = group:get_property("state")
			if state ~= nil and state ~= "" and state ~= "0" then
				me:dialog(npc, "이 안에 이미 다른 파티가 입장하여 클리어에 도전중입니다. 잠시 후에 다시 시도해보세요.")
				return
			end
			local sm, err = group:start_party(me, party, 200)
			if sm == nil then
				me:dialog(npc, "미안하지만 파티퀘스트 시스템에 현재 문제가 생겼다. 지금은 입장할 수 없다..")
				if err ~= nil then
					log("horntail pq start_party:", err)
				end
				return
			end
			remove_keys(me)
			return
		end

		if map_id == MAZE_MAP then
			local sm = me:state_machine()
			if sm == nil then
				return
			end
			local progress = tonumber(sm:get_property("stage1progress")) or 0
			if progress == 5 then
				pq.party_warp(sm, NEXT_STAGE)
			else
				me:dialog(npc, "만용을 부리는군. 어리석은 자들이여.. 아직 게임은 끝나지 않았다.")
			end
			return
		end

		if map_id == LIGHT_CAVE or map_id == DARK_CAVE then
			local sm = me:state_machine()
			if sm == nil then
				return
			end
			if not pq.has_item(me, 4001093, 6) then
				me:dialog(npc, "만용을 부리는군. 어리석은 자들이여.. 아직 게임은 끝나지 않았다.")
				return
			end
			sm:set_property("allfinish", "clear")
			pq.remove_all(4001093, me)
			local m = me:map()
			if m ~= nil then
				m:clear_effect()
				m:play_sound("Party1/Clear")
			end
			sm:finish(240050400)
		end
	end
}
