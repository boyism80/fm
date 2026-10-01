-- NPC name (String.wz/Npc.img.xml): 카이린

local MAIN_MAP = 120000101
local CLEAR_MAP = 912010200
local BLACK_CHARM = 4031059
local NECKLACE = 4031057
local TRIAL_QUEST = 195000
local DEFAULT_TEXT = "신속한 너클과 강력한 화력의 건을 사용하는 특수하고 조직적인 직업. 혹시 해적에 관심이 있으신 건가요?"

local GREETINGS = {
	[Class.Pirate] = "당신은 얼마 전에 #b해적#k이 되었던 #r#h ##k이군요, 수행은 잘 되시는지요? 더 현명한 해적이 되기 위하여 더 정진하길 바랍니다..",
	[Class.Brawler] = "당신은 얼마 전에 #b인파이터#k가 되었던 #r#h ##k이군요, 수행은 잘 되시는지요? 더 현명한 해적이 되기 위하여 더 정진하길 바랍니다..",
	[Class.Gunslinger] = "당신은 얼마 전에 #b건슬링거#k가 되었던 #r#h ##k이군요, 수행은 잘 되시는지요? 더 현명한 해적이 되기 위하여 더 정진하길 바랍니다..",
	[Class.Marauder] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b버커니어#k로 전직하였던가. #r#h ##k여 버커니어로 전직한 것을 축하하네, 더 강한 해적이 되기 위하여 더 정진하길 바라네..",
	[Class.Outlaw] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b발키리#k로 전직하였던가. #r#h ##k여 발키리로 전직한 것을 축하하네, 더 강한 해적이 되기 위하여 더 정진하길 바라네..",
	[Class.Buccaneer] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b바이퍼#k로 전직하였던가. #r#h ##k여 바이퍼로 전직한 것을 축하하네, 더 강한 해적이 되기 위하여 더 정진하길 바라네..",
	[Class.Corsair] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b캡틴#k으로 전직하였던가. #r#h ##k여 캡틴으로 전직한 것을 축하하네, 더 강한 해적이 되기 위하여 더 정진하길 바라네..",
}

local FIRST_JOB_WEAPONS = {
	{
		item = { [1482014] = 1 },
		stats = { 25, 4, 4, 4 },
	},
	{
		item = { [1492014] = 1, [2330006] = 600 },
		stats = { 4, 25, 4, 4 },
	},
}

local SECOND_JOBS = {
	[2191] = {
		class = Class.Brawler,
		text = "좋아. 너는 이제부터 #b인파이터#k야. 인파이터는 맨 몸과 맨주먹 만으로 적을 제압하는 자... 그렇기 때문에 남들보다 수련에 더 힘써야하지. 수련에 어려움이 있다면 나도 도와주겠어. 그리고 너에게 인파이터가 익힐 수 있는 스킬들이 적혀있는 책을 주었어... 그 책에는 여러가지 인파이터에 관련된 스킬들이 들어 있어. 또한 너에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어봐. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있으니까, 명심해 둬. 이제 인파이터로써 너는 더 한층 높은 해적이 되었어. 열심히 수련해 주길 바래. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록해.",
	},
	[2192] = {
		class = Class.Gunslinger,
		text = "좋아. 너는 이제부터 #b건슬링거#k야. 건슬링거는 건으로 멀리있는 적에게까지 예리한 공격으로 위협할 수 있는 자... 그렇기 때문에 남들보다 수련에 더 힘써야하지. 수련에 어려움이 있다면 나도 도와주겠어. 그리고 너에게 건슬링거가 익힐 수 있는 스킬들이 적혀있는 책을 주었어... 그 책에는 여러가지 건슬링거에 관련된 스킬들이 들어 있어. 또한 너에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어봐. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있으니까, 명심해 둬. 이제 건슬링거로써 너는 더 한층 높은 해적이 되었어. 열심히 수련해 주길 바래. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록해.",
	},
}

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

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

local function reset_stats(me, str, dex, int, luk)
	local gained = (me:base_str() - str) + (me:base_dex() - dex) + (me:base_int() - int) + (me:base_luk() - luk)
	me:ability_point(math.max(0, me:ability_point() + gained))
	me:base_str(str)
	me:base_dex(dex)
	me:base_int(int)
	me:base_luk(luk)
end

local function is_second_job(me)
	local class = me:class()
	return class == Class.Brawler or class == Class.Gunslinger
end

