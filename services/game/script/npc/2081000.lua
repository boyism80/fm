-- NPC name (String.wz/Npc.img.xml): 촌장 타타모

local DONATIONS = {
	{ item = 4000226, like = 2 },
	{ item = 4000229, like = 4 },
	{ item = 4000236, like = 3 },
	{ item = 4000237, like = 6 },
	{ item = 4000260, like = 3 },
	{ item = 4000261, like = 6 },
	{ item = 4000231, like = 7 },
	{ item = 4000238, like = 9 },
	{ item = 4000239, like = 12 },
	{ item = 4000241, like = 15 },
	{ item = 4000242, like = 20 },
	{ item = 4000234, like = 20 },
	{ item = 4000232, like = 20 },
	{ item = 4000233, like = 20 },
	{ item = 4000235, like = 100 },
	{ item = 4000243, like = 100 },
}

local function buy_seed(me, npc, feellike)
	local prompt = ""
	local price = 8000
	if feellike < 5000 then
		prompt = "자네는 아직 우리 마을에 온지 얼마 되지 않았나 보군. 무엇을 도와 주면 좋겠는가?"
		price = 30000
	elseif feellike < 24000 then
		prompt = "우리 만난적 있는가..? 낯이 많이 익군. 허허.. 무엇을 도와 주면 좋겠는가?"
		price = 27000
	elseif feellike < 50000 then
		prompt = "날씨가 참 좋구먼~ 이렇게 좋은 날씨에는 가족과 소풍을 다녀 오는건 어떻겠는가? 흠. 내가 자네를 처음 만났을때가 생각나는군. 그땐 몰랐지만, 이렇게 우리 마을을 위해 힘써주고.. 고마운 젊은이가 될 줄 누가 알았겠는가. 허허..\r\n무엇을 도와 주면 좋겠는가?"
		price = 24000
	elseif feellike < 200000 then
		prompt = "무엇을 도와주면 좋겠는가?"
		price = 18000
	elseif feellike < 800000 then
		prompt = "무엇을 도와 주면 좋겠는가?"
		price = 12000
	end
	if me:dialog_list(npc, prompt, { "마법의 씨앗을 구매하고 싶습니다." }) == nil then
		return
	end

	local input = me:dialog_input(npc, "흐음. 마법의 씨앗을 구매하고 싶단건가? 귀한 물건이기 때문에 많이 줄 수 없다네. 가격은 개당 #b" .. price .. " 메소#k 일세. 몇개 구매하고 싶은가?")
	if input == nil then
		return
	end
	local qty = tonumber(input)
	if qty == nil or qty < 1 or qty > 100 then
		me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
		return
	end
	if not me:dialog_yes_no(npc, "흐음. 그렇다면 #b#t4031346##k을 " .. qty .. "개 구매하고 싶은 것인가? 총 가격은 " .. price * qty .. " 메소라네. 정말 구매하겠는가?") then
		me:dialog(npc, "알았네. 구매할 마음이 생기면 다시 찾아오게나.")
		return
	end

	local code = me:exchange({ meso = price * qty }, { item = { [4031346] = qty } })
	if code ~= ExchangeResult.OK then
		me:dialog(npc, "메소는 충분히 갖고 있는지, 또는 인벤토리 공간이 부족한건 아닌지 다시 한번 확인해 해보게. 총 가격은 #b" .. price * qty .. "#k 메소라네.")
		return
	end
	me:dialog(npc, "좋은곳에 사용하기를 바라네.")
end

local function donate(me, npc, feellike)
	local options = {}
	for i, donation in ipairs(DONATIONS) do
		options[i] = "#t" .. donation.item .. "#"
	end
	local sel = me:dialog_list(npc, "마을에 기부해 줄 물건이라도 있는가?", options)
	if sel == nil then
		return
	end
	local donation = DONATIONS[sel]

	local input = me:dialog_input(npc, "#b#t" .. donation.item .. "##k 아이템을 마을에 기부하고 싶다 이건가..? 몇개나 기부하고 싶은가?")
	if input == nil then
		return
	end
	local qty = tonumber(input)
	if qty == nil or qty < 1 or qty > 100 then
		me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
		return
	end
	if not me:dialog_yes_no(npc, "오오. 그렇다면 #b#t" .. donation.item .. "##k 아이템을 " .. qty .. " 개 기부하고 싶다 이건가?") then
		me:dialog(npc, "흐음. 생각 해 보고 결정해 주게나.")
		return
	end

	local code = me:exchange({ item = { [donation.item] = qty } }, nil)
	if code ~= ExchangeResult.OK then
		feellike = feellike - 10
		if feellike < 100 then
			feellike = 100
		end
		me:records():set("leafre.donation", feellike)
		me:dialog(npc, "자네.. 기부해 줄 물건은 제대로 갖고 있는건가?")
		return
	end
	feellike = feellike + donation.like * qty
	if feellike > 800000 then
		feellike = 800000
	end
	me:records():set("leafre.donation", feellike)
	me:dialog(npc, "마을을 위해 물건을 기부해 주어서 정말 고맙네!")
end

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "내가 도울 일이라도 있는가?\r\n#b", { "마법의 씨앗을 구매합니다.", "리프레를 위해서 무언가 하고 싶습니다." })
		if sel == nil then
			return
		end

		local feellike = me:records():get("leafre.donation")

		if sel == 1 then
			buy_seed(me, npc, feellike)
		else
			donate(me, npc, feellike)
		end
	end
}
