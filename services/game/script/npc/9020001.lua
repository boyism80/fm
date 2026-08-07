-- NPC name (String.wz/Npc.img.xml): 클로토

local pq = require("script/lib/party_quest")

local EXIT_MAP = 103000890
local BONUS_MAP = 103000805
local COUPON_ID = 4001007
local PASS_ID = 4001008
local RANKING_QUEST = 1201

local stage1_quests = {
	{
		question = "문제입니다. #r메이플스토리에서 전사 1차 전직을 하기 위한 최소 레벨#k 만큼만 쿠폰을 모아오세요.",
		answer = 10,
	},
	{
		question = "문제입니다. #r메이플스토리에서 전사 1차 전직을 하기 위한 최소 힘#k 만큼만 쿠폰을 모아오세요.",
		answer = 35,
	},
	{
		question = "문제입니다. #r메이플스토리에서 마법사 1차 전직을 하기 위한 최소 지력#k 만큼만 쿠폰을 모아오세요.",
		answer = 20,
	},
	{
		question = "문제입니다. #r메이플스토리에서 궁수 1차 전직을 하기 위한 최소 민첩성#k 만큼만 쿠폰을 모아오세요.",
		answer = 25,
	},
	{
		question = "문제입니다. #r메이플스토리에서 도적 1차 전직을 하기 위한 최소 민첩성#k 만큼만 쿠폰을 모아오세요.",
		answer = 25,
	},
	{
		question = "문제입니다. #r메이플스토리에서 2차 전직을 하기 위한 최소 레벨#k 만큼만 쿠폰을 모아오세요.",
		answer = 30,
	},
	{
		question = "문제입니다. #r메이플스토리에서 궁수 1차 전직을 하기 위한 최소 레벨#k 만큼만 쿠폰을 모아오세요.",
		answer = 10,
	},
	{
		question = "문제입니다. #r메이플스토리에서 레벨1에서 레벨2가 되기 위해 필요한 경험치량#k 만큼만 쿠폰을 모아오세요.",
		answer = 15,
	},
	{
		question = "문제입니다. #r메이플스토리에서 도적 1차 전직을 하기 위한 최소 레벨#k 만큼만 쿠폰을 모아오세요.",
		answer = 10,
	},
	{
		question = "문제입니다. #r메이플스토리에서 마법사 1차 전직을 하기 위한 최소 레벨#k 만큼만 쿠폰을 모아오세요.",
		answer = 8,
	},
}

local bonus_rewards = {
	{ id = 2000004, n = 5 },
	{ id = 2000001, n = 100 },
	{ id = 2000002, n = 70 },
	{ id = 2000003, n = 100 },
	{ id = 2000006, n = 50 },
	{ id = 2022000, n = 15 },
	{ id = 2022003, n = 15 },
	{ id = 2040002, n = 1 },
	{ id = 2040402, n = 1 },
	{ id = 2040502, n = 1 },
	{ id = 2040505, n = 1 },
	{ id = 2040602, n = 1 },
	{ id = 2040802, n = 1 },
	{ id = 4003000, n = 30 },
	{ id = 4010000, n = 8 },
	{ id = 4010001, n = 8 },
	{ id = 4010002, n = 8 },
	{ id = 4010003, n = 8 },
	{ id = 4010004, n = 8 },
	{ id = 4010005, n = 8 },
	{ id = 4010006, n = 5 },
	{ id = 4020000, n = 8 },
	{ id = 4020001, n = 8 },
	{ id = 4020002, n = 8 },
	{ id = 4020003, n = 8 },
	{ id = 4020004, n = 8 },
	{ id = 4020005, n = 8 },
	{ id = 4020006, n = 8 },
	{ id = 4020007, n = 3 },
	{ id = 4020008, n = 3 },
	{ id = 1032002, n = 1 },
	{ id = 1032004, n = 1 },
	{ id = 1032005, n = 1 },
	{ id = 1032006, n = 1 },
	{ id = 1032007, n = 1 },
	{ id = 1032009, n = 1 },
	{ id = 1032010, n = 1 },
	{ id = 1002026, n = 1 },
	{ id = 1002089, n = 1 },
	{ id = 1002090, n = 1 },
}

