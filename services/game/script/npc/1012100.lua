-- NPC name (String.wz/Npc.img.xml): 헬레나

local LETTER = 4031010
local PROOF = 4031012
local BLACK_CHARM = 4031059
local NECKLACE = 4031057
local TRIAL_QUEST = 195000
local DEFAULT_TEXT = "밝은 눈으로 진정한 적을 가려내어 그 심장에 화살을 꼿을 수 있는 직업. 궁수는 정말 화려한 직업이지요. 혹시 궁수에 관심이 있으신 건가요?"

local GREETINGS = {
	[Class.Bowman] = "당신은 얼마 전에 #b궁수#k가 되었던 #r#h ##k이군요, 수행은 잘 되시는지요? 더 현명한 궁수가 되기 위하여 더 정진하길 바랍니다..",
	[Class.Hunter] = "당신은 얼마 전에 #b헌터#k가 되었던 #r#h ##k이군요, 수행은 잘 되시는지요? 더 현명한 궁수가 되기 위하여 더 정진하길 바랍니다..",
	[Class.Crossbowman] = "당신은 얼마 전에 #b사수#k가 되었던 #r#h ##k이군요, 수행은 잘 되시는지요? 더 현명한 궁수가 되기 위하여 더 정진하길 바랍니다..",
	[Class.Ranger] = "당신의 소식은 잘 알고 있습니다, 얼마 전에 #b레인저#k로 전직하셨다고 들었습니다. #r#h ##k여 레인저로 전직한 것을 축하합니다, 더 현명한 궁수가 되기 위하여 더 정진하길 바랍니다..",
	[Class.Sniper] = "당신의 소식은 잘 알고 있습니다, 얼마 전에 #b저격수#k로 전직하셨다고 들었습니다. #r#h ##k여 저격수로 전직한 것을 축하합니다, 더 현명한 궁수가 되기 위하여 더 정진하길 바랍니다..",
	[Class.Bowmaster] = "당신의 소식은 잘 알고 있습니다, 얼마 전에 #b보우마스터#k로 전직하셨다고 들었습니다. #r#h ##k여 보우마스터로 전직한 것을 축하합니다, 더 현명한 궁수가 되기 위하여 더 정진하길 바랍니다..",
	[Class.Crossbowmaster] = "당신의 소식은 잘 알고 있습니다, 얼마 전에 #b신궁#k으로 전직하셨다고 들었습니다. #r#h ##k여 신궁으로 전직한 것을 축하합니다, 더 현명한 궁수가 되기 위하여 더 정진하길 바랍니다..",
}

