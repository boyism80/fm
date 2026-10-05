-- NPC name (String.wz/Npc.img.xml): 하인즈

local LETTER = 4031009
local PROOF = 4031012
local BLACK_CHARM = 4031059
local NECKLACE = 4031057
local TRIAL_QUEST = 195000
local DEFAULT_TEXT = "여러가지 종류의 화려한 마법을 구사하며 적을 물리치는 직업. 마법사는 정말 정교하고 다양한 직업이지. 혹시 마법사에 관심이 있는 건가?"

local GREETINGS = {
	[Class.Magician] = "자네는 얼마 전에 #b마법사#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.FpWizard] = "자네는 얼마 전에 #b불,독 계열의 위자드#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.IlWizard] = "자네는 얼마 전에 #b얼음,번개 계열의 위자드#k가 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Cleric] = "자네는 얼마 전에 #b클레릭#k이 되었던 #r#h ##k이군, 수행은 잘 되는가? 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.FpArchMage] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b불,독 계열의 아크메이지#k로 전직하였던가. #r#h ##k여 아크메이지로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.IlArchMage] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b얼음,번개 계열의 아크메이지#k으로 전직하였던가. #r#h ##k여 아크메이지로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Bishop] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b비숍#k으로 전직하였던가. #r#h ##k여 비숍으로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.FpMage] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b불,독 계열의 메이지#k로 전직하였던가. #r#h ##k여 메이지로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.IlMage] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b얼음,번개 계열의 메이지#k으로 전직하였던가. #r#h ##k여 메이지로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
	[Class.Priest] = "자네의 소식은 잘 알고 있다네, 얼마 전에 #b프리스트#k로 전직하였던가. #r#h ##k여 프리스트로 전직한 것을 축하하네, 더 강한 존재가 되기 위하여 더 정진하길 바라네..",
}

