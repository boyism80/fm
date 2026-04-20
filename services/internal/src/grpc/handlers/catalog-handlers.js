"use strict";

function createCatalogHandlers(internalConfig, messages, grpcError) {
    return {
        ping(_call, callback) {
            try {
                const reply = new messages.PingReply();
                reply.setMessage("pong");
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        getServerCatalog(_call, callback) {
            try {
                const reply = new messages.GetServerCatalogReply();
                const worlds = internalConfig.game_servers?.worlds ?? {};
                const worldMsgs = [];
                for (const [worldId, world] of Object.entries(worlds)) {
                    const wm = new messages.WorldCatalog();
                    wm.setWorldId(Number(worldId));
                    wm.setWorldName(world.world_name || `World-${worldId}`);
                    wm.setFlag(Number(world.flag ?? 0));
                    wm.setEventMessage(world.event_message || "");
                    const channels = Array.isArray(world.channels) ? world.channels : [];
                    wm.setChannelsList(channels.map((ch) => {
                        const cm = new messages.ChannelCatalog();
                        cm.setChannelId(Number(ch.channel_id));
                        cm.setHost(ch.host);
                        cm.setPort(Number(ch.port));
                        cm.setName(ch.name || `Channel ${Number(ch.channel_id) + 1}`);
                        return cm;
                    }));
                    worldMsgs.push(wm);
                }
                reply.setWorldsList(worldMsgs);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

module.exports = { createCatalogHandlers };