local function first_job(me, npc)
	if me:dialog(npc, "해적이 되고 싶어서 저를 찾아 온 건가요? 그렇다면 제대로 찾아 왔어요.", false, true) == false then
		return
	end
	if me:dialog(npc, "해적이 되고 싶은가요? 하지만 조건이 필요한데... #b레벨이 10 이상#k이어야 한답니다. 어디볼까요... 흐음...", false, true) == false then
		return
	end
	if me:level() < 10 then
		me:dialog(npc, "당신은 아직 수련이 더 필요한 몸인 것 같네요, 좀 더 수련을 한 뒤에 찾아와 주세요.")
		return
	end
	if me:dialog_yes_no(npc, "당신은 자격이 있어 보이는 군요. 신속한 너클과 강력한 화력의 건을 사용할 수 있는... 그런 분이 필요했어요. ") == false then
		me:dialog(npc, "그런가요? 신중하게 생각해 보시고 결정하세요.")
		return
	end
	if me:dialog(npc, "좋습니다! 당신은 이제부터 해적입니다! 제가 직접 인정했으니까요... 작지만 당신에게 제가 가진 능력의 일부를 조금 보태드리도록 하겠습니다. 이야~~~ 압!!!", false, true) == false then
		return
	end
	local sel = me:dialog_list(npc, "두 가지 해적 중 어느 해적을 지향하세요?\r\n\r\n#b", {
		"너클을 사용하는 인파이터 계열",
		"건을 사용하는 건슬링거 계열",
	})
	if sel == nil then
		return
	end
	local weapon = FIRST_JOB_WEAPONS[sel]
	if me:exchange(nil, { item = weapon.item }) ~= ExchangeResult.OK then
		me:dialog(npc, "뭘 그렇게 많이 가지고 다니는거야? 장비와 소비 인벤토리를 비우고 다시 와.")
		return
	end
	me:class(Class.Pirate)
	reset_stats(me, weapon.stats[1], weapon.stats[2], weapon.stats[3], weapon.stats[4])
	if me:dialog(npc, "당신에게 약간의 #bSP#k를 드렸습니다. 왼쪽 하 단에 있는 #bSkill 메뉴#k를 열어보세요. 스킬을 올릴 수 있을 꺼에요. 단 처음부터 전부 올릴수 있는 건 아니에요... 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있어요.", false, true) == false then
		return
	end
	if me:dialog(npc, "한 가지 더 주의해야 할 점이 있어요. 초보자에서 직업을 가진 그 순간부터는 죽지 않도록 조심해야 한답니다. 만일 죽게 되면 그동안 쌓였던 경험치가 깍일 수 있으니까요.", false, true) == false then
		return
	end
	me:dialog(npc, "제가 가르쳐 드릴 수 있는건 여기까지 입니다. 이곳 저곳 여행을 하면서 자기 자신을 단련시키는 일만이 남았어요. 자신이 더욱 강해졌다고 생각되면 다시 절 찾아오세요. 당신을 계속 기다리고 있겠습니다.")
end

local function second_job(me, npc)
	for quest_id, job in pairs(SECOND_JOBS) do
		if me:quest(quest_id):completed() then
			if me:skill_point() > (me:level() - 30) * 3 then
				me:dialog(npc, "음.. SP를 아직 다 사용하지 않은 것 같은데? SP가 너무 많아 남아 아직 2차전직을 할 수 없어.")
				return
			end
			me:class(job.class)
			me:dialog(npc, job.text)
			return
		end
	end
	me:dialog(npc, "해적이 되고 싶은 자는 나에게...")
end

local function third_job_trial(me, npc, value)
	if value == "job3_trial1_1" then
		if is_second_job(me) == false or me:level() < 70 then
			me:dialog(npc, DEFAULT_TEXT)
			return
		end
		me:quest(TRIAL_QUEST):record("job3_trial1_2")
		me:dialog(npc, "널 기다리고 있었어. 몇 일 전 오시리아 대륙의 #b페드로#k님으로부터 너에 대한 이야기를 전해 들었어. 좋아... 너의 힘을 내가 시험해 줄게. 빅토리아 아일랜드 깊은숲 어딘가에 다른 차원으로 통하는 균열이 있어. 보통 사람들은 들어갈 수 없지만 너에 한해서는 들어갈 수 있도록 해 두겠어. 균열 안으로 들어가면 나의 분신을 만날 수 있는데 그 분신을 쓰러뜨리고 #b검은 부적#k을 얻어 나에게 가져와 줘.")
		return
	end
	if value == "job3_trial1_2" then
		if item_count(me, BLACK_CHARM) >= 1 then
			if me:exchange({ item = { [BLACK_CHARM] = 1 } }, { item = { [NECKLACE] = 1 } }) ~= ExchangeResult.OK then
				me:dialog(npc, "흐음. 인벤토리 공간이 부족한 것 같은데. 기타 탭을 충분히 비우고 다시 와.")
				return
			end
			me:quest(TRIAL_QUEST):record("job3_trial1_3")
			me:dialog(npc, "이럴수가... 제 분신을 쓰러뜨리고 #b검은 부적#k을 가져오셨군요! 좋아요... 이걸로 당신의 힘은 충분히 증명되었습니다. 힘에 한해서는 #b3차 전직#k을 하기에 부족함이 없어 보이는군요. 약속대로 당신에게 #b강인함의 목걸이#k를 드리겠습니다. 이 목걸이를 가지고 오시리아의 #b페드로#k님에게 돌아가면 두번째 시험을 치를 수 있겠지요. 그럼 당신이 무사히 3차 전직을 할 수 있기를 빌겠습니다.")
			return
		end
		if is_second_job(me) == false or me:level() < 70 then
			me:dialog(npc, DEFAULT_TEXT)
			return
		end
		me:dialog(npc, "빅토리아 아일랜드 어딘가의 다른 차원으로 통하는 균열로 가서, 내 분신을 쓰러뜨리고 #b검은 부적#k을 얻어 내게 가져오도록 해.")
		return
	end
	me:dialog(npc, "제가 당신에게 주었던 그 목걸이를 가지고 오시리아의 #b페드로#k님에게 돌아가면 두번째 시험을 치를 수 있겠지요. 그럼 당신이 무사히 3차 전직을 할 수 있기를 빌겠습니다.")
end

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
		if wz ~= nil and wz:id() == CLEAR_MAP then
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
		if me:class() == Class.Beginner then
			first_job(me, npc)
			return
		end
		if me:class() == Class.Pirate and me:level() >= 30 then
			second_job(me, npc)
			return
		end

		local q = me:quest(TRIAL_QUEST)
		if q:started() == false then
			q:start("0")
		end
		local value = q:record()
		if string.sub(value, 1, 11) == "job3_trial1" then
			third_job_trial(me, npc, value)
			return
		end

		if handle_training(me, npc) then
			return
		end
		me:dialog(npc, GREETINGS[me:class()] or DEFAULT_TEXT)
	end
}
