-- NPC name (String.wz/Npc.img.xml): 구옹

local pq = require("script/lib/party_quest")

local GROUP_NAME = "pirate_party_quest"
local MIN_PARTY_SIZE = 3
local MIN_LEVEL = 55
local SCALE_LEVEL = 250
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

local function try_start(me, npc)
	strip_pq_items(me)
	local party = me:party()
	if party == nil or party:leader_id() ~= me:id() then
		me:dialog(npc, "파티가 없거나, 파티장이 아니기 때문에 퀘스트를 시작하실 수 없습니다.")
		return
	end
	local map = me:map()
	if map == nil then
		return
	end
	local wz = map:wz()
	if wz == nil then
		return
	end
	local map_id = wz:id()
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
				size = size + 4
			else
				size = size + 1
			end
		end
	end
	if not ok or size < MIN_PARTY_SIZE then
		me:dialog(npc, "3명 이상이 이곳에 있어야 하며, 파티원들의 레벨이 55이상이 맞는지 확인해 주세요.")
		return
	end
	local group = state_machine(GROUP_NAME)
	if group == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		return
	end
	local state = group:get_property("state")
	if state ~= nil and state ~= "" and state ~= "0" then
		me:dialog(npc, "이미 다른 파티가 입장하여 클리어에 도전하는 중입니다. 나중에 다시 시도해 보세요.")
		return
	end
	local sm, err = group:start_party(me, party, SCALE_LEVEL)
	if sm == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		if err ~= nil then
			log("pirate_party_quest start_party:", err)
		end
	end
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "#e<파티퀘스트 : 해적 데비존>#n\r\n무엇을 원하는건가.", {
			"함께 할 파티원을 찾고 싶어요.",
			"설명을 듣고 싶어요.",
			"파티퀘스트를 하고 싶어요.",
		})
		if sel == 1 then
			me:dialog(npc, "현재 미구현된 컨텐츠입니다.\r\n\r\n컨텐츠 오픈일 : #b하루속히 빨리 오픈하겠습니다.#k\r\n\r\n파티를 구하고 싶으시면 #b파티원 찾기#k 버튼을 누르시면 됩니다.")
		elseif sel == 2 then
			me:dialog(npc, "#e<파티퀘스트 : 해적 데비존>#n\r\n도라지들이 사는 #b백초마을#k에 #r해적 데비존#k이 습격해왔다네. 도라지들의 왕인 #b우양#k님이 납치되셨어. 동료들과 해적선에 침투해 데비존을 몰아내주시오. 부탁하네.\r\n - #e레벨#n : 55이상 #r(추천레벨 : 55 ~ 250)#k\r\n - #e참가인원#n : 3~6명\r\n - #e최종 보상#n : #v2070005# #z2070005##k #b(8회 도와줄 때마다 획득)#k")
		elseif sel == 3 then
			try_start(me, npc)
		end
	end
}
