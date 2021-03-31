package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/manulife-gwam/peppermint-sparkles/backend"
	fileds "github.com/manulife-gwam/peppermint-sparkles/backend/file"
	redisds "github.com/manulife-gwam/peppermint-sparkles/backend/redis"
	"github.com/manulife-gwam/peppermint-sparkles/internal/pcf/vcap"
	"github.com/manulife-gwam/peppermint-sparkles/service"

	log "github.com/sirupsen/logrus"

	bolt "github.com/coreos/bbolt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

//	server flags
var (
	StdListenPortFlag = cli.StringFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Value:   "8080",
		Usage:   "HTTP port to listen on",
		EnvVars: []string{"PSPARKLES_HTTP_PORT"},
	}

	TlsListenPortFlag = cli.StringFlag{
		Name:    "tls-port",
		Value:   "8443",
		Usage:   "HTTPS port to listen on",
		EnvVars: []string{"PSPARKLES_HTTPS_PORT"},
	}

	TlsCertFlag = cli.StringFlag{
		Name:    "tls-cert",
		Usage:   "TLS certificate file for HTTPS",
		EnvVars: []string{"PSPARKLES_TLS_CERT"},
	}

	TlsKeyFlag = cli.StringFlag{
		Name:    "tls-key",
		Usage:   "TLS key file for HTTPS",
		EnvVars: []string{"PSPARKLES_TLS_KEY"},
	}

	DatastoreTypeFlag = cli.StringFlag{
		Name:    "datastore-type",
		Aliases: []string{"dst"},
		Value:   backend.File,
		Usage:   "backend type to be used for storage",
		EnvVars: []string{"PSPARKLES_DS_TYPE"},
	}

	DatastoreFileFlag = cli.StringFlag{
		Name:    "datastore-file",
		Aliases: []string{"dsf"},
		Value:   "/var/lib/peppermint-sparkles/psparkles.db",
		Usage:   "name / location of file for storing secrets",
		EnvVars: []string{"PSPARKLES_DS_FILE"},
	}

	DatastoreAddrFlag = cli.StringFlag{
		Name:    "datastore-addr",
		Aliases: []string{"dsa"},
		Value:   "localhost:6379",
		Usage:   "address for the remote datastore",
		EnvVars: []string{"PSPARKLES_DS_ADDR"},
	}

	Serve = &cli.Command{
		Name:    "server",
		Aliases: []string{"serve"},
		Flags: []cli.Flag{
			&StdListenPortFlag,
			&TlsListenPortFlag,
			&TlsCertFlag,
			&TlsKeyFlag,
			&DatastoreAddrFlag,
			&DatastoreFileFlag,
			&DatastoreTypeFlag,
		},
		Usage: "start the server",

		Action: func(context *cli.Context) error {
			var (
				ds  backend.Datastore
				err error
			)

			dst := context.String(DatastoreTypeFlag.Name)

			switch dst {
			case backend.Redis:
				opts := &redis.Options{Addr: context.String(DatastoreAddrFlag.Name)}

				//	check if running in PCF pull the vcap services if available
				services, err := vcap.GetServices()
				if err != nil {
					return cli.Exit(errors.Wrap(err, "unable to retrieve vcap services"), 1)
				}

				if services != nil {
					if i := services.Tagged(dst); i != nil {
						creds := i.Credentials
						opts = &redis.Options{
							Addr:     fmt.Sprintf("%s:%d", creds["host"].(string), int(creds["port"].(float64))),
							Password: creds["password"].(string),
						}
					}
				}

				if ds, err = redisds.Open(opts); err != nil {
					return cli.Exit(errors.Wrap(err, "unable to open connection to datastore"), 1)
				}

			case backend.File:

				//	FIXME ... include / handle additional bolt options (e.g. timeout, etc)
				fname := context.String(DatastoreFileFlag.Name)
				if ds, err = fileds.Open(fname, bolt.DefaultOptions); err != nil {
					return cli.Exit(errors.Wrap(err, "unable to open connection to datastore"), 1)
				}

			default:
				return cli.Exit(errors.Errorf("%s is not a supported datastore type", dst), 1)
			}

			defer ds.Close()
			log.Debug("datastore opened")

			mux := http.NewServeMux()

			//	attach current service handler
			mux = service.Handle(mux, &service.Handler{Backend: ds})

			//	start HTTPS listener in a seperate go routine since it is a blocking func
			go func() {
				cert, key := context.String(TlsCertFlag.Name), context.String(TlsKeyFlag.Name)
				if len(cert) < 1 || len(key) < 1 {
					return
				}

				if _, err := os.Stat(cert); err != nil {
					log.Error(err, "unable to access TLS cert file")
					return
				}

				if _, err := os.Stat(key); err != nil {
					log.Error(err, "unable to access TLS key file")
					return
				}

				svr := &http.Server{
					Addr:    fmt.Sprintf(":%s", context.String(TlsListenPortFlag.Name)),
					Handler: mux,
					TLSConfig: &tls.Config{
						PreferServerCipherSuites: true,
						CurvePreferences: []tls.CurveID{
							tls.CurveP256,
							tls.X25519,
						},
						MinVersion: tls.VersionTLS12,
						CipherSuites: []uint16{
							tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
							tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
							tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
							tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
							tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
							tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

							// excluding due to no forward secrecy, but leaving
							// as it might be necessary for some clients
							// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
							// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
						},
					},
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  20 * time.Second,
				}

				log.Debug("starting HTTPS listener")
				log.Fatal(svr.ListenAndServeTLS(cert, key))
			}()

			log.Debug("starting HTTP listener")

			addr := fmt.Sprintf(":%s", context.String(StdListenPortFlag.Name))
			log.Fatal(http.ListenAndServe(addr, mux))

			return nil
		},
	}
)
