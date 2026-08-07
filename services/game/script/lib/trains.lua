local ferry = require("script/lib/ferry")

local M = {}

function M.board(me, npc, departed_text, ready_text)
	ferry.board(me, npc, {
		group = "Trains",
		ticket_low = 4031073,
		ticket_high = 4031074,
		departed_text = departed_text,
		ready_text = ready_text,
	})
end

function M.sell_train_ticket(me, npc)
	ferry.sell(me, npc, {
		ticket_low = 4031073,
		ticket_high = 4031074,
		cost_low = 1000,
		cost_high = 2000,
	})
end

return M
