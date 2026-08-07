-- NPC name (String.wz/Npc.img.xml): 시종 이크

local pq = require("script/lib/party_quest")

local STAGE_HELP = "이곳은 여신의 탑의 휴게실입니다. 미네르바 여신님은 이곳에서 음악을 즐기곤 하셨지요. 여신님은 매일매일 다른 음악을 즐기셨답니다.\r\n\r\n여신님께서 즐기신 음악은 \r\n\r\n#b월요일에는 아기자기한 음악\r\n화요일에는 무서운 음악\r\n수요일에는 재미있는 음악\r\n목요일에는 우울한 음악\r\n금요일에는 싸늘한 음악\r\n토요일에는 깔끔한 음악\r\n일요일에는 웅장한 음악\r\n\r\n#k여신님은 이렇게 매일매일 음악을 바꾸어 들으셨지요.\r\n\r\n만약 그때 그 음악을 틀어주신다면 미네르바 여신님의 영혼이 신비한 일을 일으킬지도 모르겠네요."
local STAGE4_HELP = "이곳은 여신의 탑의 봉인된 방입니다. 여신 미네르바께서 자신의 소중한 물건을 안전하게 보관할 수 있을 만큼 안전한 방입니다. 그래서, 이 방을 열기 위해선 특수한 무게가 필요합니다. 여러분 다섯명이 발판을 올바르게 올라서서 정해진 무게를 맞추어야 합니다. 7번의 기회를 드리며, 그 안에 맞추시지 못한다면 안전장치로 봉인된 방에서 추방되게 됩니다."

local function clear_fx(map)
	if map == nil then
		return
	end
	map:show_effect("quest/party/clear")
	map:play_sound("Party1/Clear")
end

local function strip_pieces(me)
	for id = 4001044, 4001063 do
		pq.remove_all(id, me)
	end
end

local function ensure_stage4_rand(sm)
	local quiz = sm:get_property("stage4rand")
	if quiz ~= "" and quiz ~= "0" then
		return quiz
	end
	local a = math.random(0, 3)
	local b = math.random(0, 3 - a)
	local c = 5 - a - b
	local parts = { a, b, c }
	local order = { 1, 2, 3 }
	for i = #order, 2, -1 do
		local j = math.random(i)
		order[i], order[j] = order[j], order[i]
	end
	quiz = tostring(parts[order[1]]) .. tostring(parts[order[2]]) .. tostring(parts[order[3]])
	sm:set_property("stage4rand", quiz)
	sm:set_property("stage4try", "0")
	return quiz
end

local function handle_exit(me)
	strip_pieces(me)
	me:map(200080101)
end

local function handle_prestage(me, npc, sm)
	strip_pieces(me)
	local state = sm:get_property("prestage")
	if state == "" then
		me:dialog(npc, "안녕하세요, 저는 여신을 모시는 시종 이크라고 해요. 지금 제 모습이 보이지 않으셔도 놀라지 마세요. 여신께서 석상으로 변화되어 버렸을때 저는 제 힘을 모두 잃어버리고 말았어요. 오르비스의 구름조각을 구해오신다면 제 몸을 복구하고 여러분들 앞에 나타날 수 있답니다. #b20#k개의 구름 조각을 모아 제게 가져오시면 된답니다. 그리고, 지금 제 모습은 무척 작아서 반짝거리는 빛으로 밖에 보이지 않으실 거에요.")
		return
	end
	if state ~= "clear" then
		return
	end
	local rewarded = sm:get_property("prestageRewarded")
	if not pq.is_leader(me) then
		if rewarded == "clear" then
			me:map(920010000, 2)
		else
			me:dialog(npc, "제 몸을 다시 만들어 주셔서 너무 고마워요! 여신의 탑으로 들어가려면 파티장이 제게 말을 걸어주면 된답니다.")
		end
		return
	end
	if rewarded == "clear" then
		me:map(920010000, 2)
		return
	end
	local map = me:map()
	clear_fx(map)
	me:dialog(npc, "제 몸을 다시 만들어 주셔서 너무 고마워요! 입구로 데려다 드릴게요~!")
	sm:set_property("prestageRewarded", "clear")
	pq.party_exp(sm, 16000)
	pq.party_warp(sm, 920010000, nil, 2)
end

