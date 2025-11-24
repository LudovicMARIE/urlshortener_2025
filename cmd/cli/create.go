package cli

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"

	// Pour valider le format de l'URL

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/axellelanca/urlshortener/internal/utils"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// Driver SQLite pour GORM
)

// TODO : Faire une variable longURLFlag qui stockera la valeur du flag --url
var (
	longURLFlag string
	err         error
)

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO 1: Valider que le flag --url a été fourni.
		// On vérifie si le nom à été fourni via les flags de la commande, si non on le demande directement via un reader
		if longURLFlag == "" {
			fmt.Print("➡️  Renseigner l'URL à convertir : ")
			reader := bufio.NewReader(os.Stdin)
			longURLFlag, err = utils.ReaderLine(reader)
			if err != nil {
				log.Fatalf("‼️  Erreur de saisi: %v \n", err)
			}
		}

		// TODO Validation basique du format de l'URL avec le package url et la fonction ParseRequestURI
		// si erreur, os.Exit(1)
		_, err := url.ParseRequestURI(longURLFlag)
		if err != nil {
			log.Fatalf("‼️  Le format de l'URL n'est pas valide : %v \n", err)
			os.Exit(1)
		}

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg

		// TODO : Initialiser la connexion à la base de données SQLite.
		db, dbGormOk := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if dbGormOk != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}

		sqlDB, dbSqlOk := db.DB()
		if dbSqlOk != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		// TODO S'assurer que la connexion est fermée à la fin de l'exécution de la commande
		defer sqlDB.Close()

		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)

		linkService := services.NewLinkService(linkRepo, clickRepo)

		// TODO : Appeler le LinkService et la fonction CreateLink pour créer le lien court.
		// os.Exit(1) si erreur
		link, isLinkCreated := linkService.CreateLink(longURLFlag)
		if isLinkCreated != nil {
			log.Fatalf("‼️  Erreur lors de la création de l'URL courte: %v", err)
			os.Exit(1)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.Shortcode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.Shortcode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// TODO : Définir le flag --url pour la commande create.
	CreateCmd.Flags().StringVarP(&longURLFlag, "url", "u", "", "URL longue a convertir en URL courte")

	// TODO :  Marquer le flag comme requis
	CreateCmd.MarkFlagRequired("url")

	// TODO : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(CreateCmd)
}
