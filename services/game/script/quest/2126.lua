local quest_id = 2126

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not q:started() then
		q:start(npc, true)
		return
	end

	local file = "#fUI/UIWindow.img/QuestIcon/"
	local q3937 = me:quest(3937)
	if q3937 ~= nil and q3937:completed() then
		me:dialog(npc, "아! 누군가 했더니 오랜만이야. 이번엔 보급품 운반을 맡았나보지? 꽤 중요한 임무였는데. 수고했어.\r\n\r\n" .. file .. "4/0#\r\n#v2030000# #t2030000# 5개\r\n#v2022155# #t2022155# 5개\r\n\r\n" .. file .. "8/0#\r\n2000 exp", false, true)
		local code = me:exchange(
			{ item = { [4031624] = 1 } },
			{ item = { [2030000] = 5, [2022155] = 5 }, exp = 2000 }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "소비창에 빈 칸이 있는지 확인해 주세요.")
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		q:force_complete(npc)
	else
		me:dialog(npc, "수고했어요. 보는 눈이 많으니까 그만 가보도록해요.\r\n\r\n" .. file .. "4/0#\r\n#v2030000# #t2030000# 5개\r\n#v2022155# #t2022155# 5개\r\n\r\n" .. file .. "8/0#\r\n2000 exp", false, true)
		local code = me:exchange(
			{ item = { [4031624] = 1 } },
			{ item = { [2030000] = 10 }, exp = 2000 }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "소비창에 빈 칸이 있는지 확인해 주세요.")
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		q:force_complete(npc)
	end
end
