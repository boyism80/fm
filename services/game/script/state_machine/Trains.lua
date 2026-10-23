-- State machine (old/scripts/event/Trains.js): 오르비스↔루디브리엄 열차

return require("script/lib/ferry_route").new({
    name = "Trains",
    legs = {
        { waiting = 200000122, ride = 200090100, station = 200000121, dest = 220000110 },
        { waiting = 220000111, ride = 200090110, station = 220000110, dest = 200000121 },
    },
})
