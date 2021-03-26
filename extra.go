package main

import (
	"compress/gzip"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"strings"
)

// Send error to HTTP response
func httpError(res http.ResponseWriter, errTmpl *template.Template, config Conf, statusCode int, reason string) {
	// Write error code to response
	res.WriteHeader(statusCode)
	// Execute error template, outputting to response
	err := errTmpl.Execute(res, map[string]interface{}{
		"StatusCode": statusCode,
		"Reason":     reason,
		"Config":     config,
	})
	if err != nil {
		log.Warn().Err(err).Msg("Error occurred while handling error")
	}
}

func compressRes(res http.ResponseWriter) *gzip.Writer {
	// Set response header to reflect gzip compression
	res.Header().Set("Content-Encoding", "gzip")
	// Wrap response in gzip writer
	return gzip.NewWriter(res)
}

// Check if a string slice contains a string
func strSlcContains(slice []string, str string) bool {
	// For every value in slice
	for _, val := range slice {
		// If value is contained in provided string, return true
		if strings.Contains(str, val) {
			return true
		}
	}
	return false
}