local function handle_hub(me, npc, sm)
	local status = tonumber(sm:get_property("status")) or 0
	if status < 6 then
		me:dialog(npc, "여신상의 조각을 모아 여신님의 석상을 복구하고 여신님을 구해주세요!")
		return
	end
	if status == 6 then
		if sm:get_property("stage0clear") == "" then
			me:dialog(npc, "여신님의 석상이 복구되었군요! 이제 여신님을 구하기 위한 마지막 단계만이 남았어요! 이 모든 일의 원흉을 처치하고 이 여신의 탑의 평화를 되찾아 주세요.")
			sm:set_property("stage0clear", "clear")
			pq.party_warp(sm, 920010800, nil, 1)
		else
			me:dialog(npc, "얻어오신 생명의 풀로 여신님을 구해주세요!")
		end
		return
	end
	if status == 7 then
		me:dialog(npc, "여신님을 마침내 구해주셨군요! 여신님께 말을 걸어보세요. 분명 좋은 보상을 주실거에요.")
	end
end

local function handle_stage1(me, npc, sm, map)
	local stage = sm:get_property("stage1clear")
	if stage == "clear" then
		me:dialog(npc, "여러분은 이 방을 훌륭하게 클리어 하셨습니다. 다음 방으로 이동해 주세요.")
		return
	end
	if not pq.is_leader(me) then
		me:dialog(npc, "이곳은 여신의 탑의 산책로 입니다. 여신님은 이곳에서 산책을 즐기곤 하셨지요. 여러분은 #b#t4001050##k 30개를 모아서 제게 가져오시면 제가 조각들을 모아 #b#t4001044##k으로 바꿔드릴게요. 그럼 힘내주세요!")
		return
	end
	if stage == "" then
		sm:set_property("stage1clear", "s")
		me:dialog(npc, "이곳은 여신의 탑의 산책로 입니다. 여신님은 이곳에서 산책을 즐기곤 하셨지요. 여러분은 #b#t4001050##k 30개를 모아서 제게 가져오시면 제가 조각들을 모아 #b#t4001044##k으로 바꿔드릴게요. 그럼 힘내주세요!")
		return
	end
	if pq.has_item(me, 4001050, 30) then
		local code = me:exchange(
			{ item = { [4001050] = pq.item_count(me, 4001050) } },
			{ item = { [4001044] = 1 } }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "#b#t4001050##k 을 30개 찾아 제게 가져오시면 된답니다. 혹은 인벤토리 공간이 부족한건 아닌지 확인해 주세요.")
			return
		end
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "#b#t4001050##k 을 30개 찾아 제게 가져오시면 된답니다. 혹은 인벤토리 공간이 부족한건 아닌지 확인해 주세요.")
			return
		end
		clear_fx(map)
		sm:set_property("stage1clear", "clear")
		pq.party_exp(sm, 17500)
		me:dialog(npc, "훌륭해요! #b#t4001050##k을 30개 찾아오셨군요!")
		return
	end
	me:dialog(npc, "#b#t4001050##k 을 30개 찾아 제게 가져오시면 된답니다. 혹은 인벤토리 공간이 부족한건 아닌지 확인해 주세요.")
end

local function handle_stage2(me, npc, sm, map)
	local stage = sm:get_property("stage2clear")
	if stage == "clear" then
		me:dialog(npc, "여러분은 이 방을 훌륭하게 클리어 하셨습니다. 다음 방으로 이동해 주세요.")
		return
	end
	if not pq.is_leader(me) then
		me:dialog(npc, "이곳은 여신의 탑의 창고입니다. 하지만 지금은 샐리온의 소굴이 되고 말았답니다. 샐리온은 #b#t4001045##k을 가지고 다니고 있답니다. 샐리온을 처치하시다가 #b#t4001045##k을 찾아 제게 가져 오시면 된답니다!")
		return
	end
	if stage == "" then
		sm:set_property("stage2clear", "s")
		me:dialog(npc, "이곳은 여신의 탑의 창고입니다. 하지만 지금은 샐리온의 소굴이 되고 말았답니다. 샐리온은 #b#t4001045##k을 가지고 다니고 있답니다. 샐리온을 처치하시다가 #b#t4001045##k을 찾아 제게 가져 오시면 된답니다!")
		return
	end
	if pq.has_item(me, 4001045, 1) then
		clear_fx(map)
		sm:set_property("stage2clear", "clear")
		pq.party_exp(sm, 17500)
		me:dialog(npc, "훌륭해요! #b#t4001045##k을 찾아오셨군요!")
		return
	end
	me:dialog(npc, "#b#t4001045##k 을 찾아 제게 가져오시면 된답니다.")
end

