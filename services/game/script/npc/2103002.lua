-- NPC name (String.wz/Npc.img.xml): 왕비의 장식장

function on_click(me, npc)
	local q = me:quest(3923)
	if q == nil or not q:started() then
		return
	end
	local count = 0
	for _, it in pairs(me:item(4031578)) do
		count = count + it:count()
	end
	if count >= 1 then
		return
	end
	me:chat('asd')
	me:exchange(nil, { item = { [4031578] = 1 } })
end
