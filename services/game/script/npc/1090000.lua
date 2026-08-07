-- NPC name (String.wz/Npc.img.xml): 카이린

local MAIN_MAP = 120000101
local CLEAR_MAP = 912010200

local TEST_GROUPS = {
	[2191] = {
		"kyrin_test_1_1",
		"kyrin_test_1_2",
	},
	[2192] = {
		"kyrin_test_2_1",
		"kyrin_test_2_2",
	},
}

local TEST_DIALOG = {
	[2191] = "시험의 장소에서 파이렛 옥토를 해치우고 #b#t4031856##k 15개를 모아오겠어?",
	[2192] = "시험의 장소에서 파이렛 옥토를 해치우고 #b#t4031857##k 15개를 모아오겠어?",
}

local TRAINING_GROUPS = {
	[6370] = {
		group = "kyrin_training_ground_c",
		progress_quest = 6371,
	},
	[6330] = {
		group = "kyrin_training_ground_v",
		progress_quest = 6331,
	},
}

local function start_solo(group, me)
	local sm, err = group:start_solo(me)
	if sm ~= nil then
		return true
	end
	if err ~= nil then
		log(group:name() .. " start_solo:", err)
	end
	return false
end

local function start_test(me, npc, quest_id)
	for _, group_name in ipairs(TEST_GROUPS[quest_id]) do
		local group = state_machine(group_name)
		if group ~= nil and group:get_property("state") ~= "1" then
			if start_solo(group, me) then
				return
			end
		end
	end
	me:dialog(npc, "이미 시험장이 모두 사용 중인걸? 나중에 다시 시도해 봐.")
end

local function set_training_progress(me, npc, quest_id)
	local quest = me:quest(quest_id)
	if quest:started() then
		quest:record("2")
		return
	end
	if quest:wz() == nil then
		quest:start("2")
	else
		quest:start(npc, "2")
	end
end

local function handle_clear(me, npc)
	for quest_id, config in pairs(TRAINING_GROUPS) do
		local quest = me:quest(quest_id)
		if quest:started() then
			me:dialog(npc, "호오? 역시 기대한 대로 대단한데? 좋아, 나를 따라 나오도록 해. 네게 걸맞은 선물을 주지.")
			set_training_progress(me, npc, config.progress_quest)
			me:show_quest_completion(quest_id)
			me:map(MAIN_MAP)
			return
		end
	end
	me:map(MAIN_MAP)
end

local function handle_training(me, npc)
	for quest_id, config in pairs(TRAINING_GROUPS) do
		local quest = me:quest(quest_id)
		if quest:started() then
			local progress = me:quest(config.progress_quest)
			if progress:record() == "2" then
				me:dialog(npc, "이미 나와의 대련은 끝난 것 같은데? 더 이상 대련할 필요는 없어.")
				return true
			end
			if not me:dialog_yes_no(npc, "좋아, 내가 보내주는 수련장에서 나의 공격을 2분 이상 버텨내면 돼. 잘 해보라구.") then
				return true
			end
			local group = state_machine(config.group)
			if group == nil or group:get_property("started") == "true" or not start_solo(group, me) then
				me:dialog(npc, "지금은 훈련장을 사용할 수 없어. 잠시 후 다시 시도해 줘.")
			end
			return true
		end
	end
	return false
end

return {
	on_click = function(me, npc)
		local map = me:map()
		local wz = map ~= nil and map:wz() or nil
		if wz ~= nil and wz.id == CLEAR_MAP then
			handle_clear(me, npc)
			return
		end
		for quest_id, _ in pairs(TEST_GROUPS) do
			local quest = me:quest(quest_id)
			if quest:started() then
				if me:dialog_yes_no(npc, TEST_DIALOG[quest_id]) then
					start_test(me, npc, quest_id)
				end
				return
			end
		end
		if handle_training(me, npc) then
			return
		end
		me:dialog(npc, "신속한 너클과 강력한 화력의 건을 사용하는 해적에게 관심이 있는 건가요?")
	end
}
