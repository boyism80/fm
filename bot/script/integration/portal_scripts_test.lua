local check = require("script/integration/lib/script_check")

local BOTS = 3
local entries, missing = wz.portal_scripts()

local function enter(ctx, bot, entry)
	local where = "portal " .. entry.script .. " (map " .. entry.map .. " " .. entry.portal .. ")"
	if bot:instance_move(entry.map) == false then
		return
	end

	local warp = req.warp { target = 0xFFFFFFFF, portal_name = entry.portal }
	local p, name = bot:request(check.REPLIES, warp, check.replied, 5000)
	if p == false then
		ctx:fail(where .. ": 응답 없음")
		return
	end
	if check.is_dialog(name) then
		bot:dialog(false)
	end
end

test_suite {
	name = "Portal scripts",
	bot_count = BOTS,

	on_initialize = function(ctx)
		log("info", string.format("%d portal scripts placed in WZ, %d not placed", #entries, #missing))
		return true
	end,

	scenarios = {
		{ parallel = check.queues(BOTS, entries, enter) },
	},
}
