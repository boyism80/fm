local quest_id = 2320

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_accept(npc, "I have just one more request for you. Would you like to take a listen?") then
		me:dialog(npc, "I wanted you to personally give this piece of good news to #bBruce#k, but I understand if you're busy.", false, false)
		return
	end

	q:start(npc, true)
	local code = me:exchange({}, { item = { [4032389] = 1 } })
	if code == ExchangeResult.LackCapacity then
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	me:dialog(npc, "To be honest, these #bKiller Mushroom Spores#k are not completely out of my own work. Do you remember #bBruce#k from #bHenesys#k? I have been friends with him since childhood, and #bKiller Mushroom Spores#k was completed after he shared the results of his studies with me. This was all thanks to him, so I'd like for you to give this to him for me.", false, true)
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	me:dialog(npc, "Oh! You're here on behalf of #bScarrs#k? \r\n\r\n#fUI/UIWindow.img/QuestIcon/4/0# \r\n#fUI/UIWindow.img/QuestIcon/8/0# 8800 exp", false, true)

	local code = me:exchange({ item = { [4032389] = 1 } }, { exp = 8800 })
	if code == ExchangeResult.LackCapacity then
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	me:dialog(npc, "Ahh, so this is the #bKiller Mushroom Spores#k that I was working on in the past. I had a tough time gathering up the ingredients, so I left it in theory only, but he was able to complete it, with a sample to show for as well. Please tell him I appreciate his good work.", false, true)
	q:force_complete(npc)
end
