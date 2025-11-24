package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	api "github.com/axellelanca/urlshortener/internal/api"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/monitor"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/axellelanca/urlshortener/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// Driver SQLite pour GORM
)

var err error

// RunServerCmd représente la commande 'run-server' de Cobra.
// C'est le point d'entrée pour lancer le serveur de l'application.
var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs et les processus de fond.",
	Long: `Cette commande initialise la base de données, configure les APIs,
démarre les workers asynchrones pour les clics et le moniteur d'URLs,
puis lance le serveur HTTP.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO : créer une variable qui stock la configuration chargée globalement via cmd.cfg
		// Ne pas oublier la gestion d'erreur et faire un fatalF

		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("FATAL: La configuration n'a pas été chargée correctement")
		}

		// Initialiser la connexion à la BDD
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}
		// Fermer la connexion SQL quand la commande se termine.
		defer sqlDB.Close()

		// TODO : Initialiser les repositories.
		// Créez des instances de GormLinkRepository et GormClickRepository.
		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)

		// Laissez le log
		log.Println("Repositories initialisés.")

		// Initialiser les services métiers.
		// Créez des instances de LinkService et ClickService, en leur passant les repositories nécessaires.
		linkService := services.NewLinkService(linkRepo, clickRepo)
		clickService := services.NewClickService(clickRepo)
		_ = clickService

		// Laissez le log
		log.Println("Services métiers initialisés.")

		// Initialiser le channel ClickEventsChannel (dans api) et lancer les workers de clics.
		// Le channel est bufferisé avec la taille configurée dans la configuration analytics.buffer_size.
		if api.ClickEventsChannel == nil {
			api.ClickEventsChannel = make(chan models.ClickEvent, cfg.Analytics.BufferSize)
		}
		// Démarrer les workers de clics en arrière-plan.
		workers.StartClickWorkers(cfg.ClickWorkerConfig.WorkerCount, api.ClickEventsChannel, clickRepo)

		// Log informatif
		log.Printf("Channel d'événements de clic initialisé avec un buffer de %d. %d worker(s) de clics démarré(s).",
			cfg.Analytics.BufferSize, cfg.ClickWorkerConfig.WorkerCount)

		// Initialiser et lancer le moniteur d'URLs si activé dans la configuration.
		monitorInterval := time.Duration(cfg.Monitor.IntervalMinutes) * time.Minute
		if cfg.Monitor.Enabled {
			urlMonitor := monitor.NewUrlMonitor(linkRepo, monitorInterval)
			go urlMonitor.Start()
			log.Printf("Moniteur d'URLs démarré avec un intervalle de %v.", monitorInterval)
		} else {
			log.Println("Moniteur d'URLs désactivé par la configuration.")
		}

		// Configurer le routeur Gin et les handlers API.
		router := gin.Default()
		api.SetupRoutes(router, linkService)

		// Pas toucher au log
		log.Println("Routes API configurées.")

		// Créer le serveur HTTP Gin
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		// Démarrer le serveur Gin dans une goroutine anonyme pour ne pas bloquer.
		go func() {
			log.Printf("Démarrage du serveur HTTP sur %s...", serverAddr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("FATAL: Erreur du serveur HTTP: %v", err)
			}
		}()

		// Gère l'arrêt propre du serveur (graceful shutdown).
		// Crée un channel pour les signaux OS (SIGINT, SIGTERM), bufferisé à 1.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Attendre Ctrl+C ou signal d'arrêt

		// Bloquer jusqu'à ce qu'un signal d'arrêt soit reçu.
		<-quit
		log.Println("Signal d'arrêt reçu. Arrêt du serveur...")

		// Arrêt propre du serveur HTTP avec un timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Erreur lors de l'arrêt du serveur: %v", err)
		}

		// Fermer le channel des événements pour indiquer aux workers d'arrêter proprement.
		if api.ClickEventsChannel != nil {
			close(api.ClickEventsChannel)
		}

		// Donner un court délai pour que les workers terminent proprement.
		log.Println("Arrêt en cours... Donnez un peu de temps aux workers pour finir.")
		time.Sleep(2 * time.Second)

		log.Println("Serveur arrêté proprement.")
	},
}

func init() {
	// Ajouter la commande 'run-server' à la racine Cobra
	cmd2.RootCmd.AddCommand(RunServerCmd)
}
