package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"regexp"
)

// Function to dynamically execute template and return results
func dynamicTemplate(name string, data interface{}) (template.HTML, error) {
	// Create new buffer
	buf := &bytes.Buffer{}
	// Execute template writing to buffer with provided data
	err := templates[name].Execute(buf, data)
	if err != nil {
		return "", nil
	}
	// Return results of template execution
	return template.HTML(buf.String()), nil
}

// Wrap URL with proxy
func wrapProxy(url string) string {
	// Encode URL with base64
	b64url := base64.StdEncoding.EncodeToString([]byte(url))
	// Return /proxy/{url}
	return fmt.Sprint("/proxy/", b64url)
}

// Wrap string in template.JS to unescape JS code
func unescapeJS(s string) template.JS {
	return template.JS(s)
}

// Remove all non-alphanumeric characters
func toAlphaNum(s string) string {
	// Create regex matching everything but alphanumeric
	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	// Remove all matched characters in string, then return
	return regex.ReplaceAllString(s, "")
}

// Function to get template function map
func getFuncMap() template.FuncMap {
	// Return function map with template functions
	return template.FuncMap{
		"dyn_template": dynamicTemplate,
		"proxy":        wrapProxy,
		"unescJS": unescapeJS,
		"toAlphaNum": toAlphaNum,
 	}
}