local function area_pattern(map, count)
	local pos = ""
	local total = 0
	for i = 0, count - 1 do
		local n = map:players_in_area(i)
		total = total + n
		pos = pos .. tostring(n)
	end
	return pos, total
end

local function handle_rope_stage(me, npc, sm, map, prop_key, template, area_count, exp)
	local party = me:party()
	local is_leader = party ~= nil and party:leader_id() == me:id()
	if not is_leader then
		me:dialog(npc, template)
		return
	end
	if sm:get_property(prop_key) == "" then
		me:dialog(npc, template)
		sm:set_property(prop_key, pq.shuffle(string.rep("1", 3) .. string.rep("0", area_count - 3)))
		return
	end
	local answer = sm:get_property(prop_key)
	local pos, total = area_pattern(map, area_count)
	if total ~= 3 and not pq.is_gm(me) then
		me:dialog(npc, "아직 3개의 정답을 찾지 못하신것 같군요. 끝부분이 아니라 가운데에 정확히 서 계셔야 정답 여부가 확인됩니다.")
		return
	end
	if pq.is_gm(me) then
		pos = answer
	end
	if answer == pos then
		map:clear_effect()
		pq.party_exp(sm, exp)
		local stage = tonumber(sm:get_property("stage")) or 1
		sm:set_property("stage", tostring(stage + 1))
		me:dialog(npc, "다음 스테이지로 통하는 포탈이 열렸습니다. 서둘러 주세요.")
	else
		map:show_effect("quest/party/wrong_kor")
		map:play_sound("Party1/Failed")
	end
end

