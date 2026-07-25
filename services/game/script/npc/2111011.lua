-- NPC name (String.wz/Npc.img.xml): 벽

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
		if dist ~= nil and dist > 5000 then
			me:dialog(npc, "너무 거리가 멀어 조사할 수 없다.", false, false)
			return
		end
		local q = me:quest(3311)
		if q == nil or not q:started() then
			return
		end
		if not me:dialog_yes_no(npc, "거미줄 틈으로 보이는 벽에 뭔가 글자들이 보이는 것 같다... 벽을 살피시겠습니까?") then
			return
		end
		q:record("5")
		q:sync_progress()
		me:show_quest_completion(3311)
		me:dialog(npc, "지저분한 낙서들 틈으로 유난히 선명하게 보이는 글자가 있다.  #b그것은 펜던트의 형태로 완성되었다...#k  무슨 말일까?", false, true)
	end
}
