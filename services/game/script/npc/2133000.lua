-- NPC name (String.wz/Npc.img.xml): 엘린

local pq = require("script/lib/party_quest")

local GROUP_NAME = "ellin_party_quest"
local MIN_PARTY_SIZE = 3
local MAX_PARTY_SIZE = 5
local MIN_LEVEL = 45
local SCALE_LEVEL = 55
local PURPLE_STONE = 4001163
local MONSTER_MARBLE = 4001169
local PURIFY_BEAD = 2270004
local ALTAIR_FRAGMENT = 4001198
local ALTAIR_EARRING = 1032060
local SHINY_ALTAIR_EARRING = 1032061

local function strip_pq_items(me)
	pq.remove_all(PURPLE_STONE, me)
	pq.remove_all(MONSTER_MARBLE, me)
	pq.remove_all(PURIFY_BEAD, me)
end

local function claim_altair(me, npc)
	if not pq.has_item(me, ALTAIR_EARRING) and pq.has_item(me, ALTAIR_FRAGMENT, 10) then
		local code = me:exchange(
			{ item = { [ALTAIR_FRAGMENT] = 10 } },
			{ item = { [ALTAIR_EARRING] = 1 } }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "인벤토리 공간을 확보한 뒤 다시 말을 걸어줘.")
			return
		end
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "알테어 조각 10개를 제대로 갖고 있지 않거나, 알테어 이어링을 이미 갖고 있는 것 같은데?")
		end
		return
	end
	me:dialog(npc, "알테어 조각 10개를 제대로 갖고 있지 않거나, 알테어 이어링을 이미 갖고 있는 것 같은데?")
end

local function claim_shiny_altair(me, npc)
	if pq.has_item(me, ALTAIR_EARRING)
		and not pq.has_item(me, SHINY_ALTAIR_EARRING)
		and pq.has_item(me, ALTAIR_FRAGMENT, 10) then
		local code = me:exchange(
			{ item = { [ALTAIR_EARRING] = 1, [ALTAIR_FRAGMENT] = 10 } },
			{ item = { [SHINY_ALTAIR_EARRING] = 1 } }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "인벤토리 공간을 확보한 뒤 다시 말을 걸어줘.")
			return
		end
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "알테어 조각 10개를 제대로 갖고 있지 않거나, 알테어 이어링을 이미 갖고 있는 것 같은데?")
		end
		return
	end
	me:dialog(npc, "알테어 조각 10개를 제대로 갖고 있지 않거나, 알테어 이어링을 이미 갖고 있는 것 같은데?")
end

local function try_start(me, npc)
	local party = me:party()
	if party == nil or party:leader_id() ~= me:id() then
		me:dialog(npc, "파티를 만들고 파티장이 되면 나를 통해 퀘스트에 입장할 수 있어.")
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
	local map_id = map_wz:id()
	local ok = true
	local size = 0
	for _, mem in ipairs(party:members()) do
		if mem ~= nil then
			if mem:map_id() ~= map_id or mem:level() < MIN_LEVEL then
				ok = false
				break
			end
			local ch = map:characters()[mem:id()]
			if ch == nil then
				ok = false
				break
			end
			if pq.is_gm(ch) then
				size = size + 3
			else
				size = size + 1
			end
		end
	end
	if not ok or size < MIN_PARTY_SIZE or size > MAX_PARTY_SIZE then
		me:dialog(npc, "레벨 제한 45 이상의 파티원 3명에서 5명이 퀘스트를 시작할 수 있어.")
		return
	end
	local group = state_machine(GROUP_NAME)
	if group == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		return
	end
	local state = group:get_property("state")
	if state ~= nil and state ~= "" and state ~= "0" then
		me:dialog(npc, "이미 이 안에서 다른 파티가 퀘스트에 도전하는 중입니다. 나중에 다시 시도하세요.")
		return
	end
	local sm, err = group:start_party(me, party, SCALE_LEVEL)
	if sm == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		if err ~= nil then
			log("ellin_party_quest start_party:", err)
		end
	end
end

return {
	on_click = function(me, npc)
		strip_pq_items(me)
		local sel = me:dialog_list(npc, "#e<파티퀘스트 : 독안개의 숲>#n\r\n괴인에 의해 오염되어 버린 숲을 구해야 해! 하지만 알테어 캠프의 용사들은 모두 생업에 바빠서 움직일 수 없어. 네가 도와주지 않을래? #b45레벨 이상의 모험가#k라면, 도와줄 수 있을 거야!", {
			"알테어 이어링을 얻고 싶어.",
			"빛나는 알테어 이어링을 얻고 싶어.",
			"독안개의 숲으로 향하고 싶어.",
		})
		if sel == 1 then
			claim_altair(me, npc)
		elseif sel == 2 then
			claim_shiny_altair(me, npc)
		elseif sel == 3 then
			try_start(me, npc)
		end
	end
}
