-- NPC name (String.wz/Npc.img.xml): 메이플 운영자

local LEAF = 4001126
local CHANGED_MIND = "마음이 바뀌셨나요? 천천히 잘 생각해 보시고 결정하셔도 괜찮습니다."
local FAILED = "재료가 부족한 것은 아닌지, 인벤토리 공간이 부족한 건 아닌지 확인해 주세요. 재료를 장비하고 있을 경우 장비를 해제하여 주시기 바랍니다."
local PICK_WEAPON = "교환하실 무기를 선택해 주세요."

local WEAPONS = {
	{ material = 1302020, leaf = 2000, results = { 1302064, 1402039 } },
	{ material = 1382009, leaf = 2000, results = { 1372034, 1382039 } },
	{ material = 1452016, leaf = 2000, results = { 1452045 } },
	{ material = 1462014, leaf = 2000, results = { 1462040 } },
	{ material = 1472030, leaf = 2000, results = { 1472055 } },
	{ material = 1302030, leaf = 1500, results = { 1302064, 1402039 } },
	{ material = 1332025, leaf = 1500, results = { 1332055, 1332056 } },
	{ material = 1382012, leaf = 1500, results = { 1372034, 1382039 } },
	{ material = 1412011, leaf = 1500, results = { 1412027, 1312032 } },
	{ material = 1422014, leaf = 1500, results = { 1422029, 1322054 } },
	{ material = 1432012, leaf = 1500, results = { 1432040 } },
	{ material = 1442024, leaf = 1500, results = { 1442051 } },
	{ material = 1452022, leaf = 1500, results = { 1452045 } },
	{ material = 1462019, leaf = 1500, results = { 1462040 } },
	{ material = 1472032, leaf = 1500, results = { 1472055 } },
}

local HATS = {
	{ option = "#v1002508:# #t1002508# 얻기", leaf = 100, result = 1002508, confirm = "#t1002508#을 얻기 위해서는 다음과 같은 재료가 필요합니다. 아참! 주의사항이 있어요. 모자의 능력치는 랜덤하게 결정된답니다. 교환하시겠어요?\r\n\r\n#v4001126# 100 개" },
	{ option = "#v1002509:# 1차 업그레이드 하기", material = 1002508, leaf = 200, result = 1002509, confirm = "#t1002509#을 얻기 위해서는 다음과 같은 재료가 필요합니다. 아참! 주의사항이 있어요. 모자의 능력치는 랜덤하게 결정된답니다. 교환하시겠어요?\r\n\r\n#v4001126# 200 개 + #v1002508#" },
	{ option = "#v1002510:# 2차 업그레이드 하기", material = 1002509, leaf = 300, result = 1002510, confirm = "#t1002510#을 얻기 위해서는 다음과 같은 재료가 필요합니다. 아참! 주의사항이 있어요. 모자의 능력치는 랜덤하게 결정된답니다. 교환하시겠어요?\r\n\r\n#v4001126# 300 개 + #v1002509#" },
	{ option = "#v1002511:# 3차 업그레이드 하기", material = 1002510, leaf = 400, result = 1002511, confirm = "#t1002511#을 얻기 위해서는 다음과 같은 재료가 필요합니다. 아참! 주의사항이 있어요. 모자의 능력치는 랜덤하게 결정된답니다. 교환하시겠어요?\r\n\r\n#v4001126# 400 개 + #v1002510#" },
}

local SHIELDS = { 1092045, 1092046, 1092047 }

local SCROLLS = { 2040315, 2040912, 2043013, 2043108, 2043208, 2043308, 2043708, 2043808, 2044008, 2044108, 2044208, 2044308, 2044408, 2044508, 2044608, 2044708 }

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

local function trade(me, npc, material, leaf, result)
	local cost = { item = { [LEAF] = leaf } }
	if material ~= nil then
		cost.item[material] = 1
	end
	local code = me:exchange(cost, { item = { [result] = 1 }, random = true })
	if code ~= ExchangeResult.OK then
		me:dialog(npc, FAILED)
		return
	end
	me:dialog(npc, "즐거운 하루 되세요~")
end

local function has_maple_hat(me)
	local cap = me:equipped(EquipmentPart.Cap)
	for _, hat in ipairs(HATS) do
		if item_count(me, hat.result) > 0 then
			return true
		end
		if cap ~= nil and cap:wz():id() == hat.result then
			return true
		end
	end
	return false
end

