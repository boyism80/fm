-- NPC name (String.wz/Npc.img.xml): 다크로드

local SCRIPT = "script/npc/1052001.lua"
local TRAINING_MAP = 910300000
local TRAINING_RETURN = 103000003
local TRAINING_SECONDS = 300
local LETTER = 4031011
local PROOF = 4031012
local BLACK_CHARM = 4031059
local NECKLACE = 4031057

local INTRO = "어둠 속에서 단시간에 적의 숨통을 끊을 수 잇는 민첩함과, 항상 신속함으로 인해 즐거움을 달고 다니는 직업. 도적에 관심이 있는건가?"

local ADVANCED = {
	[Class.Thief] = "자네는 얼마 전에 #b도적#k이 되었던 #r#h ##k이군, 수행은 잘 되는가? 더 강한 도적이 되기 위하여 더 정진하길 바라네..",
	[Class.Assassin] = "자네는 얼마 전에 #b어쌔신#k이 되었던 #r#h ##k이군, 수행은 잘 되는가? 더 강한 도적이 되기 위하여 더 정진하길 바라네..",
	[Class.Bandit] = "자네는 얼마 전에 #b시프#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 더 강한 도적이 되기 위하여 더 정진하길 바라네..",
	[Class.Hermit] = "너의 소식은 잘 알고 있지, 얼마 전에 #b허밋#k으로 전직하였던가? 정말 축하한다, 더 강한 도적이 되기 위하여 더 정진하길 바란다.",
	[Class.ChiefBandit] = "너의 소식은 잘 알고 있지, 얼마 전에 #b시프마스터#k로 전직하였던가? 정말 축하한다, 더 강한 도적이 되기 위하여 더 정진하길 바란다",
	[Class.Nightlord] = "너의 소식은 잘 알고 있지, 얼마 전에 #b나이트로드#k로 전직하였던가? 정말 축하한다, 더 강한 도적이 되기 위하여 더 정진하길 바란다.",
	[Class.Shadower] = "너의 소식은 잘 알고 있지, 얼마 전에 #b섀도어#k로 전직하였던가? 정말 축하한다, 더 강한 도적이 되기 위하여 더 정진하길 바란다",
}

local SECOND_JOBS = {
	{
		class = Class.Assassin,
		info = "#b어쌔신#k은 표창을 더 자유 자재로 사용할 수 있는 도적 클레스야, #b헤이스트#k로 파티원의 이동속도와 점프력을 올려주는 스킬이나, #b드레인#k같은 적의 HP를 흡수하여 자신의 HP로 치환하는 스킬, 또한 #b부스터#k와 같은 스킬로 공격 속도를 한 단계 늘릴 수 있어.",
		confirm = "#b어쌔신#k으로 2차 전직하고 싶단 말이지? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없어. 그 결심... 틀림이 없겠지?",
		done = "#b어쌔신#k이야. 정확한 투척과 민첩성으로 적을 제압하는 자... 더욱 수련에 정진하길 바라겠어... 나의 힘으로 널 더욱 강하게 만들어 주지! 그리고 너에게 어쌔신이 익힐 수 있는 스킬들이 적혀있는 책을 주었어... 그 책에는 여러가지 어쌔신에 관련된 스킬들이 들어 있어. 또한 너에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어봐. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있으니까, 명심해 둬. 이제 어쌔신으로써 너는 더 한층 높은 도적이 되었어. 열심히 수련해 주길 바래. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록해.",
	},
	{
		class = Class.Bandit,
		info = "#b시프#k는 단검을 더 자유 자재로 사용할 수 있는 도적 클레스야, #b헤이스트#k로 파티원의 이동속도와 점프력을 올려주는 스킬이나, #b스틸#k같은 적이 소지하고있는 아이템을 일정 확률로 훔쳐오는 스킬, 또한 #b부스터#k와 같은 스킬로 공격 속도를 한 단계 늘릴 수 있어.",
		confirm = "#b시프#k로 2차 전직하고 싶단 말이지? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없어. 그 결심... 틀림이 없겠지?",
		done = "#b시프#k야. 빠른 근접 공격과 민첩성으로 적을 제압하는 자... 더욱 수련에 정진하길 바라겠어... 나의 힘으로 널 더욱 강하게 만들어 주지! 그리고 너에게 어쌔신이 익힐 수 있는 스킬들이 적혀있는 책을 주었어... 그 책에는 여러가지 시프에 관련된 스킬들이 들어 있어. 또한 너에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어봐. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있으니까, 명심해 둬. 이제 시프로써 너는 더 한층 높은 도적이 되었어. 열심히 수련해 주길 바래. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록해.",
	},
}

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

