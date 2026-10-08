-- NPC name (String.wz/Npc.img.xml): 소공

local dojo = require("script/lib/dojo")

local function start(me, npc, group, party)
	local sm, err
	if party == nil then
		sm, err = group:start_solo(me)
	else
		sm, err = group:start_party(me, party)
	end
	if sm == nil then
		me:dialog(npc, "무릉도장이 꽉 찬 것 같아. 다음에 다시 도전해 봐.")
		if err ~= nil then
			log("mu_lung_dojo start:", err)
		end
		return false
	end
	return true
end

local function challenge_solo(me, npc)
	if me:party() ~= nil then
		me:dialog(npc, "파티에 가입되어 있는 것 같은데?")
		return
	end
	local group = state_machine(dojo.GROUP)
	local rest = dojo.record(me, dojo.REST_QUEST)
	local floor = 1
	if rest > 0 then
		floor = rest * 6
	end
	group:set_property("start:" .. me:id(), tostring(floor))
	if start(me, npc, group, nil) and rest > 0 then
		dojo.set_record(me, dojo.REST_QUEST, 0)
	end
end

local function challenge_party(me, npc)
	local party = me:party()
	if party == nil then
		me:dialog(npc, "파티에 가입되어 있긴 한거야?")
		return
	end
	if party:leader_id() ~= me:id() then
		me:dialog(npc, "음, 네가 파티장이 아닌 것 같은데?")
		return
	end
	local group = state_machine(dojo.GROUP)
	group:set_property("start:" .. party:id(), "1")
	start(me, npc, group, party)
end

local function claim_belt(me, npc)
	local points = dojo.record(me, dojo.POINT_QUEST)
	local claimed = dojo.record(me, dojo.BELT_QUEST)
	local options = {}
	for i, belt in ipairs(dojo.BELTS) do
		local suffix = ""
		if i <= claimed then
			suffix = " (획득)"
		end
		table.insert(options, string.format("#i%d:# #z%d#%s", belt.item, belt.item, suffix))
	end
	local sel = me:dialog_list(npc, string.format("당신의 수련 점수는 #b%d#k 점 이에요. 사부님은 재능 있는 사람들을 좋아하세요. 수련점수를 일정 수준 이상 획득하면, 점수에 따라 허리띠를 받으실 수 있답니다.", points), options)
	local belt = dojo.BELTS[sel]
	if belt == nil then
		return
	end
	if sel <= claimed then
		me:dialog(npc, "이미 획득한 허리띠 인것 같은데? 허리띠는 한번 밖에 받을 수 없어.")
		return
	end
	if sel ~= claimed + 1 then
		me:dialog(npc, "아직 하위 단계 허리띠를 받지 않은 것 같은데? 허리띠는 순차적으로 획득할 수 있어.")
		return
	end
	if points < belt.points or me:level() < belt.level then
		local text = string.format("#i%d:# #b#t%d##k를 받기 위해서는 #b%d 레벨#k 이상이어야 하며, 누적 수련점수 %d 점이 필요해.\r\n\r\n네가 이 허리띠를 받으려면 ", belt.item, belt.item, belt.level, belt.points)
		if points < belt.points then
			text = text .. string.format("수련점수를 #r%d#k점이나 더 쌓아야 한다고.", belt.points - points)
		else
			text = text .. string.format("레벨을 #r%d#k이나 더 올려야 한다고.", belt.level - me:level())
		end
		me:dialog(npc, text)
		return
	end
	if me:exchange(nil, { item = { [belt.item] = 1 } }) ~= ExchangeResult.OK then
		me:dialog(npc, "인벤토리에 충분한 공간이 있는지 확인해 줄래?")
		return
	end
	dojo.set_record(me, dojo.BELT_QUEST, sel)
end

local function reset_points(me, npc)
	if not me:dialog_yes_no(npc, "수련 점수를 초기화 하겠다고? 수련점수를 초기화 하면 지금까지 모아놨던 수련점수가 모두 사라져. 어때? 정말 초기화 하고 싶어?") then
		return
	end
	dojo.set_record(me, dojo.POINT_QUEST, 0)
	me:dialog(npc, "수련점수를 초기화 했어.")