local SECOND_JOBS = {
	{
		class = Class.Hunter,
		info = "#b헌터#k은 활을 자유자재로 사용할 수 있도록 숙련된 궁수 클레스에요, #b소울 에로우#k라는 스킬을 이용하여 화살을 쓰지 않고도 활을 쏠 수 있게 하고, #b에로우 봄#k으로 폭탄이 들은 화살을 쏘아 폭발시켜 그 충격으로 몬스터를 기절시키는 기술, 그리고 #b부스터#k와 같은 스킬로 공격 속도를 한 단계 늘릴 수 있고, #b파이널 어택#k으로 피니쉬 어택이라는 강력한 한방을 날려줄 수 있답니다.",
		confirm = "#b헌터#k로 2차 전직하고 싶단 말씀이시죠? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없습니다. 그 결심... 틀림이 없겠죠?",
		done = "#b헌터#k입니다. 헌터는 밝은 눈으로 적의 가슴에 화살을 꼿을 수 있는 현명한 사람 ... 더욱 수련에 정진하시길 바라겠습니다. 저의 힘으로 당신을 더욱 강하게 만들어 드리겠습니다.그리고 당신에게 사수가 익힐 수 있는 스킬들이 적혀있는 책을 드렸습니다... 그 책에는 여러가지 사수와 관련된 스킬들이 들어 있습니다. 또한 당신에게 약간의 #bSP#k를 드렸으니 #bSkill 메뉴#k를 열어보세요. 스킬을 올릴 수 있을겁니다. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있으니까, 명심해 두세요. 이제 헌터로써 당신은 더 한층 높은 궁수가 되었어요. 열심히 수련해 주시길 바랍니다. 그리고 자신이 강하다고 생각할 때가 되면 저를 찾아오도록 바랍니다.",
	},
	{
		class = Class.Crossbowman,
		info = "#b사수#k는 석궁을 자유자재로 사용할 수 있도록 숙련된 궁수 클레스에요, #b소울 에로우#k라는 스킬을 이용하여 화살을 쓰지 않고도 활을 쏠 수 있게 하고, #b아이언에로우#k로 다수의 적을 관통하는 무겁고 단단한 강철 화살을 쏘는 기술, 그리고 #b부스터#k와 같은 스킬로 공격 속도를 한 단계 늘릴 수 있고, #b파이널 어택#k으로 피니쉬 어택이라는 강력한 한방을 날려줄 수 있답니다.",
		confirm = "#b사수#k로 2차 전직하고 싶단 말씀이시죠? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없습니다. 그 결심... 틀림이 없겠죠?",
		done = "#b사수#k입니다. 사수는 밝은 눈으로 적의 가슴에 화살을 꼿을 수 있는 현명한 사람 ... 더욱 수련에 정진하시길 바라겠습니다. 저의 힘으로 당신을 더욱 강하게 만들어 드리겠습니다. 그리고 당신에게 사수가 익힐 수 있는 스킬들이 적혀있는 책을 드렸습니다... 그 책에는 여러가지 사수와 관련된 스킬들이 들어 있습니다. 또한 당신에게 약간의 #bSP#k를 드렸으니 #bSkill 메뉴#k를 열어보세요. 스킬을 올릴 수 있을겁니다. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있으니까, 명심해 두세요. 이제 사수로써 당신은 더 한층 높은 궁수가 되었어요. 열심히 수련해 주시길 바랍니다. 그리고 자신이 강하다고 생각할 때가 되면 저를 찾아오도록 바랍니다.",
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
	return class == Class.Hunter or class == Class.Crossbowman
end

local function first_job(me, npc)
	if me:dialog(npc, "궁수가 되고 싶어서 저를 찾아 온 건가요? 그렇다면 제대로 찾아 왔어요.", false, true) == false then
		return
	end
	if me:dialog(npc, "궁수가 되고 싶은가요? 하지만 조건이 필요한데... #b레벨이 10 이상#k이어야 한답니다. 어디볼까요... 흐음...", false, true) == false then
		return
	end
	if me:level() < 10 then
		me:dialog(npc, "당신은 아직 수련이 더 필요한 몸인 것 같네요, 좀 더 수련을 한 뒤에 찾아와 주세요.")
		return
	end
	if me:dialog_yes_no(npc, "당신은 자격이 있어 보이는 군요. 밝은 눈으로 진정한 적을 가려내어 그 심장에 화살을 꼿을 수 있는... 그런 분이 필요했어요. 궁수로 전직하고 싶은가요?") == false then
		me:dialog(npc, "그런가요? 신중하게 생각해 보시고 결정하세요.")
		return
	end
	if me:dialog(npc, "좋습니다! 당신은 이제부터 궁수입니다! 제가 직접 인정했으니까요... 작지만 당신에게 제가 가진 능력의 일부를 조금 보태드리도록 하겠습니다. 이야~~~ 압!!!", false, true) == false then
		return
	end
	if me:empty_slots(InventoryType.Equipment) < 1 or me:empty_slots(InventoryType.Use) < 1 then
		me:dialog(npc, "장비와 소비 인벤토리를 비우고 다시 오세요.")
		return
	end
	if me:exchange(nil, { item = { [1452051] = 1, [2060000] = 2000 } }) ~= ExchangeResult.OK then
		me:dialog(npc, "장비와 소비 인벤토리를 비우고 다시 오세요.")
		return
	end
	me:class(Class.Bowman)
	reset_stats(me, 4, 25, 4, 4)
	if me:dialog(npc, "당신에게 약간의 #bSP#k를 드렸습니다. 왼쪽 하 단에 있는 #bSkill 메뉴#k를 열어보세요. 스킬을 올릴 수 있을 꺼에요. 단 처음부터 전부 올릴수 있는 건 아니에요... 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있어요.", false, true) == false then
		return
	end
	if me:dialog(npc, "한 가지 더 주의해야 할 점이 있어요. 초보자에서 직업을 가진 그 순간부터는 죽지 않도록 조심해야 한답니다. 만일 죽게 되면 그동안 쌓였던 경험치가 깍일 수 있으니까요.", false, true) == false then
		return
	end
	me:dialog(npc, "제가 가르쳐 드릴 수 있는건 여기까지 입니다. 이곳 저곳 여행을 하면서 자기 자신을 단련시키는 일만이 남았어요. 자신이 더욱 강해졌다고 생각되면 다시 절 찾아오세요. 당신을 계속 기다리고 있겠습니다.")
end

local function second_job(me, npc)
	if item_count(me, LETTER) >= 1 then
		me:dialog(npc, "아직 그를 만나지 못한겁니까? 헤네시스 근처 #b던전으로가는길#k 어딘가에 있는 #b궁수 전직 교관#k을 찾아가 보세요... 그에게 편지를 전해주면 어떻게 해야 하는지 자세한 내용을 들을 수 있을 겁니다...")
		return
	end

	if item_count(me, PROOF) < 1 then
		if me:dialog_yes_no(npc, "음... 당신 몰라보게 성장하셨군요! 예전에 그 허약한 모습은 어디로 가고 지금은 궁수로서의 위엄이 넘쳐 흐르고 있군요. 자... 어떤가요? 여기서 조금 더 강해지고 싶지 않은가요? 간단한 시험만 통과한다면 당신을 더욱 더 강하게 만들어 드리겠습니다. 해보시겠어요?") == false then
			me:dialog(npc, "흐음.. 더 강해질 수 있는 기회인데.. 아쉽군요. 마음이 바뀌면 다시 찾아오세요.")
			return
		end
		if me:exchange(nil, { item = { [LETTER] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 슬롯이 부족해 보이시는군요. 기타탭을 충분히 비우신 후 다시 찾아오세요.")
			return
		end
		if me:dialog(npc, "잘 생각하셨습니다. 당신이 강해 보이기는 하지만 그것이 정말인지 확인해 볼 필요가 있습니다. 어렵지 않은 테스트니까 당신이라면 충분히 통과할 수 있을 겁니다. 자... 우선 여기 제 편지를 하나 받도록 하세요. 잊어버리지 않도록 조심하세요.", false, true) == false then
			return
		end
		me:dialog(npc, "이 편지를 헤네시스 근처 #b던전으로가는길#k 어딘가에 있는 #b궁수 전직 교관#k에게 전해주세요. 바쁜 저를 대신하여 교관일을 해 주고 있는 고마운 사람이랍니다. 편지를 전해주면 당신을 저 대신 시험해 봐 줄 것입니다. 자세한 것은 그에게 직접 듣도록 하세요. 그럼 무사히 돌아오길 바라겠습니다.")
		return
	end

	if me:dialog(npc, "후후... 무사히 돌아 오셨군요! 당신이라면 그런 시험쯤은 간단하게 통과할 거라고 생각했습니다. 당신이 정말 뛰어난 궁수라는 것을 인정하지요. 자... 당신을 더욱 더 강하게 만들어 드리겠습니다! 그 전에... 당신은 2개의 길 중 하나를 선택해야만 합니다. 어렵긴 하겠지만... 모르는 것이 있다면 물어보도록 하세요.", false, true) == false then
		return
	end
	local sel = me:dialog_list(npc, "자... 결정이 끝났다면 가장 아래에 있는 [직업을 선택하겠습니다!]를 선택해 주세요.\r\n\r\n#b", {
		"헌터에 대해 알려주세요",
		"사수에 대해 알려주세요",
		"직업을 선택하겠습니다!",
	})
	if sel == nil then
		return
	end
	if sel <= #SECOND_JOBS then
		me:dialog(npc, SECOND_JOBS[sel].info, true, false)
		return
	end

	local pick = me:dialog_list(npc, "자... 마음을 정하셨나요? 2차 전직 하고 싶은 직업을 선택해 주세요.\r\n\r\n#b", {
		"헌터(Hunter)",
		"사수(Crossbow Man)",
	})
	if pick == nil then
		return
	end
	local job = SECOND_JOBS[pick]
	if me:dialog_yes_no(npc, job.confirm) == false then
		me:dialog(npc, "신중하게 생각해 보시고 다시 말을 걸어주세요.")
		return
	end
	if me:skill_point() > (me:level() - 30) * 3 then
		me:dialog(npc, "음.. 30레벨 이전에 얻은 SP를 모두 투자하셔야 2차 전직을 하실 수 있답니다.")
		return
	end
	if me:exchange({ item = { [PROOF] = 1 } }, nil) ~= ExchangeResult.OK then
		return
	end
	me:class(job.class)
	me:dialog(npc, "좋습니다! 당신은 이제부터 " .. job.done)
end

return {
	on_click = function(me, npc)
		local q = me:quest(TRIAL_QUEST)
		if q:started() == false then
			q:start("0")
		end
		local value = q:record()

		if me:class() == Class.Beginner then
			first_job(me, npc)
			return
		end
		if me:class() == Class.Bowman and me:level() >= 30 then
			second_job(me, npc)
			return
		end

		if value == "job3_trial1_1" then
			if is_second_job(me) == false or me:level() < 70 then
				me:dialog(npc, DEFAULT_TEXT)
				return
			end
			q:record("job3_trial1_2")
			me:dialog(npc, "당신을 기다리고 있었습니다. 몇 일 전 오시리아 대륙의 #b레네#k님으로부터 당신에 대한 이야기를 전해 들었거든요. 좋아요... 당신의 힘을 제가 시험해 드리죠. 빅토리아 아일랜드 깊은숲 어딘가에 다른 차원으로 통하는 균열이 있습니다. 보통 사람들은 들어갈 수 없지만 당신에 한해서는 들어갈 수 있도록 해 두겠어요. 균열 안으로 들어가면 저의 분신을 만날 수 있는데 그 분신을 쓰러뜨리고 #b검은 부적#k을 얻어 제게 가져와 주세요.", false, true)
			return
		end
		if value == "job3_trial1_2" then
			if item_count(me, BLACK_CHARM) >= 1 then
				if me:exchange({ item = { [BLACK_CHARM] = 1 } }, { item = { [NECKLACE] = 1 } }) ~= ExchangeResult.OK then
					me:dialog(npc, "흐음. 인벤토리 공간이 부족한 것 같은데요. 기타 탭을 충분히 비우고 다시 오세요.")
					return
				end
				q:record("job3_trial1_3")
				me:dialog(npc, "이럴수가... 제 분신을 쓰러뜨리고 #b검은 부적#k을 가져오셨군요! 좋아요... 이걸로 당신의 힘은 충분히 증명되었습니다. 힘에 한해서는 #b3차 전직#k을 하기에 부족함이 없어 보이는군요. 약속대로 당신에게 #b강인함의 목걸이#k를 드리겠습니다. 이 목걸이를 가지고 오시리아의 #b레네#k님에게 돌아가면 두번째 시험을 치를 수 있겠지요. 그럼 당신이 무사히 3차 전직을 할 수 있기를 빌겠습니다.")
				return
			end
			if is_second_job(me) == false or me:level() < 70 then
				me:dialog(npc, DEFAULT_TEXT)
				return
			end
			me:dialog(npc, "저의 분신인 만큼 그 강함은 상상할 수 없는 정도일 겁니다. 각종 고급 기술을 사용하는 데다 당신 혼자서 1:1로 싸워야 하기 때문에 쉽지만은 않을거에요. 게다가 그곳은 이 세계와는 다른 차원의 세계인 만큼 보통 인간이 오래 머무는 것은 좋지 않으므로 최대한 신속하게 쓰러뜨리는 것이 중요합니다. 미리 만반의 준비를 마친 후 도전하도록 하세요. 그럼 당신이 #b검은 부적#k을 가지고 무사히 돌아오기만을 기다리겠습니다.")
			return
		end
		if value == "job3_trial1_3" then
			me:dialog(npc, "제가 당신에게 주었던 그 목걸이를 가지고 오시리아의 #b레네#k님에게 돌아가면 두번째 시험을 치를 수 있겠지요. 그럼 당신이 무사히 3차 전직을 할 수 있기를 빌겠습니다.")
			return
		end

		me:dialog(npc, GREETINGS[me:class()] or DEFAULT_TEXT)
	end
}
