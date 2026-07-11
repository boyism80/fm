local quest_id = 2317

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_accept(npc, "Ah! If I am not mistaken, I saw the #bKiller Mushroom Spores#k way back when I was a kid in a book. Now I remember... it's made out of extracts of powerful poisons from Poison Mushrooms, which means you'll need some Poison Mushroom Caps. If you can get me those, I think I'll be able to make it.") then
		me:dialog(npc, "Breaking through the barrier will require the Poison Mushroom Cap. Talk to me when you change your mind.", false, false)
		return
	end

	q:start(npc, true)
	me:dialog(npc, "Please defeat #bPoison Mushrooms#k and bring back #b100 Poison Mushroom Caps#k in return.", false, true)
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "Have you gathered up the 100 Poison Mushroom Caps like I asked you to get?", false, true)

	local code = me:exchange({ item = { [4000500] = 100 } }, { exp = 13500 })
	if code == ExchangeResult.LackCapacity then
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	me:dialog(npc, "I am amazed that you were able to gather up these 100 Poison Mushroom Caps, which is considered a difficult feat. I think I'll be able to make #bKiller Mushroom Spores#k our of these.", false, true)
	q:force_complete(npc)
end
