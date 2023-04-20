package service

import (
	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/config"
	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/db"
	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/event/dispatcher"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/model"
	ratelimiter "github.com/ncarlier/readflow/pkg/rate-limiter"
	"github.com/ncarlier/readflow/pkg/sanitizer"
	"github.com/ncarlier/readflow/pkg/scraper"
	"github.com/ncarlier/readflow/pkg/scripting"
	"github.com/ncarlier/readflow/pkg/secret"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	// load all exporters
	_ "github.com/ncarlier/readflow/pkg/exporter/all"
)

var instance *Registry

// Registry is the structure definition of the service registry
type Registry struct {
	conf                    config.Config
	db                      db.DB
	logger                  zerolog.Logger
	downloadCache           cache.Cache
	properties              *model.Properties
	webScraper              scraper.WebScraper
	dl                      downloader.Downloader
	hashid                  *helper.HashIDHandler
	notificationRateLimiter ratelimiter.RateLimiter
	scriptEngine            *scripting.ScriptEngine
	sanitizer               *sanitizer.Sanitizer
	events                  *event.Manager
	dispatcher              dispatcher.Dispatcher
	secretsEngineProvider   secret.EngineProvider
}

// Configure the global service registry
func Configure(conf config.Config, database db.DB, downloadCache cache.Cache) error {
	webScraper, err := scraper.NewWebScraper(constant.DefaultClient, conf.Integration.ExternalWebScraperURL)
	if err != nil {
		return err
	}
	hashid, err := helper.NewHashIDHandler(conf.Global.SecretSalt)
	if err != nil {
		return err
	}
	notificationRateLimiter, err := ratelimiter.NewRateLimiter("notification", conf.RateLimiting.Notification)
	if err != nil {
		return err
	}
	blockList, err := sanitizer.NewBlockList(conf.Global.BlockList, sanitizer.DefaultBlockList)
	if err != nil {
		return err
	}
	dispatcher, err := dispatcher.NewDispatcher(conf.Integration.ExternalEventBrokerURI)
	if err != nil {
		return err
	}
	secretsEngineProvider, err := secret.NewSecretsEngineProvider(conf.Integration.SecretsEngineProvider)
	if err != nil {
		return err
	}

	instance = &Registry{
		conf:                    conf,
		db:                      database,
		logger:                  log.With().Str("component", "service").Logger(),
		downloadCache:           downloadCache,
		webScraper:              webScraper,
		dl:                      downloader.NewDefaultDownloader(downloadCache),
		hashid:                  hashid,
		notificationRateLimiter: notificationRateLimiter,
		sanitizer:               sanitizer.NewSanitizer(blockList),
		scriptEngine:            scripting.NewScriptEngine(128),
		dispatcher:              dispatcher,
		events:                  event.NewEventManager(),
		secretsEngineProvider:   secretsEngineProvider,
	}
	instance.registerEventHandlers()
	return instance.initProperties()
}

// Lookup returns the global service registry
func Lookup() *Registry {
	if instance != nil {
		return instance
	}
	log.Fatal().Msg("service registry not configured")
	return nil
}
