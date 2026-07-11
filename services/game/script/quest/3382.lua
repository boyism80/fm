-- Quest name (Quest.wz/Quest.img.xml): 유레테의 보답

local quest_id = 3382

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not q:started() then
		q:start(npc, true)
		return
	end

	local sel = me:dialog_list(npc, "제뉴미스트의 구슬과 알카드노의 구슬을 각각 25개씩 가져오면 호루스의 눈을, 둘 중 어느쪽의 구슬이라도 10개 가져오면 목걸이에 새로운 힘을 부여할 수 있는 지혜의 돌을 만들어 주겠네. 자, 어떤 아이템을 만들겠는가?\r\n#b", {
		"호루스의 눈을 만들어 주세요.",
		"제뉴미스트 구슬로 지혜의 돌을 만들어 주세요.",
		"알카드노 구슬로 지혜의 돌을 만들어 주세요.",
	})
	if sel == nil then
		return
	end

	local material
	local reward_item
	if sel == 0 then
		material = { [4001159] = 25, [4001160] = 25 }
		reward_item = 1122010
	elseif sel == 1 then
		material = { [4001159] = 10 }
		reward_item = 2041212
	else
		material = { [4001160] = 10 }
		reward_item = 2041212
	end

	local str = "교환할 아이템을 확인하게.\r\n\r\n"
	if sel == 0 then
		str = str .. "#i4001159# #t4001159# 25개\r\n"
		str = str .. "#i4001160# #t4001160# 25개\r\n"
	elseif sel == 1 then
		str = str .. "#i4001159# #t4001159# 10개\r\n"
	else
		str = str .. "#i4001160# #t4001160# 10개\r\n"
	end
	str = str .. "\r\n위 아이템으로 \r\n"
	str = str .. string.format("#i%d# #t%d# 1개로 교환하겠네. 계속하겠는가?", reward_item, reward_item)

	if not me:dialog_yes_no(npc, str) then
		me:dialog(npc, "흐음, 잘못 선택한건가? 마음이 바뀌면 다시 오게. 아마 오래 만나지 못할게야.", false, false)
		return
	end

	local code = me:exchange({ item = material }, { item = { [reward_item] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "인벤토리 공간이 부족한건 아닌지, 혹은 재료를 분명 제대로 갖고 계신건지 확인해 주시게.", false, false)
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end

	if reward_item == 1122010 then
		local q1205 = me:quest(1205)
		if q1205 ~= nil and q1205:record_ex("have") == nil then
			q1205:record_ex("have", "1")
		end
	end
	q:force_complete(npc)
end
