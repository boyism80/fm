-- NPC name (String.wz/Npc.img.xml): 알레그로

local pq = require("script/lib/party_quest")

local MIN_LEVEL = 10
local DIAMOND = 4021007
local TOKEN_FIRST = 4031683
local TOKEN_LAST = 4031692
local TOKEN_COUNT = 4
local BOX_FIRST = 2240004
local BOX_LAST = 2240015
local CRAFT_EXP = 2360

local RINGS = {
	{ box = 2240004, material = 4011007, meso = 80000 },
	{ box = 2240007, material = 4021009, meso = 40000 },
	{ box = 2240010, material = 4011006, meso = 20000 },
	{ box = 2240013, material = 4011004, meso = 10000 },
}

local INTRO = "사랑하는 사람이 있습니까? 그렇다면 영원한 사랑의 맹세인 결혼을 생각하시진 않으신지요? 만약 그러시다면 제가 몇가지 도움을 드릴 수 있을 것 같군요."

local function tokens(me)
	local owned = {}
	for id = TOKEN_FIRST, TOKEN_LAST do
		if pq.has_item(me, id) then
			table.insert(owned, id)
		end
	end
	return owned
end

local function craft(me, npc, owned)
	local names = {}
	for _, ring in ipairs(RINGS) do
		table.insert(names, "#i" .. ring.box .. ":# #t" .. ring.box .. "#")
	end
	local selected = me:dialog_list(npc, "사랑의 증표 네개를 모아오셨군요. 좋습니다. 반지의 재료에 따라 4종류의 반지를 만들 수 있답니다. 어떤 반지를 만들어 드릴까요?", names)
	if selected == nil then
		return
	end
	local ring = RINGS[selected]

	local carats = {}
	for carat = 1, 3 do
		table.insert(carats, "#e#t" .. ring.box .. "# " .. carat .. "캐럿#n\r\n#i" .. ring.material .. ":# 1개 + #i" .. DIAMOND .. ":# " .. (carat * 2) .. "개 + " .. ring.meso .. " 메소")
	end
	local carat = me:dialog_list(npc, "#t" .. ring.box .. "#에 다이아몬드를 얼마나 첨가하느냐에 따라 반지의 질이 달라진답니다. 어떤 반지를 만들고 싶은가요?", carats)
	if carat == nil then
		return
	end
	local box = ring.box + carat - 1
	if not me:dialog_yes_no(npc, "선택하신 약혼반지는 #b#i" .. box .. ":# #t" .. box .. "##k 입니다. 재료를 다시 한번 확인해 주십시오. 지금 만들고 싶으십니까?") then
		me:dialog(npc, "부족한 재료가 있으신 모양이시군요. 진정한 사랑을 위해선 철저히 준비해야 하는 법이지요.")
		return
	end

	local cost = { meso = ring.meso, item = { [ring.material] = 1, [DIAMOND] = carat * 2 } }
	for i = 1, TOKEN_COUNT do
		cost.item[owned[i]] = 1
	end
	if me:exchange(cost, { item = { [box] = 1 }, exp = CRAFT_EXP }) ~= ExchangeResult.OK then
		me:dialog(npc, "부족한 재료가 있으신건 아닌지, 메소는 충분한지, 인벤토리 공간은 충분한지 확인해 주십시오.")
		return
	end
	me:dialog(npc, "약혼반지가 완성되었습니다. 약혼하실 여성분께 프로포즈 하시면 됩니다. 사랑에 행운이 깃들기를..")
end

return {
	on_click = function(me, npc)
		if me:gender() ~= 0 then
			me:dialog(npc, INTRO .. " 하지만 제가 드릴 몇가지 안내는 남성분께서 받으실 수 있습니다. 남자분과 함께 협동해서 사랑의 증표를 모아보세요.")
			return
		end
		for id = BOX_FIRST, BOX_LAST do
			if pq.has_item(me, id) then
				me:dialog(npc, "으음, 이미 약혼 반지를 갖고계시는군요. 반지 상자를 사용해 마음에 드는 상대에게 프로포즈를 하시면 됩니다.")
				return
			end
		end
		if me:level() < MIN_LEVEL then
			me:dialog(npc, "음.. 당신은 사랑하는 사람을 지키기엔 아직 너무 약한 것 같군요. 레벨 10 이 되시면 다시 찾아오세요.")
			return
		end
		local owned = tokens(me)
		if #owned < TOKEN_COUNT then
			me:dialog(npc, INTRO .. " 반지를 원한다면 사랑의 진실함을 증명하고, 사랑의 증표 4개를 가져오세요. #b헤네시스의 마야, 엘리니아의 로웬, 페리온의 이얀, 커닝시티의 넬라, 오르비스의 에릭손, 루디브리엄의 보자관티군, 아쿠아리움의 뮤즈, 리프레의 팜, 무릉의 한태수, 아리안트의 지유르#k를 찾아가보세요.")
			return
		end
		craft(me, npc, owned)
	end
}
