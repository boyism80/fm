local ex = require("script/lib/expedition")

local M = {}

local EXIT_MAP = 105100100
local GROUND_Y = 258
local AWAKE_X = 416
local SEAL_BREAK_DAMAGE = 120000
local CLEAR_TIMER_MS = 605000
local CLEAR_WARP_MS = 5000

local MODE = {
	balrog_normal = {
		tomb = 105100400,
		clear = 105100401,
		body = 8830007,
		left = 8830008,
		right = 8830009,
		seal = 8830011,
		broken_seal = 8830013,
		sealed_x = 412,
		sealed_hp = 2000000000,
		spawn_ms = 5000,
		check_ms = 305000,
		battle_ms = 1200000,
	},
	balrog_hard = {
		tomb = 105100300,
		clear = 105100301,
		body = 8830000,
		left = 8830001,
		right = 8830002,
		seal = 8830004,
		broken_seal = 8830006,
		sealed_x = 416,
		sealed_hp = 4294967295,
		spawn_ms = 0,
		check_ms = 600000,
		battle_ms = 1800000,
	},
}

local function remove(map, id)
	for oid in pairs(map:mobs(id)) do
		map:remove_mob(oid)
	end
end

function M.machine(key)
	local mode = MODE[key]
	return {
		on_init = function(group)
			group:min_players(1)
			group:max_machines(nil)
			group:exit_map(EXIT_MAP)
		end,

		on_create = function(sm)
			return { mode.tomb, mode.clear }
		end,

		on_start = function(sm)
			sm:after("spawn", mode.spawn_ms, "on_spawn")
			sm:after("check", mode.check_ms, "on_check")
		end,

		on_player_enter = function(sm, player)
			local function on_arrive(player)
				ex.stamp(player, key)
			end
			player:map(sm:map(mode.tomb), { callback = on_arrive })
		end,

		on_spawn = function(sm)
			local tomb = sm:map(mode.tomb)
			for _, id in ipairs({ mode.body, mode.seal, mode.right }) do
				local mob = tomb:spawn_mob(id, mode.sealed_x, GROUND_Y)
				if mob ~= nil then
					mob:max_hp(mode.sealed_hp)
					mob:hp(mode.sealed_hp)
				end
			end
			sm:start_timer(mode.battle_ms)
		end,

		on_check = function(sm)
			local tomb = sm:map(mode.tomb)
			local dealt = 0
			for _, mob in pairs(tomb:mobs()) do
				dealt = dealt + mob:max_hp() - mob:hp()
			end
			if dealt <= SEAL_BREAK_DAMAGE then
				sm:message("당신이 극복하기엔 아직 발록이 너무나 강합니다.")
				sm:finish(EXIT_MAP)
				return
			end
			tomb:spawn_mob(mode.broken_seal, AWAKE_X, GROUND_Y)
			remove(tomb, mode.seal)
			remove(tomb, mode.broken_seal)
			remove(tomb, mode.body)
			remove(tomb, mode.right)
			tomb:spawn_mob(mode.body, AWAKE_X, GROUND_Y)
			tomb:spawn_mob(mode.left, AWAKE_X, GROUND_Y)
			tomb:spawn_mob(mode.right, AWAKE_X, GROUND_Y)
			sm:set_property("awake", "1")
		end,

		on_mob_die = function(sm, mob)
			if sm:get_property("awake") ~= "1" then
				return
			end
			local tomb = sm:map(mode.tomb)
			for _, id in ipairs({ mode.body, mode.left, mode.right }) do
				if tomb:mob_by_template(id) ~= nil then
					return
				end
			end
			sm:set_property("awake", "0")
			sm:message("마왕 발록을 물리쳤습니다!")
			tomb:show_effect("balog/clear/stone")
			sm:restart_timer(CLEAR_TIMER_MS)
			sm:after("clear", CLEAR_WARP_MS, "on_clear")
		end,

		on_clear = function(sm)
			sm:warp_all(mode.tomb, mode.clear)
		end,

		on_scheduled_timeout = function(sm)
			sm:finish(EXIT_MAP)
		end,
	}
end

return M
