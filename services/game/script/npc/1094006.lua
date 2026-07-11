-- NPC name (String.wz/Npc.img.xml): 풀 숲

function on_click(me, npc)
	local q = me:quest(2186)
	if q == nil or not q:started() then
		return
	end
	local function item_count(item_id)
		local count = 0
		for _, it in pairs(me:item(item_id)) do
			count = count + it:count()
		end
		return count
	end
	local rand = math.random(0, 100)
	if rand < 30 and item_count(4031853) < 1 then
		me:exchange(nil, { item = { [4031853] = 1 } })
		me:dialog(npc, "아벨의 안경을 찾았다.")
	elseif rand < 60 and item_count(4031854) < 1 then
		me:exchange(nil, { item = { [4031854] = 1 } })
		me:dialog(npc, "안경을 찾았다. 하지만 아벨의 안경이 아닌듯 하다. 아벨의 안경은 검은 뿔테라는데...")
	elseif rand <= 100 and item_count(4031855) < 1 then
		me:exchange(nil, { item = { [4031855] = 1 } })
		me:dialog(npc, "안경을 찾았다. 하지만 아벨의 안경이 아닌듯 하다. 아벨의 안경은 검은 뿔테라는데...")
	end
end
