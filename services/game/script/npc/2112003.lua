-- NPC name (String.wz/Npc.img.xml): 줄리엣

local pq = require("script/lib/party_quest")
local rj = require("script/lib/romeo_juliet")

local GROUP_NAME = "juliet_party_quest"
local ENTRY_MAP = 261000021

local function try_start(me, npc)
	rj.strip_items(me)
	local party = me:party()
	if party == nil or party:leader_id() ~= me:id() then
		me:dialog(npc, "파티가 없으시거나, 혹은 파티장이 아니신 건 아닌가요? 또한 4명의 파티만 입장할 수 있습니다.")
		return
	end
	local map = me:map()
	if map == nil or map:wz() == nil then
		return
	end
	local ok = true
	local size = 0
	local has_gm = false
	for _, mem in ipairs(party:members()) do
		if mem ~= nil then
			local ch = map:characters()[mem:id()]
			if ch == nil or mem:level() < rj.min_level then
				ok = false
				break
			end
			if pq.is_gm(ch) then
				has_gm = true
			end
			size = size + 1
		end
	end
	if not ok or (size ~= rj.min_party and not has_gm) then
		me:dialog(npc, "파티원 중 레벨 제한 71이상이 아닌 파티원이 있거나 파티원이 4명이 아닌 것 같군요.")
		return
	end
	local group = state_machine(GROUP_NAME)
	if group == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		return
	end
	local state = group:get_property("state")
	if state ~= nil and state ~= "" and state ~= "0" then
		me:dialog(npc, "이미 이 안에서 다른 파티가 퀘스트에 도전하는 중입니다. 잠시 후 다시 시도해 주세요.")
		return
	end
	local sm, err = group:start_party(me, party, rj.scale_level)
	if sm == nil then
		me:dialog(npc, "파티 퀘스트를 시작할 수 없습니다.")
		if err ~= nil then
			log("juliet_party_quest start_party:", err)
		end
	end
end

local function handle_hub(me, npc, sm, map)
	if pq.has_item(me, 4001130, 1) then
		pq.remove_all(4001130, me)
		sm:set_property("stage", "1")
		sm:notice("줄리엣은 로미오가 쓴 편지를 보고 생각에 잠겼다.")
		me:dialog(npc, "이.. 이것은 로미오가 쓴 편지..?")
		return
	end
	if not pq.is_leader(me) then
		me:dialog(npc, "알카드노와 제뉴미스트 사이의 분쟁을 멈추기 위해서는 알카드노와 제뉴미스트의 실험 자료가 필요해요. 알카드노의 실험 자료를 먼저 제게 가져다 주세요. 또한, 파티장이 제게 말을 걸어주셔야 합니다.")
		return
	end
	if pq.has_item(me, 4001134, 1) then
		pq.remove_all(4001134, me)
		sm:set_property("stage4", "1")
		me:dialog(npc, "알카드노 실험 자료를 찾아오셨군요! 이제 제뉴미스트 실험 자료를 찾으면 되겠군요.")
		return
	end
	if pq.has_item(me, 4001135, 1) and sm:get_property("stage4") == "1" then
		pq.remove_all(4001135, me)
		rj.clear_fx(map)
		map:block_gen(true)
		map:kill_all_mobs()
		pq.party_exp(sm, 10000)
		sm:set_property("stage4", "2")
		local door = map:reactor_by_name("jnr3_out3")
		if door ~= nil then
			door:hit(1)
		end
		me:dialog(npc, "제뉴미스트 실험 자료도 찾아오셨군요. 다음 스테이지로 진행해 주세요.")
		return
	end
	me:dialog(npc, "알카드노와 제뉴미스트 사이의 분쟁을 멈추기 위해서는 알카드노와 제뉴미스트의 실험 자료가 필요해요. 알카드노의 실험 자료를 먼저 제게 가져다 주세요.")
end

