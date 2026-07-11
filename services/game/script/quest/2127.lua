local quest_id = 2127

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	local qr = me:quest(7060)
	if qr == nil then
		return
	end

	local info = qr:record()
	if not qr:started() and not qr:completed() then
		if qr:wz() == nil then
			qr:start("0")
		else
			qr:start(npc, "0")
		end
		info = "0"
	elseif info == nil or info == "" then
		info = "0"
		qr:record("0")
	end

	local msgs = {
		["0"] = "떠날 준비는 다 되었나? 처음 가는 거라서 긴장을 한 것 같군. 걱정말게. 자네는 잘 할 수 있을 거야. 내 능력으로 앞으로 다섯 번은 그곳으로 보내 줄 수 있네. 지금 이동하겠나?",
		["1"] = "또 만나게 됐군. 지난 번에 갔을 때는 어땠나? 다시 가는 것을 보니 그 곳에서 하던 일이 있나보지? 자네에겐 네 번의 기회가 남아 있네. 지금 이동하겠나?",
		["2"] = "자주 보게 되는군. 무척 바쁜 모양이야? 자네에겐 세 번의 기회가 남아 있네. 지금 이동 하겠나?",
		["3"] = "사막은 매력적인 곳인가 보지? 그러고 보니 얼굴도 검게 그을린 것 같군. 자네에겐 두 번의 기회가 남아 있네. 지금 이동 하겠나?",
		["4"] = "오늘도 그곳으로 가는건가? 이름이 뭐라고 했지? 아리안트라고 했었나? 자네를 보니 나도 한 번쯤 가보고 싶어지는군. 자네에겐 한 번의 기회가 남아 있네. 지금 이동하겠나?",
	}
	local msg = msgs[info]
	if msg == nil then
		me:dialog(npc, "미안하지만 더 이상 그곳으로 보내줄 수 없네. 이미 다섯 번의 기회를 모두 사용했네.", false, false)
		if not q:completed() then
			q:force_complete(npc)
		end
		return
	end

	if not me:dialog_accept(npc, msg) then
		me:dialog(npc, "준비가 되면 다시 찾아오게.")
		return
	end
	if not me:dialog(npc, "잘 다녀오게.", false, true) then
		return
	end

	local next_info = tostring(tonumber(info) + 1)
	qr:record(next_info)
	if next_info == "5" then
		qr:force_complete(npc)
		q:force_complete(npc)
	end
	me:map(260000200)
end
