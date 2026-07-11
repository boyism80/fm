-- 칭호 도전 - 인기인! (Quest.wz/QuestData/29002.img): 칭호 도전 - 인기인!

local quest_id = 29002

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "#v1142003# #e#b#t1142003##k\r\n\r\n - 제한시간 30일\r\n - 인기도 1000상승\r\n\r\n#n이 훈장의 주인이 될 자격이 있는지 시험해 보시겠소?", false, true)
	me:dialog(npc, "자, 30일의 시간을 줄테니 목적을 이루고 나에게 돌아오시오. 제한시간 내에 나에게 와서 확인을 받아야만 인정받을 수 있다는 것을 꼭 명심하시오. 그리고 이 도전을 완료하거나 포기하지 않는 이상 다른 칭호에 도전할 수는 없다는 것도 알아두시오.", false, true)
	local deadline = now() * 1000 + 86400000 * 30
	q:start(npc, "time_" .. deadline)
	q:record_ex("popG", "1000")
	q:record_ex("popS", tostring(me:population()))
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	local pop_start = tonumber(q:record_ex("popS")) or 0
	local add_fame = me:population() - pop_start
	if add_fame < 1000 then
		if me:dialog_accept(npc, "그대의 인기도는 그동안 " .. add_fame .. "만큼 올랐소. 목표인 1000을 달성하지 못했으니 이 칭호를 받기에는 무리라고 생각하오만... 그만 포기하시겠소?") then
			q:forfeit()
			me:dialog(npc, "다음에 도전하려면 언제든지 또 오시오.", false, false)
		else
			me:dialog(npc, "계속 도전을 해주시오.", false, false)
		end
		return
	end

	local record = q:record()
	local deadline = tonumber(record and record:sub(6))
	if deadline ~= nil and deadline < now() * 1000 then
		me:dialog(npc, "이미 30일의 제한 시간이 지난 것 같소. 다시 도전하시오.", false, false)
		q:forfeit()
		return
	end

	me:dialog(npc, "호오, 대단하오. 인기도의 훈장을 받을 자격이 충분하구려.", false, true)
	local code = me:exchange({}, { item = { [1142003] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "장비 인벤토리 공간이 충분한지 확인해주게.", false, false)
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
end
