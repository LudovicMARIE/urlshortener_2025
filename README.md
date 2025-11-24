# 🚀 Projet Final Go — URL Shortener

**README – Version de notre groupe**

## 🎯 Objectif du Projet

Dans le cadre du projet final du module Go, notre groupe a développé un **service web complet de raccourcissement d’URL**, capable de :

* Générer des URLs courtes uniques,
* Rediriger instantanément vers l’URL longue,
* Enregistrer les clics de manière totalement **asynchrone**,
* Surveiller régulièrement l’état des URLs,
* Fournir une API REST complète,
* Proposer une interface CLI fonctionnelle (via Cobra).

L’objectif était de mettre en pratique l’ensemble des notions vues durant le module Go : concurrence, gestion d’API, ORM, CLI, configuration, architecture propre, multithreading, etc.

---

## 🧠 Compétences et Technologies Utilisées

Notre projet mobilise les concepts suivants :

* Syntaxe Go (structs, interfaces, erreurs, maps…)
* Concurrence : **Goroutines**, **Channels**, workers asynchrones
* API REST avec **Gin**
* CLI avec **Cobra**
* ORM **GORM** + SQLite
* Configuration dynamique via **Viper**
* Patterns d’architecture (Repository, Service)
* Manipulation JSON et gestion d’erreurs propre
* Monitoring d’URLs avec tâches planifiées

---

## 🧩 Fonctionnalités Développées

### ✔️ Fonctionnalités essentielles

1. **Raccourcissement d’URL**

   * Génération de codes courts uniques (6 caractères)
   * Gestion des collisions via retry

2. **Redirection instantanée**

   * HTTP **302**
   * Enregistrement des clics en asynchrone (via channel bufferisé + worker dédié)

3. **Monitoring automatique d’URLs**

   * Vérification périodique configurable (via Viper)
   * Notification dans les logs en cas de changement d’état

4. **API REST (Gin)**

   * `GET /health`
   * `POST /api/v1/links`
   * `GET /{shortCode}` redirection + analytics async
   * `GET /api/v1/links/{shortCode}/stats`

5. **CLI (Cobra)**

   * `run-server` → démarre serveur + workers + monitor
   * `create --url="..."` → crée une URL courte
   * `stats --code="xyz123"` → stats d’un lien
   * `migrate` → migrations de la base

### ⭐ Bonus potentiels (si temps disponible)

* Alias personnalisés
* Expiration automatique
* Rate limiting

---

## 🗂️ Architecture du Projet

Nous avons utilisé une architecture modulaire afin de bien séparer les responsabilités :

```
url-shortener/
├── cmd/
│   ├── root.go
│   ├── server/server.go
│   └── cli/
│       ├── create.go
│       ├── stats.go
│       └── migrate.go
├── internal/
│   ├── api/handlers.go
│   ├── models/
│   │   ├── link.go
│   │   └── click.go
│   ├── services/
│   │   ├── link_service.go
│   │   └── click_service.go
│   ├── workers/click_worker.go
│   ├── monitor/url_monitor.go
│   ├── config/config.go
│   ├── repository/
│   │   ├── link_repository.go
│   │   └── click_repository.go
│   └── utils/
│       └── validation.go   <-- Ajouté par notre groupe
├── configs/config.yaml
├── go.mod
├── go.sum
└── README.md
```

---

## 🆕 ⭐ Ajout spécifique de notre groupe : `internal/utils/validation.go`

Nous avons ajouté un dossier supplémentaire `utils/` dans `internal/`, dédié aux fonctions génériques de validation et nettoyage des données.

Exemple de fonction présente dans `validation.go` :

```go
// Permet de lire et nettoyer la saisie d'un utilisateur depuis un reader
func ReaderLine(reader *bufio.Reader) (string, error) {
    readerValue, _ := reader.ReadString('\n')
    readerValue = strings.TrimSpace(readerValue)
    return readerValue, nil
}
```

Ce fichier nous permet :

* d’éviter la duplication de fonctions utilitaires,
* de centraliser tout ce qui concerne la validation / nettoyage des entrées,
* d'améliorer la lisibilité de la CLI et des services.

---

## ▶️ Installation & Utilisation

### 1. Cloner le projet

```bash
git clone https://github.com/LudovicMARIE/urlshortener_2025
cd urlshortener_2025
```

### 2. Gestion des dépendances

```bash
go mod tidy
```

### 3. Compilation

```bash
go build -o url-shortener
```

### 4. Migrations

```bash
./url-shortener migrate
```

### 5. Lancer le serveur

```bash
./url-shortener run-server
```

---

## 🧪 Tests du service

### Créer une URL courte

```bash
./url-shortener create --url="https://example.com"
```

### Accéder à l’URL courte

→ Ouvrir :

```
http://localhost:8080/XYZ123
```

### Consulter les statistiques

```bash
./url-shortener stats --code="XYZ123"
```

### Vérifier l’état du serveur

```bash
curl http://localhost:8080/health
```

### Observer le moniteur d’URLs

→ Logs automatiques toutes les X minutes (X étant paramétrable dans le fichier config.yaml de viper). 

---

## 📝 Barème (rappel prof)

Conforme au barème du projet, notre README et l'organisation du code respectent :

* Architecture claire
* Concurrence implémentée correctement
* API fonctionnelle
* CLI utilisable
* Validation centralisée (via notre ajout utils/)
* Documentation complète

---

## ✏️ Auteurs

Groupe 6 :
* Valérie Song
* Ludovic Marie
* Mathias Mousset

Bon courage pour les corrections :) 

Merci pour votre implication au cours du module, c'était tip top.
