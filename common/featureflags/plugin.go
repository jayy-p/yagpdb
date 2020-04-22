package featureflags

import (
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/bot/eventsystem"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/pubsub"
)

var logger = common.GetPluginLogger(&Plugin{})

// Plugin represents the mqueue plugin
type Plugin struct {
}

// PluginInfo implements common.Plugin
func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "featureflags",
		SysName:  "featureflags",
		Category: common.PluginCategoryCore,
	}
}

// RegisterPlugin registers the mqueue plugin into the plugin system and also initializes it
func RegisterPlugin() {
	p := &Plugin{}
	common.RegisterPlugin(p)

	pubsub.AddHandler("feature_flags_updated", handleInvalidateCacheFor, nil)
}

// Invalidate the cache when the rules have changed
func handleInvalidateCacheFor(event *pubsub.Event) {
	cacheL.Lock()
	defer cacheL.Unlock()

	delete(cache, event.TargetGuildInt)
}

var _ bot.BotInitHandler = (*Plugin)(nil)

// BotInit implements bot.BotInitHandler
func (p *Plugin) BotInit() {
	eventsystem.AddHandlerAsyncLastLegacy(p, func(evt *eventsystem.EventData) {
		cacheL.Lock()
		defer cacheL.Unlock()

		delete(cache, evt.GuildDelete().ID)
	}, eventsystem.EventGuildDelete)
}