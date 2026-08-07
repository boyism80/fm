local EXIT_MAP = 280090000
local FIRE_ORE_PIECE = 4031061
local PAPER_REWARD = 2030007
local STAGE1_QUEST = 100001

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			me:play_portal_sound()
			me:map(EXIT_MAP)
			return
		end
		if sm:get_property("clear") ~= "1" then
			me:notice("아직 임무를 완료하지 못하여 다음 맵으로 이동할 수 없습니다.", Msg.PinkText)
			return
		end
		local paper = sm:get_property("paper") == "1"
		local reward = {
			item = { [FIRE_ORE_PIECE] = 1 },
			exp = 12000,
		}
		if paper then
			reward = {
				item = {
					[FIRE_ORE_PIECE] = 1,
					[PAPER_REWARD] = 5,
				},
				exp = 20000,
			}
		end
		local code = me:exchange({}, reward)
		if code == ExchangeResult.LackCapacity or code ~= ExchangeResult.OK then
			me:notice("인벤토리에 공간이 부족하여 다음 맵으로 이동할 수 없습니다.", Msg.PinkText)
			return
		end
		local q = me:quest(STAGE1_QUEST)
		if q ~= nil and q:started() then
			q:force_complete(0)
		end
		me:play_portal_sound()
		me:map(EXIT_MAP, 1)
	end
}
