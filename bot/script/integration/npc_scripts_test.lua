local check = require("script/integration/lib/script_check")

local BOTS = 3
local entries, missing = wz.npc_scripts()

local function click(ctx, bot, entry)
	local where = "npc " .. entry.npc .. " (map " .. entry.map .. ")"
	if bot:instance_move(entry.map) == false then
		return
	end
	bot:move(entry.x, entry.y, entry.foothold)

	local oid = bot:npc(entry.npc, 3000)
	if oid == nil then
		ctx:fail(where .. ": 맵에 NPC가 없음")
		return
	end

	local p, name = bot:request(check.REPLIES, req.npc_click { oid = oid }, check.replied, 5000)
	if p == false then
		ctx:fail(where .. ": 응답 없음")
		return
	end
	if check.is_dialog(name) then
		bot:dialog(false)
	end
end

test_suite {
	name = "NPC scripts",
	bot_count = BOTS,

	on_initialize = function(ctx)
		log("info", string.format("%d npc scripts placed in WZ, %d not placed", #entries, #missing))
		return true
	end,

	scenarios = {
		{ parallel = check.queues(BOTS, entries, click) },
	},
}
