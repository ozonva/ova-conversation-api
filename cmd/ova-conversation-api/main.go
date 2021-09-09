package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"

	"ova-conversation-api/internal/api"
	"ova-conversation-api/internal/kafka"
	"ova-conversation-api/internal/repo"
	conversationApi "ova-conversation-api/pkg/api/github.com/ozonva/ova-conversation-api/pkg/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Msgf("open file env: %v", err)
	}

	metricsServer := initMetrics()

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()

	jaegerRunner, err := initJaeger()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	defer jaegerRunner.Close()

	run()
}

func run() {
	grpcPort := os.Getenv("GRPC_PORT")
	l, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal().Msgf("listening to TCP: %v", err)
	}
	log.Info().Msgf("ova-conversation-api: gRPC server started on the port %s", grpcPort)

	dbConn := ConnectToDB(dsnForDB())

	kafkaProducer, err := kafka.NewProducer([]string{os.Getenv("KAFKA_ADDR")}, "conversations")
	if err != nil {
		log.Fatal().Msgf("connecting to kafka: %v", err)
	}

	service := grpc.NewServer()
	conversationApi.RegisterConversationApiServer(service, api.NewConversationApiServer(repo.NewRepo(dbConn), kafkaProducer))
	if err = service.Serve(l); err != nil {
		log.Fatal().Msgf("failed to serve: %s", err)
	}
}

func ConnectToDB(dsn string) *sqlx.DB {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatal().Msgf("connect do db error %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Msgf("failed to connect to db: %v", err)
	}

	return db
}

func dsnForDB() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_ADDR"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

func initMetrics() *http.Server {
	mux := http.DefaultServeMux
	mux.Handle("/metrics", promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    os.Getenv("PROM_PORT"),
		Handler: mux,
	}

	return metricsServer
}

func initJaeger() (io.Closer, error) {
	cfgMetrics := &jaegerCfg.Configuration{
		ServiceName: "ova-conversation-api",
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := cfgMetrics.NewTracer(
		jaegerCfg.Logger(jaeger.StdLogger),
	)
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}
