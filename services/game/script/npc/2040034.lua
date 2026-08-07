-- NPC name (String.wz/Npc.img.xml): 표지판

local pq = require("script/lib/party_quest")
local GROUP_NAME = "ludibrium_party_quest"
local MIN_PARTY_SIZE = 3
local MAX_PARTY_SIZE = 6
local MIN_LEVEL = 35
local MAX_LEVEL = 250
local SCALE_LEVEL = 50
local REPEAT_QUEST = 199600
local PASS_ID = 4001022
local KEY_ID = 4001023

local function try_start(me, npc)
	pq.remove_all(PASS_ID, me)
	pq.remove_all(KEY_ID, me)
	local party = me:party()
	if party == nil then
		me:dialog(npc, "다른 차원에서 온 괴물이 이 안에 숨어들어 시간을 어지렵히고 있어요. 이 위에서 생긴 시공의 균열을 닫고 사악한 몬스터를 없애주실 분들을 찾고 있어요. 한번 도전해 보시지 않겠어요? 그렇다면 파티를 만들고 파티원을 모아보세요. 또는 다른 파티에 참여해보세요!")
		return
	end
	if party:leader_id() ~= me:id() then
		me:dialog(npc, "시공의 균열을 닫고 사악한 몬스터를 없애주시고 싶으시다구요? 그렇다면 당신의 파티의 #b파티장#k에게 퀘스트를 시작할 것을 요청하세요.")
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
	local members = party:members()
	local ok = true
	local in_map = 0
	local count = 0
	for _, mem in ipairs(members) do
		if mem ~= nil then
			count = count + 1
			if mem:level() < MIN_LEVEL or mem:level() > MAX_LEVEL then
				ok = false
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
	if count > MAX_PARTY_SIZE or in_map < MIN_PARTY_SIZE then
		ok = false
	end
	if not ok then
		me:dialog(npc, "파티원 전원이 레벨 "
			.. MIN_LEVEL
			.. " 이상이어야 하며, "
			.. MIN_PARTY_SIZE
			.. "~"
			.. MAX_PARTY_SIZE
			.. "명이 같은 맵에 모여 있어야 합니다.")
		return
	end
	local group = state_machine(GROUP_NAME)
	if group == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		return
	end
	local state = group:get_property("state")
	if state ~= nil and state ~= "" and state ~= "0" then
		me:dialog(npc, "이 안에 이미 다른 파티가 입장하여 클리어에 도전중입니다. 잠시 후에 다시 시도해보세요.")
		return
	end
	local sm, err = group:start_party(me, party, SCALE_LEVEL)
	if sm == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		if err ~= nil then
			log("ludibrium_party_quest start_party:", err)
		end
		return
	end
	pq.remove_all(PASS_ID, me)
	pq.remove_all(KEY_ID, me)
end

local function claim_earring(me, npc)
	local q = me:quest(REPEAT_QUEST)
	local count = 0
	if q ~= nil then
		count = tonumber(q:record()) or 0
	end
	if count >= 5 then
		local item = pq.gain_item(me, 1032064, 1)
		if item == nil then
			me:dialog(npc, "인벤토리 공간을 확보하신 후 다시 말을 걸어주세요.")
			return
		end
		if q:started() then
			q:record("0")
		else
			q:start("0")
		end
		me:dialog(npc, "그동안 도와주셔서 감사합니다. 총 5번 도와 주셔서 #b귀걸이#k를 1번 받으실 수 있습니다. #b귀걸이#k을 한개 드렸습니다. 앞으로 5번을 더 하시면 #b귀걸이#k 한개를 더 받으실 수 있습니다.")
		return
	end
	me:dialog(npc, "저를 5번 도와주실 때 마다 #v1032064# #b#z1032064##k을 1개씩 드리고 있습니다. 현재 도전 횟수는 #b"
		.. tostring(count)
		.. "#k회이며, 저를 5번 도와주시면 #b귀걸이#k을 받으실 수 있습니다.")
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "#e<파티퀘스트 : 차원의 균열>#n\r\n이 위부터는 엄청나게 위험한 존재들로 가득 차 있어 더이상 올라 가실 수 없어요. 파티원들과 함께 힘을 모아 퀘스트를 해결해 보시지 않겠습니까? 도전해 보고 싶다면 #b당신이 속한 파티의 파티장#k을 통해 저에게 말을 걸어 주세요.", {
			"함께 할 파티원을 찾고 싶어요",
			"시크릿 뉴비 귀걸이가 받고 싶어요",
			"설명을 듣고 싶어요",
			"파티퀘스트에 참가하고 싶습니다.",
		})
		if sel == 1 then
			me:dialog(npc, "현재 미구현된 컨텐츠입니다.\r\n\r\n컨텐츠 오픈일 : #b하루속히 빨리 오픈하겠습니다.#k\r\n\r\n파티를 구하고 싶으시면 #b파티원 찾기#k 버튼을 누르시면 됩니다.")
		elseif sel == 2 then
			claim_earring(me, npc)
		elseif sel == 3 then
			me:dialog(npc, "#e<파티퀘스트 : 차원의 균열>#n\r\n#b루디브리엄#k에 차원의 균열이 생겨났습니다! 이곳으로부터 침입한 몬스터들을 막아내려면 용감한 모험가들의 자발적인 도움이 절실해요. 믿을 수 있는 동료들과 힘을 합하여 루디브리엄을 구해주세요! 몬스터를 퇴치하거나 퀴즈를 풀어나가는 각종 난관을 해결하고 #r알리샤르#k에게 승리해야 한답니다.\r\n - #e레벨#n : 30이상 #r(추천레벨 : 35 ~ 50)#k\r\n - #e제한시간#n : 60분\r\n - #e참가인원#n : 3~6명\r\n - #e획득 아이템#n : #v1022073# #z1022073##k #b(15회 도와줄 때마다 획득)#k")
		elseif sel == 4 then
			try_start(me, npc)
		end
	end
}