local function enter_training(me, npc)
	local ok, count = run_on_map(TRAINING_MAP, SCRIPT, "character_count")
	if ok and count ~= nil and count > 0 then
		me:dialog(npc, "흠.. 훈련장에 이미 누군가가 들어간 것 같구만. 나중에 다시 찾아오게나.")
		return
	end
	local q = me:quest(110114)
	if q:started() == false then
		q:start("0")
	end
	local last = tonumber(q:record()) or 0
	if last + TRAINING_SECONDS >= now() then
		me:dialog(npc, "훈련장은 5분에 한번씩만 입장할 수 있다네. 나중에 다시 찾아와보게나.")
		return
	end
	q:record(tostring(now()))
	run_on_map(TRAINING_MAP, SCRIPT, "reset_map")
	local function on_arrive(me)
		me:clock(TRAINING_SECONDS, function(me)
			me:map(TRAINING_RETURN)
		end)
	end
	me:map(TRAINING_MAP, 0, { callback = on_arrive })
end

local function first_job(me, npc)
	if me:dialog(npc, "도적이 되고 싶어서 나를 찾아 온 건가...그렇다면 제대로 찾아 왔어...", false, true) == false then
		return
	end
	if me:dialog(npc, "도적이 되고 싶어? 하지만 조건이 필요한데 말야... #b레벨이 10 이상#k이어야 해. 어디보자... 너는 말이지...", false, true) == false then
		return
	end
	if me:level() < 10 then
		me:dialog(npc, "내가 보니까, 너는 아직 애송이에 불과하군, 좀 더 수련을 한 뒤에 찾아오면 그때 받아주지.")
		return
	end
	if me:dialog_yes_no(npc, "너는 충분히 도적이 될 수 있어 보이는군... 아직 미숙한 너지만 너의 몸에선 놀라운 민첩함이 느껴져... 너 도적이 되지 않겠어?") == false then
		me:dialog(npc, "그런가? 천천히 생각해 보고 다시 오게나.")
		return
	end
	if me:dialog(npc, "너는 이제부터 도적이야.. 더욱 정진하길 바래... 열심히 수련하라는 뜻에서 너에게 힘을 추가시켜 줄게.. 자아... 하압~!", false, true) == false then
		return
	end
	if me:exchange(nil, { item = { [1332063] = 1, [1472061] = 1, [2070015] = 500 } }) ~= ExchangeResult.OK then
		me:dialog(npc, "장비와 소비 인벤토리를 비우고 다시 오게.")
		return
	end
	me:class(Class.Thief)
	local total = me:base_str() + me:base_dex() + me:base_int() + me:base_luk()
	me:ability_point(me:ability_point() + total - (4 + 25 + 4 + 4))
	me:base_str(4)
	me:base_dex(25)
	me:base_int(4)
	me:base_luk(4)
	if me:dialog(npc, "너에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어보도록 해. 스킬을 올릴 수 있을거야. 단 처음부터 전부 올릴수 있는 건 아냐. 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있어. 명심해 둬.", false, true) == false then
		return
	end
	if me:dialog(npc, "한가지 더 주의해야 할 점이 있어. 초보자에서 직업을 가진 그 순간부터는 죽지 않도록 조심해야 해... 만일 죽게 되면 그동안 쌓였던 경험치가 깍일 수 있으니까 말야... ", false, true) == false then
		return
	end
	me:dialog(npc, "내가 너에게 가르쳐 줄 수 있는건 여기까지야... 이제는 혼자서 자기 자신을 더욱 단련시키는 일만 남은 거겠지. 자신이 더욱 강해졌다고 생각하면 다시 날 찾아오도록 해.")
