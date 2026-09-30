-- NPC name (String.wz/Npc.img.xml): 주먹펴고 일어서

local LETTER = 4031008
local PROOF = 4031012
local BLACK_CHARM = 4031059
local NECKLACE = 4031057
local TRIAL_QUEST = 195000
local DEFAULT_TEXT = "강인한 힘을 가지고, 가볍고 강력한 갑옷을 통한 방어력을 구사하는 전사는 항상 강인한 직업이지, 어때? 전사에 관심이 있는가 ?"

local GREETINGS = {
	[Class.Warrior] = "자네는 얼마 전에 #b전사#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Fighter] = "자네는 얼마 전에 #b파이터#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Page] = "자네는 얼마 전에 #b페이지#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Spearman] = "자네는 얼마 전에 #b스피어맨#k이 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Crusader] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b크루세이더#k로 전직하였던가. #r#h ##k여 크루세이더로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.WhiteKnight] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b나이트#k로 전직하였던가. #r#h ##k여 나이트로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.DragonKnight] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b용기사#k로 전직하였던가. #r#h ##k여 용기사로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Hero] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b히어로#k로 전직하였던가. #r#h ##k여 히어로로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Paladin] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b팔라딘#k으로 전직하였던가. #r#h ##k여 팔라딘으로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.DarkKnight] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b다크나이트#k로 전직하였던가. #r#h ##k여 다크나이트로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
}

