package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Create struct to store template context
type TemplateData struct {
	Username string
	Config   Conf
	User     User
	Error    string
}

func registerRoutes(app App) {
	// Root endpoint, home page
	app.Route.Path("/").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Get session by name from config
		session, _ := app.Session.Get(req, app.Config.Session.Name)
		// Attempt to get loggedInAs from session
		loggedInAs, ok := session.Values["loggedInAs"].(string)
		// If user not logged in and login is required
		if !ok && app.Config.LoginRequired {
			// Redirect to login page
			http.Redirect(res, req, "/login", http.StatusFound)
		} else if !ok && !app.Config.LoginRequired {
			// If not logged in and login is not required
			// Set logged in to public user
			loggedInAs = "_public_"
			// Set logged in user to session
			session.Values["loggedInAs"] = loggedInAs
			// Save session
			_ = session.Save(req, res)
		}
		// Create template context
		tmplData := TemplateData{Username: loggedInAs, Config: app.Config, User: app.Config.Users[loggedInAs]}
		// Execute home template with provided context
		err := app.Templates["home"].Execute(res, tmplData)
		if err != nil {
			httpError(res, app.Templates["error"], app.Config, http.StatusInternalServerError, "Error executing home template")
			Log.Warn().Err(err).Msg("Error executing home template")
			return
		}
	})

	// /login endpoint, login page
	app.Route.Path("/login").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Check request method
		switch req.Method {
		// If GET request
		case http.MethodGet:
			// Get session by name from config
			session, _ := app.Session.Get(req, app.Config.Session.Name)
			// Attempt to get loggedInAs from session
			loggedInAs, ok := session.Values["loggedInAs"].(string)
			// If logged in and not as public user
			if ok && loggedInAs != "_public_" {
				// Redirect back to home page
				http.Redirect(res, req, "/", http.StatusFound)
			}
			// Get query parameter error
			urlErr := req.URL.Query().Get("error")
			// Execute login template with provided context
			err := app.Templates["login"].Execute(res, TemplateData{Config: app.Config, Error: urlErr, Username: loggedInAs})
			if err != nil {
				httpError(res, app.Templates["error"], app.Config, http.StatusInternalServerError, "Error executing login template")
				Log.Warn().Err(err).Msg("Error executing login template")
				return
			}
		// If POST request
		case http.MethodPost:
			// Parse form in POST request body
			_ = req.ParseForm()
			// Get password from form
			password := req.PostForm.Get("password")
			// Get username from form
			username := req.PostForm.Get("username")
			// Get user from config by username
			user, ok := app.Config.Users[username]
			// If user not found
			if !ok {
				// Redirect to login page with error parameter set to usr
				http.Redirect(res, req, "/login?error=usr", http.StatusFound)
				return
			}
			// Compare hash stored in config to password in form
			err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
			// If password was incorrect
			if err != nil {
				// Redirect to login page with error parameter set to pwd
				http.Redirect(res, req, "/login?error=pwd", http.StatusFound)
				return
			}
			// Get session by name from config
			session, _ := app.Session.Get(req, app.Config.Session.Name)
			// Set loggedInAs value in session to username from form
			session.Values["loggedInAs"] = username
			// Save session
			err = session.Save(req, res)
			if err != nil {
				httpError(res, app.Templates["error"], app.Config, http.StatusInternalServerError, "Error saving session")
				Log.Warn().Err(err).Msg("Error saving session")
				return
			}
			// Redirect to homepage
			http.Redirect(res, req, "/", http.StatusFound)
		}
	})

	// /logout endpoint, logout and redirect
	app.Route.Path("/logout").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Get session by name from config
		session, _ := app.Session.Get(req, app.Config.Session.Name)
		// Remove loggedInAs value from session
		delete(session.Values, "loggedInAs")
		// Save session
		err := session.Save(req, res)
		if err != nil {
			httpError(res, app.Templates["error"], app.Config, http.StatusInternalServerError, "Error saving session")
			Log.Warn().Err(err).Msg("Error while handling error")
			return
		}
		// Redirect to login page
		http.Redirect(res, req, "/login", http.StatusFound)
	})

	// /status/:b64url endpoint, return status of base64 encoded URL as JSON
	app.Route.Path("/status/{b64url}").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Get path variables
		vars := mux.Vars(req)
		// Create JSON encoder writing to response
		enc := json.NewEncoder(res)
		// Decode base64 URL string
		url, _ := base64.StdEncoding.DecodeString(vars["b64url"])
		// Create new HEAD request to check status without downloading whole page
		headReq, err := http.NewRequest(http.MethodHead, string(url), nil)
		if err != nil {
			// Encode down status and write to response
			_ = enc.Encode(map[string]interface{}{"code": 0, "down": true})
			Log.Warn().Err(err).Msg("Error creating HEAD request for status check")
			return
		}
		// Create new HTTP client with 5 second timeout
		client := http.Client{Timeout: 5 * time.Second}
		// Use client to do request created above
		headRes, err := client.Do(headReq)
		if err != nil {
			// Encode down status and write to response
			_ = enc.Encode(map[string]interface{}{"code": 0, "down": true})
			Log.Warn().Err(err).Msg("Error executing HEAD request for status check")
			return
		}
		// Encode returned status and write to response
		_ = enc.Encode(map[string]interface{}{"code": headRes.StatusCode, "down": false})
	})

	// /proxy/:b64url endpoint, proxy HTTP connection bypassing CORS
	app.Route.Path("/proxy/{b64url}").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Get path variables
		vars := mux.Vars(req)
		// Decode base64 URL string
		url, _ := base64.StdEncoding.DecodeString(vars["b64url"])
		// If URL is allowed for proxy
		if strSlcContains(app.Config.AllowProxy, string(url)) {
			// Create new HTTP request with the same parameters as sent to endpoint
			proxyReq, err := http.NewRequest(req.Method, string(url), req.Body)
			if err != nil {
				httpError(res, app.Templates["error"], app.Config, http.StatusInternalServerError, "Proxying connection failed")
				Log.Warn().Err(err).Msg("Error creating request for proxy")
				return
			}
			// Create new HTTP client with 5 second timeout
			client := http.Client{Timeout: 5 * time.Second}
			// Use client to do request created above
			proxyRes, err := client.Do(proxyReq)
			if err != nil {
				httpError(res, app.Templates["error"], app.Config, http.StatusInternalServerError, "Proxying connection failed")
				Log.Warn().Err(err).Msg("Error executing request for proxy")
				return
			}
			// Close proxy response body at end of function
			defer proxyRes.Body.Close()
			// Copy data from proxy response to response
			io.Copy(res, proxyRes.Body)
		} else {
			httpError(res, app.Templates["error"], app.Config, http.StatusBadRequest, "This url is not in allowedProxy")
			Log.Warn().Str("url", string(url)).Msg("URL is disallowed for proxy")
			return
		}
	})

	// Catch-all route, gzip-compressing file server
	app.Route.PathPrefix("/").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Create OS-specific path to requested file
		filePath := filepath.Join("resources", "public", req.URL.Path)
		// Open requested file
		file, err := os.Open(filePath)
		if err != nil {
			Log.Warn().Str("file", filePath).Msg("File not found")
			httpError(res, app.Templates["error"], app.Config, http.StatusNotFound, "This file was not found")
			return
		}
		// Close file at end of function
		defer file.Close()
		// Compress response
		gzRes := compressRes(res)
		// Close compressed response at end of function
		defer gzRes.Close()
		// Copy file contents to compressed response
		io.Copy(gzRes, file)
	})
}
