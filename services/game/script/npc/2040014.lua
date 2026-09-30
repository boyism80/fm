-- NPC name (String.wz/Npc.img.xml): 치코

local BOARDS = {
	[100000203] = {
		prompt = "후후.. 난 미니게임의 마스터 #b카이지#k라고 하지. 어때..? 게임판을 만들어 보고 싶은거야?\r\n\r\n#b",
		list = {
			{ 4080000, { { 4030000, 1 }, { 4030001, 1 }, { 4030009, 1 } } },
			{ 4080001, { { 4030000, 1 }, { 4030010, 1 }, { 4030009, 1 } } },
			{ 4080002, { { 4030000, 1 }, { 4030011, 1 }, { 4030009, 1 } } },
			{ 4080003, { { 4030010, 1 }, { 4030001, 1 }, { 4030009, 1 } } },
			{ 4080004, { { 4030011, 1 }, { 4030010, 1 }, { 4030009, 1 } } },
			{ 4080005, { { 4030011, 1 }, { 4030001, 1 }, { 4030009, 1 } } },
			{ 4080100, { { 4030012, 15 } } },
		},
	},
	[220000300] = {
		prompt = "게임판을 만들어 보고 싶은거야?\r\n\r\n#b",
		list = {
			{ 4080006, { { 4030013, 1 }, { 4030014, 1 }, { 4030009, 1 } } },
			{ 4080007, { { 4030013, 1 }, { 4030016, 1 }, { 4030009, 1 } } },
			{ 4080008, { { 4030014, 1 }, { 4030016, 1 }, { 4030009, 1 } } },
			{ 4080009, { { 4030015, 1 }, { 4030013, 1 }, { 4030009, 1 } } },
			{ 4080010, { { 4030015, 1 }, { 4030014, 1 }, { 4030009, 1 } } },
			{ 4080011, { { 4030015, 1 }, { 4030016, 1 }, { 4030009, 1 } } },
			{ 4080100, { { 4030012, 15 } } },
		},
	},
}

return {
	on_click = function(me, npc)
		local field = me:map()
		local wz = field ~= nil and field:wz() or nil
		if wz == nil then
			return
		end
		local boards = BOARDS[wz:id()]
		if boards == nil then
			return
		end

		local options = {}
		for i, entry in ipairs(boards.list) do
			options[i] = " #t" .. entry[1] .. "#"
		end
		local sel = me:dialog_list(npc, boards.prompt, options)
		if sel == nil then
			return
		end

		local board = boards.list[sel]
		local text = "#b#i" .. board[1] .. "# #t" .. board[1] .. "##k 게임을 만들어 보고 싶은거야? 흐음... 재료는 \r\n"
		for _, mat in ipairs(board[2]) do
			text = text .. "#b#i" .. mat[1] .. "# #t" .. mat[1] .. "##k " .. mat[2] .. " 개\r\n"
		end
		text = text .. "\r\n\r\n인데, 정말 만들어 보고 싶어? 특별히 제작비는 받지 않도록 하지."
		if not me:dialog_yes_no(npc, text) then
			return
		end

		local cost_items = {}
		for _, mat in ipairs(board[2]) do
			cost_items[mat[1]] = mat[2]
		end
		local code = me:exchange({ item = cost_items }, { item = { [board[1]] = 1 } })
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "흐음.. 인벤토리 공간이 부족한 것 같은데? 다시 한번 확인해 줄래?")
			return
		end
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "부족한 재료가 있는건 아닌지 다시 한번 확인해 줄래?")
			return
		end
		me:dialog(npc, "좋아.. 완성되었어. 재밌게 게임을 즐기라구.")
	end
}