end

local function second_job(me, npc)
	local job = nil
	while job == nil do
		local sel = me:dialog_list(npc, "자... 결정이 끝났다면 가장 아래에 있는 [직업을 선택하겠습니다!]를 선택해...\r\n\r\n#b", {
			"어쌔신에 대해 알려주세요",
			"시프에 대해 알려주세요",
			"직업을 선택하겠습니다!",
		})
		if sel == nil then
			return
		end
		if sel == 3 then
			local pick = me:dialog_list(npc, "자... 마음을 정했어? 2차 전직 하고 싶은 직업을 선택해 봐.\r\n\r\n#b", {
				"어쌔신(Assassin)",
				"시프(Thirf)",
			})
			if pick == nil then
				return
			end
			job = SECOND_JOBS[pick]
		else
			me:dialog(npc, SECOND_JOBS[sel].info, true, false)
		end
	end
	if me:dialog_yes_no(npc, job.confirm) == false then
		me:dialog(npc, "신중하게 생각해 본 후 다시 말을 걸어주게나.")
		return
	end
	if me:skill_point() > (me:level() - 30) * 3 then
		me:dialog(npc, "음.. SP를 아직 다 사용하지 않은 것 같군? 자네는 SP가 너무 많아 남아 아직 2차전직을 할 수 없다네.")
		return
	end
	if me:exchange({ item = { [PROOF] = 1 } }, nil) ~= ExchangeResult.OK then
		return
	end
	me:class(job.class)
	me:dialog(npc, "좋아. 너는 이제부터 " .. job.done)
end

local function second_job_test(me, npc)
	if item_count(me, LETTER) > 0 then
		me:dialog(npc, "아직 그를 만나지 못한거야? 커닝시티 근처 #b커닝시티북쪽공사장#k 어딘가에 있는 #b도적 전직 교관#k을 찾아가 봐... 그에게 편지를 전해주면 어떻게 해야 하는지 자세한 내용을 들을 수 있을 거야...")
		return
	end
	if item_count(me, PROOF) > 0 then
		if me:dialog(npc, "후후... 무사히 돌아왔군! 너라면 그런 시험쯤은 간단하게 통과할 거라고 생각했어. 네가 정말 뛰어난 도적이라는 것을 인정하지. 자... 널 더욱 더 강하게 만들어 주겠어! 그 전에... 너는 2개의 길 중 하나를 선택해야만 해. 어렵긴 하겠지만... 모르는 것이 있다면 물어보도록 해.", false, true) == false then
			return
		end
		second_job(me, npc)
		return
	end
	if me:dialog_yes_no(npc, "흠... 너 몰라보게 강해졌군! 예전에 그 허약한 모습은 어디로 가고 지금은 도적으로서의 위엄이 넘쳐 흐르고 있군... 자... 어때? 여기서 조금 더 강해지고 싶지 않은가? 간단한 시험만 통과한다면 널 더욱 더 강하게 만들어 주겠네! 해보겠나?") == false then
		me:dialog(npc, "흠..그런가? 아쉽군. 더욱 강해질 수 있는 기회인데. 아직 다른 할 일이 있나 보군?")
		return
	end
	if me:exchange(nil, { item = { [LETTER] = 1 } }) ~= ExchangeResult.OK then
		me:dialog(npc, "흠. 인벤토리 공간이 부족한 것 같군. 기타 탭의 슬롯을 충분히 비운 후 다시 찾아오게나.")
		return
	end
	if me:dialog(npc, "잘 생각했어... 하지만 넌 강해 보이기는 하지만 그것이 정말인지 확인해 볼 필요가 있어. 어렵지 않은 테스트니까 너라면 충분히 통과할 수 있을 거야... 자... 우선 여기 내 편지를 하나 받도록 해. 잊어버리지 않도록 조심하라구.", false, true) == false then
		return
	end
	me:dialog(npc, "이 편지를 커닝시티 근처 #b커닝시티북쪽공사장#k 어딘가에 있는 #b도적 전직 교관#k에게 전해줘. 바쁜 나를 대신하여 교관일을 해 주고 있는 고마운 사람이야. 편지를 전해주면 널 내 대신 시험해 봐 줄 거야. 자세한 것은 그에게 직접 듣도록 해. 그럼 무사히 돌아오길 바라겠어.")
