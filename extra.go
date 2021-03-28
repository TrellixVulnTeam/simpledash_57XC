/*
	Copyright 2021 Arsen Musayelyan

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

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
