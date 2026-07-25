local quest_id = 2156

local function morph_source(me)
	local b = me:buff(BuffFlag.Morph)
	if not b then return -1 end
	local wz = b:wz()
	if type(wz.id) == "number" then return wz.id end
	return wz:id()
end

local function has_item(me, item_id, need)
	local slots = me:item(item_id)
	local count = 0
	for _, it in pairs(slots) do
		count = count + it:count()
	end
	return count >= need
end

return {
	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not q:started() then
			q:start(npc, true)
			return
		end

		local file = "#fUI/UIWindow.img/QuestIcon/"
		if not has_item(me, 2210006, 1) then
			me:dialog(npc, "흠... #b#t2210006##k은 아직 구하지 못한 거야? 해안가 풀숲으로 가서 #b#o2220000##k를 잡아 보라니까?")
			return
		end

		me:dialog(npc, "좋아. 이게 바로 무지개색 달팽이 껍질이라 이거지? \r\n\r\n" .. file .. "4/0#\r\n\r\n\r\n" .. file .. "8/0# 7500 exp", false, true)

		local is_morph = morph_source(me) == 2210006
		local code
		if is_morph then
			code = me:exchange(
				{ item = { [2210006] = 1 }, population = 1 },
				{ meso = 30000, exp = 10000 }
			)
		else
			code = me:exchange(
				{ item = { [2210006] = 1 } },
				{ meso = 30000, exp = 7500, population = 3 }
			)
		end
		if code == ExchangeResult.LackCapacity then
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end

		if is_morph then
			me:dialog(npc, "뭐, 뭐야 그 모습은? 왜 달팽이가 되어서 온 거야?! 뭐? 무지개색 달팽이 껍질을 쓰니까 이렇게 되었다고? 그러고 보니 달팽이 전설에서도 보물의 위험에 대해 경고하는 말이 있었지.. 휴우. 다행이다. 안 먹어서. 욕심을 부리니까 그런 꼴이 된 거라고.", false, true)
		else
			me:dialog(npc, "무지개색 달팽이 껍질을 나누기로 하지 않았냐고? 하지만 이걸 반으로 갈랐다가는 효능이 없어질지도 모르잖아? 먼저 정보를 준 건 이 쪽이니까, 당연히 내가 가져야지! 후훗!", false, true)
			local map = me:map()
			if map ~= nil then
				for _, n in pairs(map:npcs()) do
					if n:id() == npc then
						n:show_effect("act2156")
						break
					end
				end
			end
		end

		me:show_effect(EffectType.QuestCompletion)
		q:force_complete(npc)
	end
}
