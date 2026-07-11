local quest_id = 3250

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "우왓! 이 녀석은 먹였고, 이 녀석도 먹였고, 이 녀석도 먹였던가? 아냐! 3번 녀석은 안 먹였어! 자, 어서 먹이를 먹어랏! 휴우~ 다들 입을 벌리고 있으니 어느 녀석이 배가 고픈지 모르겠네... 아, 왔어?", false, true)
	me:dialog(npc, "휴우... 바빠서 정신이 하나도 없어. 얼마 전에 연구원들이 말하길, #o5220003#는 파풀라투스의 지배를 받아 변해 버렸을 뿐 잘만 기르면 사람을 잘 따르는 귀여운 새가 된다지 뭐야. 그래서 한가한 틈을 타서 새를 기르기로 했어.", false, true)
	if not me:dialog_yes_no(npc, "그런데 이거... 한 두마리도 아니고, 열 마리를 동시에 기르려니 무진장 피곤해... 먹이 구하는 것도 쉽지 않고... 뭐, 그래도 귀엽긴 하지만. 너도 한 마리 길러볼래? 분양해 줄게.") then
		me:dialog(npc, "흠... 동물은 별로 좋아하지 않는 모양이구나. 이렇게 귀여운데... 이 녀석들에게라면 햄버거도 양보해 줄 수 있다구.", false, false)
		return
	end

	me:dialog(npc, "#o5220003#는 다른 세계의 생물이라, 다 기르면 원래 세계로 돌려보내줘야 해. 그러니 #o5220003#가 #b완전히 자라거든 다시 데려다 줘#k. 그럼 부탁할게.", false, true)

	local code = me:exchange({}, { item = { [4220046] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "소비창에 빈 칸이 있는지 확인해 주세요.")
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	local q7067 = me:quest(7067)
	if q7067 ~= nil then
		if q7067:wz() == nil then
			q7067:start("0")
		else
			q7067:start(npc, "0")
		end
	end
	q:start(npc, true)
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	local qr = me:quest(7067)
	if qr == nil then
		return
	end
	if qr:record() ~= "3000" then
		me:dialog(npc, "아직 타이머를 다 못키운거야? 시계탑 몬스터들을 사냥하면 나오는 태엽벌레를 타이머에게 먹여봐.", false, false)
		return
	end

	local file = "#fUI/UIWindow.img/QuestIcon/"
	me:dialog(npc, "어때? 타이머 기르는 건 재미있어?", false, true)
	me:dialog(npc, "응? 타이머가 벌써 다 자라 버렸다고? 헉... 태엽벌레를 열심히 먹인 모양이구나... 그럼 이제 그만 타이머를 돌려줘. 이제 그만 녀석을 원래세계로 보내줘야지...\r\n\r\n" .. file .. "4/0#\r\n\r\n" .. file .. "6/0# 11", false, true)

	local code = me:exchange({ item = { [4220046] = 1 } }, { population = 11 })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "소비창에 빈 칸이 있는지 확인해 주세요.")
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
	me:dialog(npc, "아쉽기는 하지만 이 녀석들은 원래 다른 차원의 생물이니까 원래 세계가 녀석들에게 살기 좋을거야.", false, false)
end
