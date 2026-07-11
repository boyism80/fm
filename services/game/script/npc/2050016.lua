-- NPC name (String.wz/Npc.img.xml): 운석3

function on_click(me, npc)
	local npc_id = npc
	if type(npc) ~= "number" then
		npc_id = npc:id()
	end
	local q = me:quest(3421)
	if q == nil or not q:started() then
		return
	end
	local nCount = npc_id - 2050014
	local info = q:record()
	if info == nil or info == "" then
		info = "000000"
	end
	if info == "111111" then
		return
	end
	local nInvite = info:sub(nCount + 1, nCount + 1)
	if nInvite == "1" then
		me:dialog(npc, "이 운석의 샘플은 이미 채취했다. 다른 운석을 찾아보자.")
		return
	end
	local code = me:exchange(nil, { item = { [4031117] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "인벤토리 공간이 부족한 것 같다.")
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	me:dialog(npc, "운석 샘플을 채취했다.")
	local z = ""
	for i = 0, 5 do
		if nCount == i then
			z = z .. "1"
		else
			z = z .. info:sub(i + 1, i + 1)
		end
	end
	q:record(z)
	q:sync_progress()
end
