-- NPC name (String.wz/Npc.img.xml): 카이지

local henesys = {
	{ 4080000, { { 4030000, 1 }, { 4030001, 1 }, { 4030009, 1 } } },
	{ 4080001, { { 4030000, 1 }, { 4030010, 1 }, { 4030009, 1 } } },
	{ 4080002, { { 4030000, 1 }, { 4030011, 1 }, { 4030009, 1 } } },
	{ 4080003, { { 4030010, 1 }, { 4030001, 1 }, { 4030009, 1 } } },
	{ 4080004, { { 4030011, 1 }, { 4030010, 1 }, { 4030009, 1 } } },
	{ 4080005, { { 4030011, 1 }, { 4030001, 1 }, { 4030009, 1 } } },
	{ 4080100, { { 4030012, 15 } } },
}

local ludibrium = {
	{ 4080006, { { 4030013, 1 }, { 4030014, 1 }, { 4030009, 1 } } },
	{ 4080007, { { 4030013, 1 }, { 4030016, 1 }, { 4030009, 1 } } },
	{ 4080008, { { 4030014, 1 }, { 4030016, 1 }, { 4030009, 1 } } },
	{ 4080009, { { 4030015, 1 }, { 4030013, 1 }, { 4030009, 1 } } },
	{ 4080010, { { 4030015, 1 }, { 4030014, 1 }, { 4030009, 1 } } },
	{ 4080011, { { 4030015, 1 }, { 4030016, 1 }, { 4030009, 1 } } },
	{ 4080100, { { 4030012, 15 } } },
}

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil then
			return
		end
		local boards = nil
		local prompt = ""
		if map:template_id() == 100000203 then
			boards = henesys
			prompt = "후후.. 난 미니게임의 마스터 #b카이지#k라고 하지. 어때..? 게임판을 만들어 보고 싶은거야?"
		elseif map:template_id() == 220000300 then
			boards = ludibrium
			prompt = "게임판을 만들어 보고 싶은거야?"
		else
			return
		end
		local choices = {}
		for i, board in ipairs(boards) do
			choices[i] = "#t" .. board[1] .. "#"
		end
		if #choices == 0 then
			return
		end
		local sel = me:dialog_list(npc, prompt, choices)
		if sel == nil then
			return
		end
		local board = boards[sel]
		if board == nil then
			return
		end
		local ask = "#b#i" .. board[1] .. "# #t" .. board[1] .. "##k 게임을 만들어 보고 싶은거야? 흐음... 재료는 \r\n"
		local cost = { item = {} }
		for _, mat in ipairs(board[2]) do
			ask = ask .. "#b#i" .. mat[1] .. "# #t" .. mat[1] .. "##k " .. mat[2] .. " 개\r\n"
			cost.item[mat[1]] = mat[2]
		end
		ask = ask .. "\r\n\r\n인데, 정말 만들어 보고 싶어? 특별히 제작비는 받지 않도록 하지."
		if not me:dialog_yes_no(npc, ask) then
			return
		end
		local code = me:exchange(cost, { item = { [board[1]] = 1 } })
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
