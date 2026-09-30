-- NPC name (String.wz/Npc.img.xml): 전사 미니무르

local class_names = require("script/lib/class")

local AWAKEN = 4111006
local LEAF = 4001126
local KEY_SLOT = 8

return {
	on_click = function(me, npc)
		local skill = me:skill(AWAKEN)
		if skill ~= nil and skill:level() > 0 then
			me:dialog(npc, "각성스킬이 성공적으로 조작키에 입력되었습니다.")
			me:put_key(KEY_SLOT, 1, AWAKEN)
			return
		end
		local sel = me:dialog_list(npc, "아르카나온라인에 어서오세요! 혹시 각성에 관심이 있으신가요? #b[캡틴,나로는 배우지않는것을 권장.] #r각성 가능 레벨:150#k", {
			"네",
		})
		if sel == nil then
			return
		end
		if me:level() < 120 then
			me:dialog(npc, "당신은 각성에 대해 알 자격이 없는 것 같아요. 좀더 성장 하세요.")
			return
		end
		if me:dialog(npc, "당신은 충분히 각성에 대해 알 자격이 있는거같네요. 각성에 대해 알려드릴게요.", false, true) == false then
			return
		end
		local leaves = 0
		for _, it in pairs(me:item(LEAF)) do
			leaves = leaves + it:count()
		end
		if leaves < 1 or me:meso() <= 1 then
			me:dialog(npc, "재료가 부족합니다.")
			return
		end
		local again = me:dialog_list(npc, "그럼 알려드릴게요. 각성은 자신의 가능성이에요. 보통 스스로 깨우치긴 힘든데 어느정도 키워드만 알게되면 쉽게 발현가능한  힘이거든요.우선 직업을 알아야하는데...감정해볼게요.당신은 " .. class_names.name(me:class()) .. "군요. 키워드가 입력되었으니 각성 스킬 조작키에 추가를 눌러보세요!", {
			"각성 스킬 조작키",
		})
		if again == nil then
			return
		end
		local code = me:exchange({ item = { [LEAF] = 1 }, meso = 1 }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료가 부족합니다.")
			return
		end
		if me:class() == Class.Shadower then
			local learned = me:add_skill(AWAKEN)
			if learned ~= nil then
				learned:level(20, 20)
			end
			me:put_key(KEY_SLOT, 20, AWAKEN)
		end
		me:dialog(npc, "각성스킬이 성공적으로 조작키에 입력되었습니다.")
	end
}
