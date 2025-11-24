package config

import (
	"fmt"
	"log" // Pour logger les informations ou erreurs de chargement de config

	"github.com/spf13/viper" // La bibliothèque pour la gestion de configuration
)

// TODO Créer Config qui est la structure principale qui mappe l'intégralité de la configuration de l'application.
// Les tags `mapstructure` sont utilisés par Viper pour mapper les clés du fichier de config
// (ou des variables d'environnement) aux champs de la structure Go.
type Config struct {
	Server            ServerConfig      `mapstructure:"server"`
	Database          DatabaseConfig    `mapstructure:"database"`
	Analytics         AnalyticsConfig   `mapstructure:"analytics"`
	Monitor           MonitorConfig     `mapstructure:"monitor"`
	ClickWorkerConfig ClickWorkerConfig `mapstructure:"clickworker"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	BaseURL string `mapstructure:"base_url"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type AnalyticsConfig struct {
	Enabled    bool `mapstructure:"enabled"`
	BufferSize int  `mapstructure:"buffer_size"`
}

type MonitorConfig struct {
	Enabled         bool `mapstructure:"enabled"`
	IntervalMinutes int  `mapstructure:"interval_minutes"`
}
type ClickWorkerConfig struct {
	ChannelSize int `mapstructure:"channel_size"`
	WorkerCount int `mapstructure:"worker_count"`
}

// LoadConfig charge la configuration de l'application en utilisant Viper.
// Elle recherche un fichier 'config.yaml' dans le dossier 'configs/'.
// Elle définit également des valeurs par défaut si le fichier de config est absent ou incomplet.
func LoadConfig() (*Config, error) {
	// TODO Spécifie le chemin où Viper doit chercher les fichiers de config.
	// on cherche dans le dossier 'configs' relatif au répertoire d'exécution.
	viper.AddConfigPath("configs/")

	// TODO Spécifie le nom du fichier de config (sans l'extension).
	viper.SetConfigName("config")

	// TODO Spécifie le type de fichier de config.
	viper.SetConfigType("yaml")

	// TODO : Définir les valeurs par défaut pour toutes les options de configuration.
	// Ces valeurs seront utilisées si les clés correspondantes ne sont pas trouvées dans le fichier de config
	// ou si le fichier n'existe pas.
	// server.port, server.base_url etc.
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.base_url", "http://localhost")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "mydb")
	viper.SetDefault("database.user", "user")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("analytics.enabled", true)
	viper.SetDefault("analytics.buffer_size", 100)
	viper.SetDefault("monitor.enabled", true)
	viper.SetDefault("monitor.interval_minutes", 5)
	viper.SetDefault("clickworker.worker_count", 1)

	// TODO : Lire le fichier de configuration.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No config file found, using default values")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("fatal error config file: %w", err)
		}
	}

	// TODO 4: Démapper (unmarshal) la configuration lue (ou les valeurs par défaut) dans la structure Config.
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Log  pour vérifier la config chargée
	log.Printf("Configuration loaded: Server Port=%d, DB Name=%s, Analytics Buffer=%d, Monitor Interval=%dmin",
		cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)

	return &cfg, nil // Retourne la configuration chargée
}
