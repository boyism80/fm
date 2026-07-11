-- NPC name (String.wz/Npc.img.xml): 알카드노의 책장

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
	if dist ~= nil and dist > 5000 then
		me:notice("너무 멀어 조사할 수 없다.", 5)
		return
	end
	local q = me:quest(3309)
	if q == nil or not q:started() then
		return
	end
	local count = 0
	for _, it in pairs(me:item(4031708)) do
		count = count + it:count()
	end
	if count >= 1 then
		return
	end
	me:dialog(npc, "어둠 속에서 책장이 만져진다... 눈에 힘을 주고 잘 살펴보자, 이상한 서류 뭉치가 보인다... 이게 바로 베딘이 말한 그 서류인 것 같다. 서류를 챙겼으니 베딘에게 돌아가자.", false, true)
	me:exchange(nil, { item = { [4031708] = 1 } })
end
