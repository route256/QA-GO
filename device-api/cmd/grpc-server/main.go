package main

import (
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/repo"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/retranslator"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/sender"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/config"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/database"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/pkg/logger"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/server"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/tracer"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	batchSize uint = 2
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

	log.Info().
		Str("version", cfg.Project.Version).
		Str("commitHash", cfg.Project.CommitHash).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	// default: zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Project.Debug {
		zapLogger, _ := zap.NewDevelopment()
		logger.SetLogger(zapLogger.Sugar())
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := database.NewPostgres(dsn, cfg.Database.Driver)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed init postgres")
	}
	defer db.Close()

	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			log.Error().Err(err).Msg("Migration failed")

			return
		}
	}

	tracing, err := tracer.NewTracer(&cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed init tracing")

		return
	}
	defer tracing.Close()

	go setupRetranslator(db, cfg.Kafka)

	if err := server.NewGrpcServer(db, batchSize).Start(&cfg); err != nil {
		log.Error().Err(err).Msg("Failed creating gRPC server")

		return
	}
}

func setupRetranslator(db *sqlx.DB, kafkaCfg config.Kafka) {
	sigs := make(chan os.Signal, 1)

	sndr, err := sender.NewEventSender(kafkaCfg.Brokers, kafkaCfg.Topic)
	if err != nil {
		log.Fatal().Err(err)
	}

	rtrCfg := retranslator.Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 10 * time.Second,
		ProducerCount:  28,
		WorkerCount:    2,

		Repo:   repo.NewEventRepo(db, batchSize),
		Sender: sndr,
	}

	rtr := retranslator.NewRetranslator(rtrCfg)
	rtr.Start()
	defer rtr.Close()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
