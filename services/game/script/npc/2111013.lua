-- NPC name (String.wz/Npc.img.xml): 액자

local function distance_sq(me, npc_id)
	local map = me:map()
	if map == nil then
		return nil
	end
	local nx, ny
	for _, n in pairs(map:npcs()) do
		if n:id() == npc_id then
			nx, ny = n:position()
			break
		end
	end
	if nx == nil then
		return nil
	end
	local mx, my = me:position()
	local dx = mx - nx
	local dy = my - ny
	return dx * dx + dy * dy
end

function on_click(me, npc)
	local dist = distance_sq(me, npc)
	if dist ~= nil and dist > 7000 then
		me:dialog(npc, "너무 거리가 멀어 조사할 수 없다.", false, false)
		return
	end
	local q = me:quest(3322)
	if q ~= nil and q:started() then
		local count = 0
		for _, it in pairs(me:item(4031697)) do
			count = count + it:count()
		end
		if count < 1 then
			me:dialog(npc, "액자 뒤에 있는 고리를 풀어 열었다. 액자 속의 비밀 공간에 은색 펜던트가 들어 있다. 조심스럽게 펜던트를 꺼낸 후 액자를 닫고 테이블 위에 놓았다.", false, false)
			me:exchange(nil, { item = { [4031697] = 1 } })
			return
		end
	end
	me:dialog(npc, "아무것도 없는 것 같다.", false, false)
end