return {
	on_click = function(me, npc)
		local menu = me:dialog_list(npc, "안녕하세요~ 믿어지세요? #e#r메이플스토리#k#n가 벌써 #e#r4주년#k#n이 되었어요. 4년 동안 변함없는 사랑을 보내 주신 모험가 여러분을 위해서 특별히 제작한 4주년 무기를 나누어 드리고 있습니다.", {
			"4주년 무기를 얻고 싶어요",
			"메이플 모자를 얻고 싶어요.",
			"4주년 메이플 실드를 얻고 싶어요.",
			"메이플 이어링을 얻고 싶어요.",
			"4주년 전용 주문서를 얻고 싶어요.#k",
		})
		if menu == nil then
			return
		end

		if menu == 1 then
			local options = {}
			for i, w in ipairs(WEAPONS) do
				options[i] = "#v4001126# #b" .. w.leaf .. " 개#k + #v" .. w.material .. "#"
			end
			local sel = me:dialog_list(npc, "4주년 무기를 얻기 위해서는 #b메이플 무기#k와 #b단풍잎#k이 필요하답니다. 어떤 무기를 가져오셨나요?", options)
			if sel == nil then
				return
			end
			local w = WEAPONS[sel]
			local results = {}
			for i, result in ipairs(w.results) do
				results[i] = "#v4001126# #b" .. w.leaf .. " 개#k + #v" .. w.material .. "# = #v" .. result .. ":#"
			end
			local pick = me:dialog_list(npc, PICK_WEAPON, results)
			if pick == nil then
				return
			end
			if not me:dialog_yes_no(npc, "아참! 주의사항이 있어요. 무기의 능력치는 랜덤하게 결정된답니다. 그리고 만약 교환하실 무기를 2개 이상 가지고 계시다면 인벤토리 가장 앞쪽에 있는 무기가 교환된다는 것도 알아 두시길 바래요. 정말 교환하시겠어요?") then
				me:dialog(npc, CHANGED_MIND)
				return
			end
			trade(me, npc, w.material, w.leaf, w.results[pick])
		elseif menu == 2 then
			local options = {}
			for i, hat in ipairs(HATS) do
				options[i] = hat.option
			end
			local sel = me:dialog_list(npc, "어떤 모자를 원하시나요? 원하는 메이플 모자를 선택하세요.", options)
			if sel == nil then
				return
			end
			local hat = HATS[sel]
			if not me:dialog_yes_no(npc, hat.confirm) then
				me:dialog(npc, CHANGED_MIND)
				return
			end
			if hat.material == nil and has_maple_hat(me) then
				me:dialog(npc, "이미 메이플 모자를 가지고 계시군요? 메이플 모자는 2개 이상 소유할 수 없답니다.")
				return
			end
			trade(me, npc, hat.material, hat.leaf, hat.result)
		elseif menu == 3 then
			local options = {}
			for i, shield in ipairs(SHIELDS) do
				options[i] = "#v4001126# #b1000 개#k + #v1092030# = #v" .. shield .. ":#"
			end
			local sel = me:dialog_list(npc, "메이플 실드와 단풍잎만 있다면 #b4주년 메이플 실드#k로 교환하실 수 있답니다. 어떤 방패로 교환하시겠어요?", options)
			if sel == nil then
				return
			end
			if not me:dialog_yes_no(npc, "아참! 주의사항이 있어요. 방패의 능력치는 랜덤하게 결정된답니다.  그리고 만약 교환하실 방패를 2개 이상 가지고 계시다면 인벤토리 가장 앞쪽에 있는 방패가 교환된다는 것도 알아 두시길 바래요.정말 교환하시겠어요?") then
				me:dialog(npc, CHANGED_MIND)
				return
			end
			trade(me, npc, 1092030, 1000, SHIELDS[sel])
		elseif menu == 4 then
			me:dialog(npc, "메이플 이어링이라... 그건 리스항구의 #b쿤#k에게 가보는 것이 좋겠군요.")
		elseif menu == 5 then
			local options = {}
			for i, scroll in ipairs(SCROLLS) do
				options[i] = "#t" .. scroll .. "#"
			end
			local sel = me:dialog_list(npc, "#b#t4001126# 1000 장#k만 있다면 4주년 무기에 사용할 수 있는 전용 주문서를 얻을 수 있답니다. 목록에서 원하는 주문서를 선택해보세요.", options)
			if sel == nil then
				return
			end
			if not me:dialog_yes_no(npc, "#t4001126# 1000 장을 #b#t" .. SCROLLS[sel] .. "##k로 바꾸시겠어요?") then
				me:dialog(npc, CHANGED_MIND)
				return
			end
			trade(me, npc, nil, 1000, SCROLLS[sel])
		end
	end
}
