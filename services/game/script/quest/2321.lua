local quest_id = 2321

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_accept(npc, "Now you'll be able to penetrate the spiny vine barrier of Mushroom Forest, but before that, #bMinister of Home Affairs#k wants to have a word with you. Please go see him immediately.") then
		me:dialog(npc, "You don't seem to follow instructions well. Come see me when you are ready.", false, false)
		return
	end

	q:start(npc, true)
	me:dialog(npc, "Good luck.", false, true)
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "I have been keeping up on your fabulour work. I am aware that you have successfully created the #bKiller Mushroom Spores#k, which penetrates through the unpenetrable barrier of the forest. Congratulations!", false, true)

	local code = me:exchange({}, { exp = 2500 })
	if code ~= ExchangeResult.OK then
		return
	end
	me:dialog(npc, "The problem now is to figure out how to enter the castle.", false, true)
	q:force_complete(npc)
end