end

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	reset_map = function(map)
		map:reset()
	end,

	on_click = function(me, npc)
		if me:quest(6141):started() then
			enter_training(me, npc)
			return
		end
		local class = me:class()
		if class == Class.Beginner then
			first_job(me, npc)
			return
		end
		if class == Class.Thief and me:level() >= 30 then
			second_job_test(me, npc)
			return
		end

		local q = me:quest(195000)
		local value = q:record()
		local third_ready = (class == Class.Assassin or class == Class.Bandit) and me:level() >= 70
		if value == "job3_trial1_1" then
			if third_ready == false then
				me:dialog(npc, INTRO)
				return
			end
			q:record("job3_trial1_2")
			me:dialog(npc, "자네를 기다리고 있었네. 몇 일 전 오시리아 대륙의 #b아레크#k님으로부터 자네에 대한 이야기를 전해 들었지. 좋아... 자네의 힘을 내가 시험해 주지. 빅토리아 아일랜드 개미굴 어딘가에 다른 차원으로 통하는 균열이 있다네. 보통 사람들은 들어갈 수 없지만 자네에 한해서는 들어갈 수 있도록 해 두겠네. 균열 안으로 들어가면 나의 분신을 만날 수 있는데 그 분신을 쓰러뜨리고 #b검은 부적#k을 얻어 내게 가져와 주게나.")
		elseif value == "job3_trial1_2" then
			if item_count(me, BLACK_CHARM) > 0 then
				if me:exchange({ item = { [BLACK_CHARM] = 1 } }, { item = { [NECKLACE] = 1 } }) ~= ExchangeResult.OK then
					me:dialog(npc, "흐음. 인벤토리 공간이 부족한 것 같은데. 기타 탭을 충분히 비우고 다시 오게.")
					return
				end
				q:record("job3_trial1_3")
				me:dialog(npc, "이럴수가... 나의 분신을 쓰러뜨리고 #b검은 부적#k을 가져왔군. 그래! 좋아... 이걸로 너의 힘은 충분히 증명되었다네. 힘에 한해서는 #b3차 전직#k을 하기에 부족함이 없어 보이는군. 약속대로 너에게 #b강인함의 목걸이#k를 주겠네. 이 목걸이를 가지고 오시리아의 #b아레크#k님에게 돌아가면 두번째 시험을 치를 수 있겠지. 그럼 네가 무사히 3차 전직을 할 수 있기를 빌겠네.")
			elseif third_ready then
				me:dialog(npc, "나의 분신인 만큼 그 강함은 상상할 수 없는 정도일 걸세. 각종 고급 기술을 사용하는 데다 자네 혼자서 1:1로 싸워야 하기 때문에 쉽지만은 않을거야. 게다가 그곳은 이 세계와는 다른 차원의 세계인 만큼 보통 인간이 오래 머무는 것은 좋지 않으므로 최대한 신속하게 쓰러뜨리는 것이 중요하다네. 미리 만반의 준비를 마친 후 도전하도록 하게나. 그럼 자네가 #b검은 부적#k을 가지고 무사히 돌아오기만을 기다리겠네.")
			else
				me:dialog(npc, INTRO)
			end
		elseif value == "job3_trial1_3" then
			me:dialog(npc, "내가 너에게 주었던 그 목걸이를 가지고 오시리아의 #b아레크#k님에게 돌아가면 두번째 시험을 치를 수 있겠지. 그럼 네가 무사히 3차 전직을 할 수 있기를 빌겠네.")
		else
			me:dialog(npc, ADVANCED[class] or INTRO)
		end
	end
}