local SECOND_JOBS = {
	{
		class = Class.Fighter,
		info = "#b파이터#k는 검과 도끼에 숙련된 전사 클레스일세, #b분노#k와 같은 보조 스킬은 방어력은 떨어지지만 공격력을 더 강하게 만들어 주는 보조 스킬이지, 그리고 #b파워 가드#k라는 스킬은 데미지를 감소시켜주는 스킬이라네. 또한 #b파이널 어택#k과 같은 스킬로 기본 공격 외에 치명타를 한 방 먹일 수 있으며, #b부스터#k와 같은 스킬로 공격 속도를 한 단계 늘릴 수 있다네.",
		confirm = "#b파이터#k로 2차 전직하고 싶단 말이지? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없네. 그 결심... 틀림이 없는가?",
		done = "#b파이터#k일세. 파이터는 강함을 추구하면서 끊임없이 싸우는 자... 결코 그 의지를 꺽지 말고 앞으로 앞으로 나아가게나. 나의 힘으로 자네를 더욱 강하게 만들어 주겠네. 그리고 자네에게 파이터가 익힐 수 있는 스킬들이 적혀있는 책을 주었네... 그 책에는 여러가지 파이터와 관련된 스킬들이 들어 있네. 또한 자네에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어보게. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있다네. 명심해 두게나. 이제 파이터로써 자네는 더 한층 높은 전사가 되었으니. 열심히 수련해 주길 바라네. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록!",
	},
	{
		class = Class.Page,
		info = "#b페이지#k는 검과 둔기에 숙련된 전사 클레스일세, #b위협#k과 같은 보조 스킬은 몬스터를 겁을 주게 하여, 그 공격력을 약화시키는 스킬이지, 그리고 #b파워 가드#k라는 스킬은 데미지를 감소시켜주는 스킬이라네. 또한 #b파이널 어택#k과 같은 스킬로 기본 공격 외에 치명타를 한 방 먹일 수 있으며, #b부스터#k와 같은 스킬로 공격 속도를 한 단계 늘릴 수 있다네.",
		confirm = "#b페이지#k로 2차 전직하고 싶단 말이지? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없네. 그 결심... 틀림이 없는가?",
		done = "#b페이지#k일세. 페이지는 철벽 같은 방어를 추구하면서 체계적으로 싸우는 자... 결코 그 의지를 꺽지 말고 앞으로 앞으로 나아가게나. 나의 힘으로 자네를 더욱 강하게 만들어 주겠네.  그리고 자네에게 페이지가 익힐 수 있는 스킬들이 적혀있는 책을 주었네... 그 책에는 여러가지 페이지와 관련된 스킬들이 들어 있네. 또한 자네에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어보게. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있다네. 명심해 두게나. 이제 페이지로써 자네는 더 한층 높은 전사가 되었으니. 열심히 수련해 주길 바라네. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록!",
	},
	{
		class = Class.Spearman,
		info = "#b스피어맨#k은 창과 폴암에 숙련된 전사 클레스일세, 이 클레스는 방어력을 중심으로 움직이지. #b아이언 월#k와 같은 방어력을 올려주는 스킬과, 그리고 #b하이퍼 바디#k 스킬처럼 자신의 최대 체력을 일시적으로 증가시켜주는 스킬이 있다네. 또한 #b파이널 어택#k과 같은 스킬로 기본 공격 외에 치명타를 한 방 먹일 수 있으며, #b부스터#k와 같은 스킬로 공격 속도를 한 단계 늘릴 수 있다네.",
		confirm = "#b스피어맨#k으로 2차 전직하고 싶단 말이지? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없네. 그 결심... 틀림이 없는가?",
		done = "#b스피어맨#k일세. 스피어맨은 길다란 창이나 폴암을 사용하여 넓은 범위의 근접 공격을 추구하는 자... 결코 그 의지를 꺽지 말고 앞으로 앞으로 나아가게나. 나의 힘으로 자네를 더욱 강하게 만들어 주겠네.  그리고 자네에게 스피어맨이 익힐 수 있는 스킬들이 적혀있는 책을 주었네... 그 책에는 여러가지 스피어맨과 관련된 스킬들이 들어 있네. 또한 자네에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어보게. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있다네. 명심해 두게나. 이제 스피어맨으로써 자네는 더 한층 높은 전사가 되었으니. 열심히 수련해 주길 바라네. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록!",
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
	return class == Class.Fighter or class == Class.Page or class == Class.Spearman
end

local function first_job(me, npc)
	if me:dialog(npc, "전사가 되고 싶어서 나를 찾아 온 건가...그렇다면 제대로 찾아 왔네...", false, true) == false then
		return
	end
	if me:dialog(npc, "전사가 되고 싶은가? 하지만 조건이 필요한데 말야... #b레벨이 10 이상#k이어야 하네. 어디보자... 자네는 말이지...", false, true) == false then
		return
	end
	if me:level() < 10 then
		me:dialog(npc, "자네는 아직 수련이 더 필요한 몸인 것 같네, 좀 더 수련을 한 뒤에 찾아오게나~!")
		return
	end
	if me:dialog_yes_no(npc, "자네는 충분히 전사가 될 수 있어 보이는군... 아직 미숙한 자네지만 그 몸에서 풍기는 위압감... 충분해! 자네 전사가 되지 않겠나?") == false then
		me:dialog(npc, "그런가? 천천히 생각해 보고 다시 오게나.")
		return
	end
	if me:dialog(npc, "자네는 이제부터 전사일세!! 더욱 정진하길 바라네... 열심히 수련하라는 뜻에서 자네를 더욱 단련시켜 주지! 자아... 하압~!", false, true) == false then
		return
	end
	if me:empty_slots(InventoryType.Equipment) < 1 then
		me:dialog(npc, "장비 인벤토리를 비우고 다시 오게.")
		return
	end
	if me:exchange(nil, { item = { [1302077] = 1 } }) ~= ExchangeResult.OK then
		me:dialog(npc, "장비 인벤토리를 비우고 다시 오게.")
		return
	end
	me:class(Class.Warrior)
	reset_stats(me, 35, 4, 4, 4)
	if me:dialog(npc, "자네는 한층 강해졌다네. 자네에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어보게. 스킬을 올릴 수 있을거야. 단 처음부터 전부 올릴수 있는 건 아냐. 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있다네. 명심해 두게나.", false, true) == false then
		return
	end
	if me:dialog(npc, "한가지 더 주의해야 할 점이 있네. 초보자에서 직업을 가진 그 순간부터는 죽지 않도록 조심해야 하네... 만일 죽게 되면 그동안 쌓였던 경험치가 깍일 수 있으니까 말야...", false, true) == false then
		return
	end
	me:dialog(npc, "내가 자네에게 가르쳐 줄 수 있는건 여기까지 일세... 이제는 혼자서 자기 자신을 더욱 단련시키는 일만 남은 거겠지. 자신이 더욱 강해졌다고 생각하면 다시 날 찾아오게나.")
end

local function second_job(me, npc)
	if item_count(me, LETTER) >= 1 then
		me:dialog(npc, "아직 그를 만나지 못한건가... 페리온 근처 #b서쪽바위산4#k 어딘가에 있는 #b전사 전직 교관#k을 찾아가 보게. 그에게 편지를 전해주면 어떻게 해야 하는지 자세한 내용을 들을 수 있을 거야...")
		return
	end

	if item_count(me, PROOF) < 1 then
		if me:dialog_yes_no(npc, "헛... 자네 몰라보게 성장했군! 예전에 그 비리비리한 모습은 어디로 가고 지금은 전사로서의 위엄이 넘쳐 흐르는 걸! 자... 어때? 여기서 조금 더 강해지고 싶지 않은가? 간단한 시험만 통과한다면 자네를 더욱 더 강하게 만들어 주겠네! 해볼텐가?") == false then
			me:dialog(npc, "흠..그런가? 아쉽군. 더욱 강해질 수 있는 기회인데. 아직 다른 할 일이 있나 보군?")
			return
		end
		if me:exchange(nil, { item = { [LETTER] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "흠. 인벤토리 공간이 부족한 것 같군. 기타 탭의 슬롯을 충분히 비운 후 다시 찾아오게나.")
			return
		end
		if me:dialog(npc, "잘 생각했네... 자네가 강해 보이기는 하지만 그것이 정말인지 확인해 볼 필요가 있네. 어렵지 않은 테스트니까 자네라면 충분히 통과할 수 있을 것일세... 자... 우선 여기 내 편지를 하나 받도록 하게나. 잊어버리지 않도록 조심하게.", false, true) == false then
			return
		end
		me:dialog(npc, "이 편지를 페리온 근처 #b서쪽바위산4#k 어딘가에 있는 #b전사 전직 교관#k에게 전해주게. 바쁜 나를 대신하여 교관일을 해 주고 있는 고마운 사람이지. 편지를 전해주면 자네를 나 대신 시험해 봐 줄 것일세. 자세한 것은 그에게 직접 듣도록. 그럼 무사히 돌아오길 바라겠네.")
		return
	end

	if me:dialog(npc, "오오... 무사히 돌아왔군! 자네라면 그런 시험쯤은 간단하게 통과할 거라고 생각했네. 자네가 정말 강한 전사라는 것을 인정하지. 자... 자네를 더욱 더 강하게 만들어 주겠네. 그 전에...! 자네는 3개의 길 중 하나를 선택해야만 하네. 어렵긴 하겠지만... 모르는 것이 있다면 물어봐도 좋아.", false, true) == false then
		return
	end
	local sel = me:dialog_list(npc, "자... 결정이 끝났다면 가장 아래에 있는 [직업을 선택하겠습니다!]를 선택해 주게...\r\n\r\n#b", {
		"파이터에 대해 알려주세요",
		"페이지에 대해 알려주세요",
		"스피어맨에 대해 알려주세요",
		"직업을 선택하겠습니다!",
	})
	if sel == nil then
		return
	end
	if sel <= #SECOND_JOBS then
		me:dialog(npc, SECOND_JOBS[sel].info, true, false)
		return
	end

	local pick = me:dialog_list(npc, "자... 마음을 정했는가? 2차 전직 하고 싶은 직업을 선택해 보게나.\r\n\r\n#b", {
		"파이터(Fighter)",
		"페이지(Page)",
		"스피어맨(Spearman)",
	})
	if pick == nil then
		return
	end
	local job = SECOND_JOBS[pick]
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
	me:dialog(npc, "좋아! 자네는 이제부터 " .. job.done)
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
		if me:class() == Class.Warrior and me:level() >= 30 then
			second_job(me, npc)
			return
		end

		if value == "job3_trial1_1" then
			if is_second_job(me) == false or me:level() < 70 then
				me:dialog(npc, DEFAULT_TEXT)
				return
			end
			q:record("job3_trial1_2")
			me:dialog(npc, "자네를 기다리고 있었네. 몇 일 전 오시리아 대륙의 #b타일러스#k님으로부터 자네에 대한 이야기를 전해 들었지. 좋아... 자네의 힘을 내가 시험해 주지. 빅토리아 아일랜드 개미굴 어딘가에 다른 차원으로 통하는 균열이 있다네. 보통 사람들은 들어갈 수 없지만 자네에 한해서는 들어갈 수 있도록 해 두겠네. 균열 안으로 들어가면 나의 분신을 만날 수 있는데 그 분신을 쓰러뜨리고 #b검은 부적#k을 얻어 내게 가져와 주게나.", false, true)
			return
		end
		if value == "job3_trial1_2" then
			if item_count(me, BLACK_CHARM) >= 1 then
				if me:exchange({ item = { [BLACK_CHARM] = 1 } }, { item = { [NECKLACE] = 1 } }) ~= ExchangeResult.OK then
					me:dialog(npc, "흐음. 인벤토리 공간이 부족한 것 같은데. 기타 탭을 충분히 비우고 다시 오게.")
					return
				end
				q:record("job3_trial1_3")
				me:dialog(npc, "이럴수가... 나의 분신을 쓰러뜨리고 #b검은 부적#k을 가져왔군. 그래! 좋아... 이걸로 너의 힘은 충분히 증명되었다네. 힘에 한해서는 #b3차 전직#k을 하기에 부족함이 없어 보이는군. 약속대로 너에게 #b강인함의 목걸이#k를 주겠네. 이 목걸이를 가지고 오시리아의 #b타일러스#k님에게 돌아가면 두번째 시험을 치를 수 있겠지. 그럼 네가 무사히 3차 전직을 할 수 있기를 빌겠네.")
				return
			end
			if is_second_job(me) == false or me:level() < 70 then
				me:dialog(npc, DEFAULT_TEXT)
				return
			end
			me:dialog(npc, "나의 분신인 만큼 그 강함은 상상할 수 없는 정도일 걸세. 각종 고급 기술을 사용하는 데다 자네 혼자서 1:1로 싸워야 하기 때문에 쉽지만은 않을거야. 게다가 그곳은 이 세계와는 다른 차원의 세계인 만큼 보통 인간이 오래 머무는 것은 좋지 않으므로 최대한 신속하게 쓰러뜨리는 것이 중요하다네. 미리 만반의 준비를 마친 후 도전하도록 하게나. 그럼 자네가 #b검은 부적#k을 가지고 무사히 돌아오기만을 기다리겠네.")
			return
		end
		if value == "job3_trial1_3" then
			me:dialog(npc, "내가 자네에게 주었던 그 목걸이를 가지고 오시리아의 #b타일러스#k님에게 돌아가면 두번째 시험을 치를 수 있겠지. 그럼 네가 무사히 3차 전직을 할 수 있기를 빌겠네.")
			return
		end

		me:dialog(npc, GREETINGS[me:class()] or DEFAULT_TEXT)
	end
}