end

local function show_record(me, npc)
	local floor = dojo.record(me, dojo.BEST_FLOOR_QUEST)
	if floor == 0 then
		me:dialog(npc, "아직 한 층도 통과하지 못했잖아? 기록을 남기고 싶으면 먼저 도전부터 해 봐.")
		return
	end
	local text = string.format("지금까지 #b%d층#k까지 통과했어.", floor)
	local seconds = dojo.record(me, dojo.BEST_TIME_QUEST)
	if seconds > 0 then
		text = text .. string.format("\r\n1층부터 옥상까지 가장 빨리 올라간 기록은 #b%s#k이야.", dojo.format_time(seconds))
	end
	me:dialog(npc, text)
end

local function spar(me, npc)
	local sm, err = state_machine(dojo.TUTORIAL_GROUP):start_solo(me)
	if sm ~= nil then
		return
	end
	me:dialog(npc, "지금은 대련할 수 없어. 다음에 다시 와.")
	if err ~= nil then
		log("mu_lung_dojo_tutorial start:", err)
	end
end

local function lobby(me, npc)
	local sel = me:dialog_list(npc, "우리 사부님은 무릉에서 최고로 강한 분이지. 그런 분에게 네가 도전하겠다고? 나중에 후회하지마.", {
		"혼자 도전해볼게.",
		"같이 도전해볼게.",
		"허리띠를 받고 싶어.",
		"수련점수를 초기화 할게.",
		"내 기록을 보고 싶어.",
		"너와 대련해 보고 싶어.",
		"무릉 도장이 뭐지?",
	})
	if sel == 1 then
		challenge_solo(me, npc)
	elseif sel == 2 then
		challenge_party(me, npc)
	elseif sel == 3 then
		claim_belt(me, npc)
	elseif sel == 4 then
		reset_points(me, npc)
	elseif sel == 5 then
		show_record(me, npc)
	elseif sel == 6 then
		spar(me, npc)
	elseif sel == 7 then
		me:dialog(npc, "우리 사부님은 무릉에서 가장 강한 분이야. 그런 사부님께서 만드신 곳이 바로 이 무릉 도장이라는 것이지. 무릉 도장은 38층이나 되는 높은 건물이야. 하나하나 올라가면서 자신을 수련할 수 있어. 물론 너의 실력으로는 끝까지 가기 힘들겠지만.")
	end
end

local function give_up(me, npc, sm)
	if not me:dialog_yes_no(npc, "결국 포기하는거야? 정말 나가겠어?") then
		return
	end
	if sm:leader() == me then
		sm:finish(dojo.EXIT)
		return
	end
	me:map(dojo.EXIT)
end

local function rest_floor(me, npc, sm, floor)
	local sel = me:dialog_list(npc, "여기까지 무사히 잘왔다니 놀랍네. 하지만 앞으로는 쉽지 않을걸? 어때 계속 도전해 보겠어?", {
		"계속 도전해볼게.",
		"밖으로 나가겠어.",
		"진행 상황을 저장하겠어.",
	})
	if sel == 1 then
		if sm:leader() ~= me then
			me:dialog(npc, "파티장만 도전을 신청할 수 있어.")
			return
		end
		dojo.advance(sm, floor)
	elseif sel == 2 then
		give_up(me, npc, sm)
	elseif sel == 3 then
		if sm:party() ~= nil then
			me:dialog(npc, "같이 도전중일때는 진행상황을 저장할 수 없어.")
			return
		end
		dojo.set_record(me, dojo.REST_QUEST, math.floor(floor / 6))
		me:dialog(npc, "진행상황이 저장되었어. 다음에 다시 도전할때, 이 층으로 바로 올라올 수 있을거야.")
	end
end

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		local floor = dojo.current_floor(me)
		if sm == nil or floor == nil then
			lobby(me, npc)
		elseif dojo.is_rest(floor) then
			rest_floor(me, npc, sm, floor)
		else
			give_up(me, npc, sm)
		end
	end
}
