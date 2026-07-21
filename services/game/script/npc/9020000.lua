-- NPC name (String.wz/Npc.img.xml): 라케리스

local pq = require("script/lib/party_quest")
local config = require("script/lib/kerning_city_party_quest")

local GROUP_NAME = "kerning_city_party_quest"
local MIN_LEVEL = 21
local COUPON_ID = 4001007
local PASS_ID = 4001008

function on_click(me, npc)
	local party = me:party()
	if party == nil then
		me:dialog(npc, "파티원들과 함께 힘을 모아 퀘스트에 해결해 보시고 싶지 않으세요? 이 안에는 서로 힘을 합치지 않으면 해결할 수 없는 장애물들이 많이 있답니다. 언제든지 도전해 보고 싶으시면 #b파티장#k에게 제게 말을 걸어 달라고 해주세요.")
		return
	end
	if party:leader_id() ~= me:id() then
		me:dialog(npc, "퀘스트를 시작하고 싶으세요? 그렇다면 당신의 파티장에게 말을 걸어달라고 해주세요.")
		return
	end

	local map = me:map()
	if map == nil then
		return
	end
	local map_wz = map:wz()
	if map_wz == nil then
		return
	end
	local map_id = map_wz.id
	local members = party:members()
	local ok = true
	local in_map = 0
	local count = 0
	for _, mem in ipairs(members) do
		if mem ~= nil then
			count = count + 1
			if mem:level() < MIN_LEVEL then
				ok = false
			end
			if mem:map_id() == map_id and mem:channel_index() ~= nil then
				local ch = map:characters()[mem:id()]
				if ch ~= nil then
					in_map = in_map + 1
				end
			end
		end
	end
	if count ~= config.required_party_size
		or in_map ~= config.required_party_size then
		ok = false
	end
	if not ok then
		me:dialog(npc, "당신이 속한 파티의 파티원이 "
			.. config.required_party_size
			.. "명이 아니거나 자신 혹은 파티원 중에서 레벨 21이상이 아닌 캐릭터가 있습니다. 혹은 파티원 전원이 현재 맵에 모여있는지 다시 한 번 확인해 주세요.")
		return
	end

	local group = state_machine(GROUP_NAME)
	if group == nil then
		me:dialog(npc, "This PQ is not currently available.")
		return
	end
	local state = group:get_property("state")
	if state == nil or state == "" or state == "0" then
		local sm, err = group:start_party(me, party, 30)
		if sm == nil then
			me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
			if err ~= nil then
				log("kpq start_party:", err)
			end
		end
	else
		me:dialog(npc, "이 안에 이미 다른 파티가 들어가서 클리어에 도전하고있습니다. 잠시후에 다시 시도해 주세요")
	end
	pq.remove_all(PASS_ID, me)
	pq.remove_all(COUPON_ID, me)
end