local function handle_stage3(me, npc, sm, map)
	if not pq.is_leader(me) then
		me:dialog(npc, STAGE_HELP)
		return
	end
	local music = sm:get_property("stage3_music")
	if music == "" then
		me:dialog(npc, STAGE_HELP)
		return
	end
	if sm:get_property("stage3clear") == "clear" then
		me:dialog(npc, "여신님께서 즐겨들으시던 음악이 나오고 있어요.")
		return
	end
	local today = os.date("*t").wday
	if tonumber(music) == today then
		clear_fx(map)
		local box = map:reactor(2002011)
		if box ~= nil then
			box:hit(1)
		end
		sm:set_property("stage3clear", "clear")
		pq.party_exp(sm, 17500)
		me:dialog(npc, "바로 이 음악이에요! 여신님은 이 날에 이 음악을 즐겨 듣곤 하셨지요.")
		return
	end
	me:dialog(npc, "이 음악이 아닌 것 같은데요...")
end

local function handle_stage4(me, npc, sm, map)
	local quiz = ensure_stage4_rand(sm)
	local clear = sm:get_property("stage4clear")
	if clear == "clear" then
		me:dialog(npc, "이미 이 미션을 훌륭하게 클리어 하신 것 같아요. 다른 방에 도전해 보세요!")
		return
	end
	if not pq.is_leader(me) then
		me:dialog(npc, STAGE4_HELP)
		return
	end
	local a1 = map:players_in_area(0)
	local a2 = map:players_in_area(1)
	local a3 = map:players_in_area(2)
	if a1 + a2 + a3 ~= 3 then
		me:dialog(npc, "아직 파티원 3명이 정답 발판 3개를 찾지 못한것 같군요. 아슬아슬하게 서 계시지 말고 가운데에 올바르게 서 계셔야 정답으로 인정되니 주의해 주세요!")
		return
	end
	local ans1 = tonumber(quiz:sub(1, 1)) or -1
	local ans2 = tonumber(quiz:sub(2, 2)) or -1
	local ans3 = tonumber(quiz:sub(3, 3)) or -1
	if a1 == ans1 and a2 == ans2 and a3 == ans3 then
		clear_fx(map)
		local box = map:reactor(2002012)
		if box ~= nil then
			box:hit(1)
		end
		sm:set_property("stage4clear", "clear")
		pq.party_exp(sm, 17500)
		me:dialog(npc, "축하해요! 정답을 맞추셨어요!")
		return
	end
	local tries = (tonumber(sm:get_property("stage4try")) or 0) + 1
	sm:set_property("stage4try", tostring(tries))
	if tries == 6 then
		me:dialog(npc, "마지막 기회에요! 틀리지 않게 조심해 주세요! 행운을 빌어요~")
		return
	end
	if tries == 7 then
		sm:notice("봉인된 방에서 추방되었습니다.")
		me:mkitem(4001047, 1)
		sm:set_property("stage4rand", "0")
		pq.party_warp(sm, 920010100, nil, 13)
		return
	end
	local ok = 0
	if a1 == ans1 then
		ok = ok + 1
	end
	if a2 == ans2 then
		ok = ok + 1
	end
	if a3 == ans3 then
		ok = ok + 1
	end
	me:dialog(npc, "정답이 아니에요! 지금까지 " .. tostring(tries) .. " 번 도전 하셨어요. 총 7번의 기회가 있으니 천천히 도전해보세요~\r\n\r\n#e파티원들이 " .. tostring(ok) .. " 개의 발판에 올바르게 올라섰군요.#n")
end

local function handle_stage5(me, npc, sm, map)
	local stage = sm:get_property("stage5clear")
	if stage == "clear" then
		me:dialog(npc, "여러분은 이 방을 훌륭하게 클리어 하셨습니다. 다음 방으로 이동해 주세요.")
		return
	end
	if not pq.is_leader(me) then
		me:dialog(npc, "이곳은 여신의 탑의 손님들이 머물던 숙박실 입니다. 여러분은 40개의 #b#t4001052##k을 모아 제게 가져오시면 됩니다. 그럼 행운을 빌어요!")
		return
	end
	if stage == "" then
		sm:set_property("stage5clear", "s")
		me:dialog(npc, "이곳은 여신의 탑의 손님들이 머물던 숙박실 입니다. 여러분은 40개의 #b#t4001052##k을 모아 제게 가져오시면 됩니다. 그럼 행운을 빌어요!")
		return
	end
	if pq.has_item(me, 4001052, 40) then
		local code = me:exchange(
			{ item = { [4001052] = pq.item_count(me, 4001052) } },
			{ item = { [4001048] = 1 } }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "#b#t4001052##k 을 40개 찾아 제게 가져오시면 된답니다. 혹은 인벤토리 공간이 부족한건 아닌지 확인해 주세요.")
			return
		end
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "#b#t4001052##k 을 40개 찾아 제게 가져오시면 된답니다. 혹은 인벤토리 공간이 부족한건 아닌지 확인해 주세요.")
			return
		end
		clear_fx(map)
		sm:set_property("stage5clear", "clear")
		pq.party_exp(sm, 17500)
		me:dialog(npc, "훌륭해요! #b#t4001052##k을 40개 찾아오셨군요!")
		return
	end
	me:dialog(npc, "#b#t4001052##k 을 40개 찾아 제게 가져오시면 된답니다. 혹은 인벤토리 공간이 부족한건 아닌지 확인해 주세요.")
