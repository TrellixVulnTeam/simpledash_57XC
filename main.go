package main

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gorilla/mux"
	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/wader/gormstore/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Set global logger to ConsoleWriter
var Log = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

type App struct {
	Route     *mux.Router
	Templates map[string]*template.Template
	Session   *gormstore.Store
	Config    Conf
}

// Create new empty map to store templates
var templates = map[string]*template.Template{}

func main() {
	// Create command-line flags
	addr := flag.IPP("addr", "a", net.ParseIP("0.0.0.0"), "Bind address for HTTP server")
	port := flag.IntP("port", "p", 8080, "Bind port for HTTP server")
	config := flag.StringP("config", "c", "simpledash.toml", "TOML config file")
	// Parse flags
	flag.Parse()

	// Create new router
	router := mux.NewRouter().StrictSlash(true)

	// Create OS-specific glob for all templates
	path := filepath.Join("resources", "templates", "*.html")
	// Create OS-specific glob for all card templates
	cardGlob := filepath.Join("resources", "templates", "cards", "*.html")
	// Get all template paths
	tmplMatches, _ := filepath.Glob(path)
	cardMatches, _ := filepath.Glob(cardGlob)
	matches := append(tmplMatches, cardMatches...)
	// For each template path
	for _, match := range matches {
		// Get name of file without path or extension
		fileName := strings.TrimSuffix(filepath.Base(match), filepath.Ext(match))
		// If file is called base
		if fileName == "base" {
			// Skip
			continue
		}
		var err error
		// Parse detected template and base template, add to templates map
		templates[fileName], err = template.New(
			filepath.Base(match)).Funcs(
			sprig.FuncMap()).Funcs(
			getFuncMap()).ParseFiles(
			"resources/templates/base.html", match)
		if err != nil {
			Log.Fatal().Str("template", fileName).Err(err).Msg("Error parsing template")
		}
	}

	// Open sqlite database called sessions.db for storing sessions
	sessionDB, _ := gorm.Open(sqlite.Open("sessions.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	// Create session store from database
	sessionStore := gormstore.New(sessionDB, []byte(""))

	// Create channel to stop periodic cleanup
	quitCleanup := make(chan struct{})
	// Clean up expired sessions every hour
	go sessionStore.PeriodicCleanup(1*time.Hour, quitCleanup)

	// Open config file
	configFile, err := os.Open(filepath.Clean(*config))
	if err != nil {
		Log.Fatal().Err(err).Msg("Error opening config file")
	}
	// Create new TOML decoder
	dec := toml.NewDecoder(configFile)
	// Create new nil variable to store decoded config
	var decodedConf Conf
	// Decode config into variable
	err = dec.Decode(&decodedConf)
	if err != nil {
		Log.Fatal().Err(err).Msg("Error decoding config file")
	}

	// Register HTTP routes
	registerRoutes(App{
		Route:     router,
		Templates: templates,
		Session:   sessionStore,
		Config:    decodedConf,
	})

	// Create string address from flag values
	strAddr := fmt.Sprint(*addr, ":", *port)
	// Create listener on IPv4 using address created above
	ln, err := net.Listen("tcp4", strAddr)
	if err != nil {
		Log.Fatal().Err(err).Msg("Error creating listener")
	}

	// Log HTTP server start
	Log.Info().Str("addr", strAddr).Msg("Starting HTTP server")
	// Start HTTP server using previously-created router and listener
	err = http.Serve(ln, router)
	if err != nil {
		Log.Fatal().Err(err).Msg("Error while serving")
	}
}
