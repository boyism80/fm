-- NPC name (String.wz/Npc.img.xml): 러셀론의 책상

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

return {
	on_click = function(me, npc)
		local dist = distance_sq(me, npc)
		if dist ~= nil and dist > 7000 then
			me:dialog(npc, "조사하기에는 너무 멀다.", false, false)
			return
		end
		local q = me:quest(3314)
		if q ~= nil and q:started() then
			local count = 0
			for _, it in pairs(me:item(2022198)) do
				count = count + it:count()
			end
			if count < 1 then
				me:dialog(npc, "책상 위에는 작은 알약들이 여러 개 놓여 있다. 한 개만 가져 가자...", false, false)
				me:exchange(nil, { item = { [2022198] = 1 } })
				return
			end
		end
		me:dialog(npc, "책상 위에는 작은 알약들이 여러 개 놓여 있다.", false, false)
	end
}