local SECOND_CLASSES = {
	{
		class = Class.FpWizard,
		info = "#b불, 독계열의 위자드#k는 불과 독마법과 공격 마법에 숙련된 마법사 클레스일세, #b메디테이션#k과 같이 자신과 파티원의 마력을 올려주는 스킬이나, #bMP 이터#k같은 적의 MP를 흡수하여 자신의 MP로 치환하는 스킬, 또한 #b파이어 에로우#k나 #b포이즌 브레스#k 같은 속성 공격 스킬이 있다네. 이 스킬들은 몬스터의 속성에 따라 데미지가 증폭되기도 하고, 반감되기도 하니, 몬스터의 속성을 잘 이용해야 한다네.",
		confirm = "#b불,독 계열의 위자드#k로 2차 전직하고 싶단 말이지? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없네. 그 결심... 틀림이 없는가?",
		done = "#b불,독 계열의 위자드#k일세. 위자드는 높은 지능을 바탕으로 초자연의 힘으로 적을 제압하는 자... 더욱 학업에 정진하도록 바라네... 나의 힘으로 자네를 더욱 강하게 만들어 주겠네. 그리고 자네에게 불, 독 계열의 위자드가 익힐 수 있는 스킬들이 적혀있는 책을 주었네... 그 책에는 여러가지 위자드와 관련된 스킬들이 들어 있네. 또한 자네에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어보게. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있다네. 명심해 두게나. 이제 위자드로써 자네는 더 한층 높은 마법사가 되었으니. 열심히 수련해 주길 바라네. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록!",
	},
	{
		class = Class.IlWizard,
		info = "#b얼음, 번개계열의 #k은 얼음과 번개마법과 공격 마법에 숙련된 마법사 클레스일세, #b메디테이션#k과 같이 자신과 파티원의 마력을 올려주는 스킬이나, #bMP 이터#k같은 적의 MP를 흡수하여 자신의 MP로 치환하는 스킬, 또한 #b콜드 빔#k이나 #b썬더 볼트#k 같은 속성 공격 스킬이 있다네. 이 스킬들은 몬스터의 속성에 따라 데미지가 증폭되기도 하고, 반감되기도 하니, 몬스터의 속성을 잘 이용해야 한다네.",
		confirm = "#b얼음,번개 계열의 위자드#k로 2차 전직하고 싶단 말이지? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없네. 그 결심... 틀림이 없는가?",
		done = "#b얼음,번개 계열의 위자드#k일세. 위자드는 높은 지능을 바탕으로 초자연의 힘으로 적을 제압하는 자... 더욱 학업에 정진하도록 바라네... 나의 힘으로 자네를 더욱 강하게 만들어 주겠네. 그리고 자네에게 얼음, 번개 계열의 위자드가 익힐 수 있는 스킬들이 적혀있는 책을 주었네... 그 책에는 여러가지 위자드와 관련된 스킬들이 들어 있네. 또한 자네에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어보게. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있다네. 명심해 두게나. 이제 위자드로써 자네는 더 한층 높은 마법사가 되었으니. 열심히 수련해 주길 바라네. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록!",
	},
	{
		class = Class.Cleric,
		info = "#b클레릭#k은 성 속성 마법과 보조 마법에 숙련된 마법사 클레스일세, 적의 물리 공격 데미지를 비례적으로 줄여주는 #b인빈서블#k이나 자신과 파티원 모두 능력치를 올려주는 #b블레스#k같은 보조 스킬이 있다네, 그리고 #b힐#k은 파티원의 체력을 회복시키기도 하지만, 언데드 속성의 적에게는 데미지를 줄 수 있다네. 또한 #b홀리 에로우#k라는 스킬이 있는데 이 스킬은 몬스터의 속성에 따라 데미지가 증폭되기도 하고, 반감되기도 하니, 몬스터의 속성을 잘 이용해야 한다네.",
		confirm = "#b클레릭#k으로 2차 전직하고 싶단 말이지? 한번 결정하면 다른 2차 전직 직업으로는 전직할 수 없네. 그 결심... 틀림이 없는가?",
		done = "#b클레릭#k일세. 클레릭은 신을 섬기는 마음으로 만물에게 생명의 힘을 주는 자... 더욱 깊은 신앙심을 가지도록 바라네... 나의 힘으로 자네를 더욱 강하게 만들어 주겠네. 그리고 자네에게 클레릭이 익힐 수 있는 스킬들이 적혀있는 책을 주었네... 그 책에는 여러가지 클레릭에 관련된 스킬들이 들어 있네.  또한 자네에게 약간의 #bSP#k를 주었으니 #bSkill 메뉴#k를 열어보게. 스킬을 올릴 수 있을거야. 참고로 1차 전직 때처럼 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있다네. 명심해 두게나. 이제 클레릭으로써 자네는 더 한층 높은 마법사가 되었으니. 열심히 수련해 주길 바라네. 그리고 자신이 강하다고 생각할 때가 되면 나를 찾아오도록!",
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

local function is_second_class(me)
	local class = me:class()
	return class == Class.FpWizard or class == Class.IlWizard or class == Class.Cleric
end

local function first_class(me, npc)
	if me:dialog(npc, "마법사가 되고 싶은가? 하지만 조건이 필요한데 말야... #b레벨이 8 이상#k이어야 하네. 어디보자... 자네는 말이지...", false, true) == false then
		return
	end
	if me:level() < 8 then
		me:dialog(npc, "자네는 아직 수련이 더 필요한 몸인 것 같군, 좀 더 수련을 한 뒤에 찾아와 주게.")
		return
	end
	if me:dialog_yes_no(npc, "자네는 충분히 마법사가 될 수 있어 보이는군... 아직 미숙한 자네지만 그 몸에서는 마나의 기운이 느껴진다.. 자네 마법사로 전직하고 싶은가?") == false then
		me:dialog(npc, "그런가? 천천히 생각해 보고 다시 오게나.")
		return
	end
	if me:dialog(npc, "좋네! 자네는 이제부터 마법사가 되었네! 이제 마나를 자유적으로 다룰 수 있을 걸세. 작지만 그대에게 내가 가진 능력의 일부를 조금 보태주지.. 허어~~~ 업!!!", false, true) == false then
		return
	end
	if me:empty_slots(InventoryType.Equipment) < 1 then
		me:dialog(npc, "장비 인벤토리를 비우고 다시 오게.")
		return
	end
	if me:exchange(nil, { item = { [1372043] = 1 } }) ~= ExchangeResult.OK then
		me:dialog(npc, "장비 인벤토리를 비우고 다시 오게.")
		return
	end
	me:class(Class.Magician)
	reset_stats(me, 4, 4, 25, 4)
	if me:dialog(npc, "자네에게 약간의 #bSP#k를 주었네. 왼쪽 하 단에 있는 #bSkill 메뉴#k를 열어보게나. 스킬을 올릴 수 있을 것이네. 단 처음부터 전부 올릴수 있는 건 아니라네... 다른 스킬들을 어느 정도 익혀야만 배울수 있는 스킬도 있다는거 명심해주게나.", false, true) == false then
		return
	end
	if me:dialog(npc, "한 가지 더 주의해야 할 점이 있네. 초보자에서 직업을 가진 그 순간부터는 죽지 않도록 조심해야 한다네. 만일 죽게 되면 그동안 쌓였던 경험치가 깍일 수 있게 되거든.", false, true) == false then
		return
	end
	me:dialog(npc, "내가 가르쳐 줄 수 있는건 여기까지 이네. 이곳 저곳 여행을 하면서 필요한 것은 당연 수련일세. 자신이 더욱 강해졌다고 생각되면 다시 날 찾아와주게. 그대를 계속 지켜보고 있겠네.")
end

local function second_class(me, npc)
	if item_count(me, LETTER) >= 1 then
		me:dialog(npc, "아직 그를 만나지 못한건가... 엘리니아 근처 #b엘리니아 북쪽숲4#k 어딘가에 있는 #b마법사 전직 교관#k을 찾아가 보게... 그에게 편지를 전해주면 어떻게 해야 하는지 자세한 내용을 들을 수 있을 거야...")
		return
	end

	if item_count(me, PROOF) < 1 then
		if me:dialog_yes_no(npc, "흠... 자네 몰라보게 성장했군! 예전에 그 허약한 모습은 어디로 가고 지금은 마법사로서의 위엄이 넘쳐 흐르고 있군... 자... 어떤가? 여기서 조금 더 강해지고 싶지 않은가? 간단한 시험만 통과한다면 자네를 더욱 더 강하게 만들어 주겠네! 해볼텐가?") == false then
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
		me:dialog(npc, "이 편지를 엘리니아 근처 #b엘리니아 북쪽숲4#k 어딘가에 있는 #b마법사 전직 교관#k에게 전해주게. 바쁜 나를 대신하여 교관일을 해 주고 있는 고마운 사람이지. 편지를 전해주면 자네를 나 대신 시험해 봐 줄 것일세. 자세한 것은 그에게 직접 듣도록 하게... 그럼 무사히 돌아오길 바라겠네.")
		return
	end

	if me:dialog(npc, "오오... 무사히 돌아왔군! 자네라면 그런 시험쯤은 간단하게 통과할 거라고 생각했네. 자네가 정말 강한 마법사라는 것을 인정하지. 자... 자네를 더욱 더 강하게 만들어 주겠네. 그 전에... 자네는 3개의 길 중 하나를 선택해야만 하네. 어렵긴 하겠지만... 모르는 것이 있다면 물어봐도 좋아.", false, true) == false then
		return
	end
	local sel = me:dialog_list(npc, "자... 결정이 끝났다면 가장 아래에 있는 [직업을 선택하겠습니다!]를 선택해 주게...\r\n\r\n#b", {
		"위자드(불, 독)에 대해 알려주세요",
		"위자드(얼음, 번개)에 대해 알려주세요",
		"클레릭에 대해 알려주세요",
		"직업을 선택하겠습니다!",
	})
	if sel == nil then
		return
	end
	if sel <= #SECOND_CLASSES then
		me:dialog(npc, SECOND_CLASSES[sel].info, true, false)
		return
	end

	local pick = me:dialog_list(npc, "자... 마음을 정했는가? 2차 전직 하고 싶은 직업을 선택해 보게나...\r\n\r\n#b", {
		"불, 독 계열의 위자드(Wizard)",
		"얼음, 번개 계열의 위자드(Wizard)",
		"클레릭(Clelic)",
	})
	if pick == nil then
		return
	end
	local next_class = SECOND_CLASSES[pick]
	if me:dialog_yes_no(npc, next_class.confirm) == false then
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
	me:class(next_class.class)
	me:dialog(npc, "좋네. 자네는 이제부터 " .. next_class.done)
end

return {
	on_click = function(me, npc)
		local q = me:quest(TRIAL_QUEST)
		if q:started() == false then
			q:start("0")
		end
		local value = q:record()

		if me:class() == Class.Beginner then
			first_class(me, npc)
			return
		end
		if me:class() == Class.Magician and me:level() >= 30 then
			second_class(me, npc)
			return
		end

		if value == "job3_trial1_1" then
			if is_second_class(me) == false or me:level() < 70 then
				me:dialog(npc, DEFAULT_TEXT)
				return
			end
			q:record("job3_trial1_2")
			me:dialog(npc, "자네를 기다리고 있었네. 몇 일 전 오시리아 대륙의 #b로베이라#k님으로부터 자네에 대한 이야기를 전해 들었지. 좋아... 자네의 힘을 내가 시험해 주지. 빅토리아 아일랜드 개미굴 어딘가에 다른 차원으로 통하는 균열이 있다네. 보통 사람들은 들어갈 수 없지만 자네에 한해서는 들어갈 수 있도록 해 두겠네. 균열 안으로 들어가면 나의 분신을 만날 수 있는데 그 분신을 쓰러뜨리고 #b검은 부적#k을 얻어 내게 가져와 주게나.", false, true)
			return
		end
		if value == "job3_trial1_2" then
			if item_count(me, BLACK_CHARM) >= 1 then
				if me:exchange({ item = { [BLACK_CHARM] = 1 } }, { item = { [NECKLACE] = 1 } }) ~= ExchangeResult.OK then
					me:dialog(npc, "흐음. 인벤토리 공간이 부족한 것 같은데. 기타 탭을 충분히 비우고 다시 오게.")
					return
				end
				q:record("job3_trial1_3")
				me:dialog(npc, "이럴수가... 나의 분신을 쓰러뜨리고 #b검은 부적#k을 가져왔군. 그래! 좋아... 이걸로 너의 힘은 충분히 증명되었다네. 힘에 한해서는 #b3차 전직#k을 하기에 부족함이 없어 보이는군. 약속대로 너에게 #b강인함의 목걸이#k를 주겠네. 이 목걸이를 가지고 오시리아의 #b로베이라#k님에게 돌아가면 두번째 시험을 치를 수 있겠지. 그럼 네가 무사히 3차 전직을 할 수 있기를 빌겠네.")
				return
			end
			if is_second_class(me) == false or me:level() < 70 then
				me:dialog(npc, DEFAULT_TEXT)
				return
			end
			me:dialog(npc, "나의 분신인 만큼 그 강함은 상상할 수 없는 정도일 걸세. 각종 고급 기술을 사용하는 데다 자네 혼자서 1:1로 싸워야 하기 때문에 쉽지만은 않을거야. 게다가 그곳은 이 세계와는 다른 차원의 세계인 만큼 보통 인간이 오래 머무는 것은 좋지 않으므로 최대한 신속하게 쓰러뜨리는 것이 중요하다네. 미리 만반의 준비를 마친 후 도전하도록 하게나. 그럼 자네가 #b검은 부적#k을 가지고 무사히 돌아오기만을 기다리겠네.")
			return
		end
		if value == "job3_trial1_3" then
			me:dialog(npc, "내가 자네에게 주었던 그 목걸이를 가지고 오시리아의 #b로베이라#k님에게 돌아가면 두번째 시험을 치를 수 있겠지. 그럼 네가 무사히 3차 전직을 할 수 있기를 빌겠네.")
			return
		end

		me:dialog(npc, GREETINGS[me:class()] or DEFAULT_TEXT)
	end
}
