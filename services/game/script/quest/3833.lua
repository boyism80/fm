-- Quest name (Quest.wz/Quest.img.xml): 모자라도 약재 구하기

local quest_id = 3833

local function item_count(me, item_id)
	local slots = me:item(item_id)
	local count = 0
	for _, it in pairs(slots) do
		count = count + it:count()
	end
	return count
end

local function tier_for(n)
	if n >= 1000 then
		return {
			dialog = "허... 허허허허. 자네는 인간인가?! 이렇게 많이 구해 오다니, 도대체 늙은 도라지를 몇 마리나 학살한 겐가...?! 험. 아무튼 고맙네. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 54000,
			reward = { [2000005] = 50, [2040501] = 1 },
			ok = "이 정도면 태상이 기절하지 않을까...",
		}
	elseif n >= 900 then
		return {
			dialog = "아니... 정말 이게 모두 자네가 구한 것인가? 괴, 굉장하군... 이 정도까지 구하다니... 자네는 진정 굉장한 모험가일세. 정말 훌륭해! 이 정도면 이런 걸 줘도 아깝지 않지. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 54000,
			reward = { [2020013] = 50, [2040502] = 1 },
			ok = "이렇게 많은 걸 무슨 수로 무릉까지 가져간다...?",
		}
	elseif n >= 700 then
		return {
			dialog = "호오... 정말 굉장하군. 여기까지 모으는 게 쉽지 않았을 텐데. 정말 고맙네. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 54000,
			reward = { [2000004] = 50, [2040500] = 1 },
			ok = "흠. 잘못하다간 도라지를 배에 모두 싣기도 어렵겠는걸.",
		}
	elseif n >= 500 then
		return {
			dialog = "오오... 많군. 이 정도면 태상도 할 말 없겠지. 자네 보기보다 제법이군. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 54000,
			reward = { [2020013] = 50 },
			ok = "도라지를 배까지 옮기는 것도 쉽지 않겠는걸. 황선장에게 옮겨달라고 부탁해야겠군...",
		}
	elseif n >= 300 then
		return {
			dialog = "흠... 뭐 이 정도면 당분간 태상도 별 말 않겠지. 고맙네. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 51000,
			reward = { [2020012] = 50 },
			ok = "이 정도면 태상이 만족하겠지.",
		}
	elseif n >= 200 then
		return {
			dialog = "좋아. 이 정도면 그럭저럭 태상을 달래놓을 수 있을 것 같군. 감사하지. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 48000,
			reward = { [2001001] = 50 },
			ok = "자아, 그럼 도라지를 잘 말려봐야겠군. 약재로 쓰려면 말리는 것도 중요하지.",
		}
	elseif n >= 100 then
		return {
			dialog = "흐음. 좀 모자라긴 하지만 뭐, 당장 급한 건 이쪽이니 어쩔 수 없지. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 45000,
			reward = { [2020008] = 50 },
			ok = "쯧. 태상이 한동안 투덜대겠군..",
		}
	elseif n >= 50 then
		return {
			dialog = "에잉. 겨우 이 정도인가? 아주 적은 건 아니지만 만족스럽진 않군. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 10000,
			reward = { [2020007] = 50 },
			ok = "약탈만 아니었더라도... 이런 식으로 구할 필요 없을 텐데.",
		}
	elseif n >= 10 then
		return {
			dialog = "겨우 이게 단가? 허음... 자네 능력이 이것뿐이 안 된다니 뭐라 할 수는 없겠지만... 쯧쯧쯧. 젊은이가 이래서야 원. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 1000,
			reward = { [2022144] = 10 },
			ok = "이래가지고는 태상을 볼 면목이 없군...",
		}
	elseif n >= 1 then
		return {
			dialog = "...뭐, 그래. 어쨌든 가져다 주긴 하는구만. 자네가 구한 #b#t4000294# %d개#k를 어서 나에게 주게.\r\n\r\n",
			exp = 10,
			reward = { [2000000] = 1 },
			ok = "한동안 태상과 연락을 하지 않는 게 좋을 듯 하군...",
		}
	end
	return nil
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not q:started() then
		q:start(npc, true)
		return
	end

	local n_item = item_count(me, 4000294)
	local tier = tier_for(n_item)
	if tier == nil then
		me:dialog(npc, "아직 100년 묵은 도라지는 구하지 못한 모양이군. 그건 늙은 도라지를 잡아 얻을 수 있다네.", false, false)
		return
	end

	local file = "#fUI/UIWindow.img/QuestIcon/"
	local str = string.format(tier.dialog, n_item)
	str = str .. file .. "4/0#\r\n" .. file .. "5/0#\r\n\r\n" .. file .. "8/0# " .. tier.exp .. " exp"
	if not me:dialog(npc, str, false, true) then
		return
	end

	local code = me:exchange(
		{ item = { [4000294] = n_item } },
		{ item = tier.reward, exp = tier.exp }
	)
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "뭘 그렇게 많이 들고 다니는건가? 인벤토리에 빈 칸이 있는지 확인해 주게.", false, false)
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	q:force_complete(npc)
	me:dialog(npc, tier.ok, false, false)
end