local function give_bonus(me, npc)
	local reward = bonus_rewards[math.random(1, #bonus_rewards)]
	local item = pq.gain_item(me, reward.id, reward.n)
	if item == nil then
		me:dialog(npc, "인벤토리 공간을 확보하신 후 다시 말을 걸어주세요.")
		return
	end
	me:map(BONUS_MAP)
end

local function handle_stage1(me, npc, sm, map)
	local party = me:party()
	local is_leader = party ~= nil and party:leader_id() == me:id()
	if is_leader then
		if sm:get_property("stage1p") == "" then
			me:dialog(npc, "안녕하세요. 첫번째 스테이지에 오신 것을 환영합니다. 주변을 둘러보면 리게이터가 돌아다니고 있는 것을 볼 수 있을 겁니다. 리게이터는 쓰러뜨리면 꼭 한개의 쿠폰을 떨어뜨립니다. 파티장을 제외한 파티원 전원은 각각 저에게 말을 걸어 문제를 받고 문제의 답에 해당하는 수 만큼 리게이터가 주는 쿠폰을 모아와야 합니다. \r\n만일 정답만큼 쿠폰을 모아왔다면 저는 그 파티원에게 #b통행권#k을 드리게 됩니다. 파티장을 제외한 모든 파티원이 통행권을 얻어 파티장에게 넘겨주면 파티장이 그렇게 모은 #b통행권#k을 저에게 넘겨줌으로써 스테이지를 클리어 하게 됩니다. 되도록 빨리 해결해야 더 많은 스테이지에 도전할 수 있으므로 서둘러 주세요. 그럼 행운을 빕니다.")
			sm:set_property("stage1p", "1")
			return
		end
		local psize = #sm:players() - 1
		if pq.has_item(me, PASS_ID, psize) then
			map:clear_effect()
			pq.party_exp(sm, 1500)
			pq.remove_all(PASS_ID, me)
			local stage = tonumber(sm:get_property("stage")) or 1
			sm:set_property("stage", tostring(stage + 1))
			me:dialog(npc, "다음 스테이지로 통하는 포탈이 열렸습니다. 서둘러 주세요.")
		else
			me:dialog(npc, "죄송합니다. 통행증의 개수가 부족합니다. 파티장을 제외한 파티원의 인원수 만큼의 통행증이 필요합니다. 문제를 해결하고, 얻는 통행증을 제게 주시기 바랍니다.")
		end
		return
	end

	local key = "stage1_" .. tostring(me:id())
	local val = sm:get_property(key)
	if val == "" then
		local idx = math.random(0, #stage1_quests - 1)
		sm:set_property(key, tostring(idx))
		val = tostring(idx)
		if not me:dialog(npc, "안녕하세요. 첫번째 스테이지에 오신 것을 환영합니다. 주변을 둘러보면 리게이터가 돌아다니고 있는 것을 볼 수 있을 겁니다. 리게이터는 쓰러뜨리면 꼭 한개의 쿠폰을 떨어뜨립니다. 파티장을 제외한 파티원 전원은 각각 저에게 말을 걸어 문제를 받고 문제의 답에 해당하는 수 만큼 리게이터가 주는 쿠폰을 모아와야 합니다. \r\n만일 정답만큼 쿠폰을 모아왔다면 저는 그 파티원에게 #b통행권#k을 드리게 됩니다. 파티장을 제외한 모든 파티원이 통행권을 얻어 파티장에게 넘겨주면 파티장이 그렇게 모은 #b통행권#k을 저에게 넘겨줌으로써 스테이지를 클리어 하게 됩니다. 되도록 빨리 해결해야 더 많은 스테이지에 도전할 수 있으므로 서둘러 주세요. 그럼 행운을 빕니다.", false, true) then
			return
		end
		if not me:dialog(npc, "여기, 리게이터를 잡고 개인적으로 제가 내 드리는 문제의 답 만큼의 #b쿠폰#k을 가져오시면 통행증으로 바꿔드립니다. 그것을 모아서 제게 가져다 주시기 바랍니다.", false, true) then
			return
		end
		local q = stage1_quests[tonumber(val) + 1]
		me:dialog(npc, q.question)
		return
	end
	local idx = tonumber(val)
	if idx == nil then
		return
	end
	if idx >= #stage1_quests then
		me:dialog(npc, "제가 드리는 문제를 훌륭히 완수하셨습니다. 다른 파티원이 문제를 해결할 때 까지 잠시만 기다려 주세요.")
		return
	end
	local q = stage1_quests[idx + 1]
	if pq.item_count(me, COUPON_ID) == q.answer then
		local item = pq.gain_item(me, PASS_ID, 1)
		if item == nil then
			me:dialog(npc, "인벤토리 공간이 부족합니다. 인벤토리 공간의 여유를 다시 확인해보세요.")
			return
		end
		sm:set_property(key, tostring(#stage1_quests + 1))
		pq.remove_all(COUPON_ID, me)
		me:dialog(npc, "정답을 맞추셨습니다! 이 #b통행증#k을 파티장에게 건네주세요.")
	else
		me:dialog(npc, "정답이 아닙니다.\r\n\r\n" .. q.question)
	end
end

local function handle_stage5(me, npc, sm, map)
	local party = me:party()
	if party == nil or party:leader_id() ~= me:id() then
		me:dialog(npc, "마지막 스테이지에 오신것을 환영합니다. 주위를 둘러보시면 몬스터가 보일겁니다. 그 몬스터를 잡고, 나온 통행증을 제게 가져다 주시기 바랍니다. 파티원들이 통행증을 먹으면, 파티장은 그것을 모아 제게 넘겨주시면 됩니다. 이 몬스터들은 친숙한 몬스터 일수도 있지만, 훨씬 더 강력하므로 조심해 주시기 바랍니다. 그리고 맨 아래층에는, 무시무시한 킹슬라임이 기다리고 있습니다. 부디 행운을 빕니다!")
		return
	end
	if sm:get_property("stage4p") == "" then
		sm:set_property("stage4p", "1")
		me:dialog(npc, "마지막 스테이지에 오신것을 환영합니다. 주위를 둘러보시면 몬스터가 보일겁니다. 그 몬스터를 잡고, 나온 통행증을 제게 가져다 주시기 바랍니다. 파티원들이 통행증을 먹으면, 파티장은 그것을 모아 제게 넘겨주시면 됩니다. 이 몬스터들은 친숙한 몬스터 일수도 있지만, 훨씬 더 강력하므로 조심해 주시기 바랍니다. 그리고 맨 아래층에는, 무시무시한 킹슬라임이 기다리고 있습니다. 부디 행운을 빕니다!")
		return
	end
	if not pq.has_item(me, PASS_ID, 10) then
		me:dialog(npc, "아직 통행증 10장을 모으지 못하신 것 같군요. 몬스터들을 잡고 제게 통행증 10장을 가져다 주세요.")
		return
	end
	map:clear_effect()
	pq.party_exp(sm, 3500)
	local stage = tonumber(sm:get_property("stage")) or 5
	sm:set_property("stage", tostring(stage + 1))
	pq.remove_all(PASS_ID, me)
	for _, p in ipairs(sm:players()) do
		if p ~= nil then
			p:end_party_quest(RANKING_QUEST)
		end
	end
	me:dialog(npc, "모든 문제를 훌륭히 해결하셨습니다. 모든 스테이지를 클리어 하셨으므로 보너스 스테이지로 이동됩니다. 남은 시간동안 마음껏 사냥하실 수 있습니다. 하지만 중간에 나가고 싶으시면 NPC를 통해 밖으로 나가실 수 있습니다. 저에게 다시 말을 걸어 주시면 작은 보상을 드리도록 할게요.")
end

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			me:map(EXIT_MAP)
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
		local field_id = wz.id
		if sm:get_property("stage") == "" then
			sm:set_property("stage", "1")
		end
		local curstage = tonumber(sm:get_property("stage")) or 1
		local map_stage = (field_id % 10) + 1

		if curstage > map_stage and field_id ~= 103000804 then
			me:dialog(npc, "다음 스테이지로 통하는 포탈이 열렸습니다. 서둘러 주세요.")
			return
		end
		if field_id == 103000804 and curstage > map_stage then
			give_bonus(me, npc)
			return
		end

		if field_id == 103000800 then
			handle_stage1(me, npc, sm, map)
		elseif field_id == 103000801 then
			handle_rope_stage(
				me,
				npc,
				sm,
				map,
				"stage2r",
				"두번째 스테이지에 대해 설명해 드리겠습니다. 옆에 밧줄들이 보일 것입니다. 이 밧줄들 중에서 #b3개가 다음 스테이지로 향하는 포탈#k과 통해 있습니다. 파티원 중에서 #b3 명이 정답 줄을 찾아 매달리면 됩니다.#k\r\n단, 줄 끝에 아슬아슬하게 매달리시지 말고 줄 가운데에 매달려 계셔야 정답으로 인정되니 이점 주의해 주시기 바랍니다. 그리고 반드시 3 명만 줄에 매달려 계셔야 있어야 합니다. 파티원이 줄에 올라서면 파티장은 #b저를 더블클릭하여 정답인지 아닌지 확인#k해야 합니다. 그럼 힘내 주세요!",
				4,
				1200
			)
		elseif field_id == 103000802 then
			handle_rope_stage(
				me,
				npc,
				sm,
				map,
				"stage3r",
				"세번째 스테이지에 대해 설명해 드리겠습니다. 옆에 나무 발판들이 보일 것입니다. 이 발판들 중에서 #b3개가 다음 스테이지로 향하는 포탈#k과 통해 있습니다. 파티원 중에서 #b3 명이 정답 발판을 찾아 올라서면 됩니다.#k\r\n단, 발판 끝에 아슬아슬하게 서계시지 말고 발판 가운데에 정확하게 서 계셔야 정답으로 인정되니 이점 주의해 주시기 바랍니다. 그리고 반드시 3명만 발판에 서계셔야 있어야 합니다. 파티원이 발판에 올라서면 파티장은 #b저를 더블클릭하여 정답인지 아닌지 확인#k해야 합니다. 그럼 힘내 주세요!",
				5,
				1400
			)
		elseif field_id == 103000803 then
			handle_rope_stage(
				me,
				npc,
				sm,
				map,
				"stage4r",
				"네번째 스테이지에 대해 설명해 드리겠습니다. 옆에 나무통들이 보일 것입니다. 이 통들 중에서 #b3개가 다음 스테이지로 향하는 포탈#k과 통해 있습니다. 파티원 중에서 #b3 명이 정답 통을 찾아 올라서면 됩니다.#k\r\n단, 통 끝에 아슬아슬하게 서계시지 말고 통 가운데에 정확하게 서 계셔야 정답으로 인정되니 이점 주의해 주시기 바랍니다. 그리고 반드시 3명만 통 위에 서계셔야 있어야 합니다. 파티원이 통 위에 올라서면 파티장은 #b저를 더블클릭하여 정답인지 아닌지 확인#k해야 합니다. 그럼 힘내 주세요!",
				6,
				1800
			)
		elseif field_id == 103000804 then
			handle_stage5(me, npc, sm, map)
		end
	end
}