end

local function handle_stage6(me, npc, sm, map)
	local stage = sm:get_property("stage6clear")
	if stage == "clear" then
		me:dialog(npc, "여러분은 이 방을 훌륭하게 클리어 하셨습니다. 다음 방으로 이동해 주세요.")
		return
	end
	if not pq.is_leader(me) then
		me:dialog(npc, "이곳은 여신의 탑의 꼭대기로 올라가는 곳입니다. 맵 위쪽에 문을 제어하는 다섯개의 레버가 있습니다. 레버 중 두개의 레버만 문을 제어하는 레버이니 레버를 두개 당기신 후 제게 오시면 정답인지 아닌지 확인해 드리겠습니다. 그럼 힘내주세요! ")
		return
	end
	if stage == "" then
		sm:set_property("stage6clear", "s")
		sm:set_property("stage6_ans", pq.shuffle("11000"))
		me:dialog(npc, "이곳은 여신의 탑의 꼭대기로 올라가는 곳입니다. 맵 위쪽에 문을 제어하는 다섯개의 레버가 있습니다. 레버 중 두개의 레버만 문을 제어하는 레버이니 레버를 두개 당기신 후 제게 오시면 정답인지 아닌지 확인해 드리겠습니다. 그럼 힘내주세요! ")
		return
	end
	local cur = ""
	for i = 1, 5 do
		local lever = map:reactor_by_name(tostring(i))
		local state = 0
		if lever ~= nil then
			state = lever:state()
		end
		cur = cur .. tostring(state)
	end
	if cur == sm:get_property("stage6_ans") then
		local box = map:reactor(2002013)
		if box ~= nil then
			box:hit(1)
		end
		clear_fx(map)
		sm:set_property("stage6clear", "clear")
		pq.party_exp(sm, 17500)
		me:dialog(npc, "축하합니다. 정답 레버를 찾으셨군요! 나타난 상자에서 여신의 조각상 조각을 회수해 주세요~")
		return
	end
	me:dialog(npc, "정답 레버가 아닙니다. 다시 한번 시도해 보세요.")
end

local function handle_garden(me, npc, sm)
	local boss = sm:get_property("stagebossclear")
	if boss == "clear" then
		me:dialog(npc, "얻은 생명의 풀로 여신님을 되살려 주세요!")
		return
	end
	if pq.is_leader(me) and boss == "" then
		sm:set_property("stagebossclear", "s")
	end
	me:dialog(npc, "이곳은 여신의 탑의 정원입니다. #b생명의 풀#k은 이곳에 있는 화분에서만 키울 수 있답니다. 생명의 풀을 기를 씨앗만 있으면 될텐데... 참! 여신님을 석상에 가둔 #b파파픽시#k가 이곳에 있을지 모르니 조심하세요!")
end

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz().id
		if map_id == 920011200 then
			handle_exit(me)
			return
		end
		local sm = me:state_machine()
		if sm == nil then
			me:dialog(npc, "잠시 후 다시 시도해 주세요.")
			return
		end
		if map_id == 920010000 then
			handle_prestage(me, npc, sm)
		elseif map_id == 920010900 then
			me:dialog(npc, "이곳은 여신 미네르바께서 죄인들을 가두어 두었던 지하 감옥입니다. 이곳에 여신의 조각상 조각은 없답니다. 하지만 무언가가 숨겨져 있을지도 모르겠군요.")
		elseif map_id == 920010100 then
			handle_hub(me, npc, sm)
		elseif map_id == 920010200 then
			handle_stage1(me, npc, sm, map)
		elseif map_id == 920010300 then
			handle_stage2(me, npc, sm, map)
		elseif map_id == 920010400 then
			handle_stage3(me, npc, sm, map)
		elseif map_id == 920010500 then
			handle_stage4(me, npc, sm, map)
		elseif map_id == 920010600 then
			handle_stage5(me, npc, sm, map)
		elseif map_id == 920010700 then
			handle_stage6(me, npc, sm, map)
		elseif map_id == 920010800 then
			handle_garden(me, npc, sm)
		elseif not pq.is_leader(me) then
			me:dialog(npc, "파티장이 제게 말을 걸어 주어야 합니다.")
		end
	end
}
