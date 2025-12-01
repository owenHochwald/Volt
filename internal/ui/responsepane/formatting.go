package responsepane

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

// formatJSON attempts to format JSON content with indentation
func formatJSON(content string) string {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, []byte(content), "", "    ")
	if err != nil {
		return content
	}
	return pretty.String()
}

// highlightContent applies syntax highlighting using chroma
func highlightContent(content, lexer string) string {
	var buf bytes.Buffer
	err := quick.Highlight(&buf, content, lexer, "terminal256", "monokai")
	if err != nil {
		return content
	}
	return buf.String()
}

// getContentLexer maps content type to chroma lexer name
func getContentLexer(contentType string) string {
	switch {
	case strings.Contains(contentType, "application/json"):
		return "json"
	case strings.Contains(contentType, "text/html"):
		return "html"
	case strings.Contains(contentType, "text/plain"):
		return "plaintext"
	case strings.Contains(contentType, "application/xml"), strings.Contains(contentType, "text/xml"):
		return "xml"
	default:
		return ""
	}
}

// formatContentByType formats response body based on content type
func formatContentByType(body, contentType string) string {
	switch {
	case strings.Contains(contentType, "application/json"):
		formatted := formatJSON(body)
		return highlightContent(formatted, "json")

	case strings.Contains(contentType, "image/jpeg"), strings.Contains(contentType, "image/"):
		return fmt.Sprintf("Sorry, we don't support %s yet!", contentType)

	case strings.Contains(contentType, "text/html"):
		return highlightContent(body, "html")

	case strings.Contains(contentType, "text/plain"):
		return highlightContent(body, "plaintext")

	case strings.Contains(contentType, "application/xml"), strings.Contains(contentType, "text/xml"):
		return highlightContent(body, "xml")

	case strings.Contains(contentType, "application/graphql"):
		return "Sorry, we don't support GraphQL yet!"

	case strings.Contains(contentType, "multipart/form-data"):
		return "Sorry, we don't support multipart/form-data yet!"

	default:
		return fmt.Sprintf("Unhandled Content-Type: %s\n", contentType)
	}
}