local function handle_reward(me, npc)
	me:dialog(npc, "이렇게 유레테도 마음을 고쳐먹었고, 제뉴미스트와 알카드노의 전쟁도 조금은 수그러 들 수 있을까요? 비록 저희 사랑이 이뤄지기까지 아직도 험난한 고난이 기다리고 있지만, 결코 포기하지 않고 로미오를 지키겠습니다.")
	local code = me:exchange({}, { item = { [4001160] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "인벤토리 공간이 부족하신 건 아닌지 확인해 주시겠어요?")
		return
	end
	if code ~= ExchangeResult.OK then
		me:dialog(npc, "인벤토리 공간이 부족하신 건 아닌지 확인해 주시겠어요?")
		return
	end
	rj.strip_items(me)
	me:end_party_quest(rj.ranking_quest)
	local q = me:quest(199603)
	if q ~= nil then
		local n = tonumber(q:record()) or 0
		q:record(tostring(n + 1))
	end
	me:map(926110700)
end

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz().id
		if map_id == ENTRY_MAP then
			local sel = me:dialog_list(npc, "#e<파티퀘스트 : 로미오와 줄리엣>#n\r\n마가티아는 지금 크나큰 위기를 맞이하고 있습니다. 용감한 메이플의 모험가들께서 저희를 도와주시지 않으시겠어요?\r\n\r\n#b", {
				"줄리엣의 이야기를 듣는다",
				"함께 할 파티를 찾고 싶어요.",
				"퀘스트를 시작한다.",
			})
			if sel == 1 then
				me:dialog(npc, "#e<파티퀘스트 : 로미오와 줄리엣>#n\r\n로미오와 줄리엣을 도와주세요.\r\n - #e레벨#n : 71이상 #r(추천레벨 : 71 ~ 85)#k\r\n - #e참가인원#n : 4명\r\n - #e최종 보상#n : #v1122010# #z1122010##k #b(알카드노, 제뉴미스트의 구슬 각각 25개로 교환)#k")
			elseif sel == 2 then
				me:dialog(npc, "현재 미구현된 컨텐츠입니다.\r\n\r\n컨텐츠 오픈일 : #b하루속히 빨리 오픈하겠습니다.#k\r\n\r\n파티를 구하고 싶으시면 #b파티원 찾기#k 버튼을 누르시면 됩니다.")
			elseif sel == 3 then
				try_start(me, npc)
			end
			return
		end
		local sm = me:state_machine()
		if map_id == 926110000 then
			me:dialog(npc, "이 연구실에서 가끔 수상한 소리가 들린다는 소문이 있었어요. 분명 이 근처 어딘가에 수상한 소리의 근원이 있을 거에요.")
		elseif map_id == 926110001 then
			me:dialog(npc, "모든 몬스터를 물리쳐 주세요!")
		elseif map_id == 926110100 then
			me:dialog(npc, "이곳의 비커들은 깨져서 내용물이 새고 있답니다. 비커에 내용물이 모두 새기 전에 빠르게 수상한 액체를 채워 넣어주세요.")
		elseif map_id == 926110200 then
			if sm == nil then
				return
			end
			handle_hub(me, npc, sm, map)
		elseif map_id == 926110300 then
			me:dialog(npc, "각 파티원이 각각 연구실을 통해 꼭대기에 있는 중앙 연구실 까지 올라가야 해요.")
		elseif map_id == 926110400 then
			me:dialog(npc, "준비가 되면, 제 사랑을 구하러 가요.")
		elseif map_id == 926110401 then
			if sm ~= nil then
				pq.party_warp(sm, 926110500)
			end
		elseif map_id == 926100401 then
			me:dialog(npc, "구해주셔서 진심으로 감사드립니다.")
		elseif map_id == 926110600 then
			handle_reward(me, npc)
		elseif map_id == 926100600 then
			me:dialog(npc, "저를 구해주신 모험가님. 진심으로 감사드립니다.")
		end
	end
}
