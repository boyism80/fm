-- NPC name (String.wz/Npc.img.xml): 사서 위즈

local BOOKS = {
	{ quest = 3615, item = 4031235 },
	{ quest = 3616, item = 4031236 },
	{ quest = 3617, item = 4031237 },
	{ quest = 3618, item = 4031238 },
	{ quest = 3630, item = 4031270 },
	{ quest = 3633, item = 4031280 },
	{ quest = 3639, item = 4031298 },
}

return {
	on_click = function(me, npc)
		local counter = 0
		local books = ""
		for _, book in ipairs(BOOKS) do
			local q = me:quest(book.quest)
			if q ~= nil and q:completed() then
				counter = counter + 1
				books = books .. "\r\n#v" .. book.item .. "# #b#t" .. book.item .. "##k"
			end
		end
		if counter == 0 then
			me:dialog(npc, "#b#h ##k님은 아직 회수하신 책이 없군요.")
			return
		end
		if not me:dialog(npc, "어디보자... #b#h ##k님은 총 #b" .. counter .. "#k권의 책을 회수하셨군요. 다음과 같은 책을 회수하셨습니다.:\r\n" .. books, false, true) then
			return
		end
		me:dialog(npc, "저희 도서관은 모험가님께 진심으로 감사를 드리고 있습니다. #b#h ##k님의 도움으로, 잃어버린 책들과 이야기들을 다시 찾았답니다. 앞으로도 많은 도움을 부탁드립니다.", true, true)
	end
}
