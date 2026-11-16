package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

var shareTemplate = template.Must(template.New("share").Parse(`<!doctype html>
<html lang="fr">
<head>
<meta charset="utf-8">
<title>{{.Title}} - Demso</title>
<meta name="description" content="{{.Description}}">
<meta property="og:site_name" content="Demso">
<meta property="og:type" content="website">
<meta property="og:locale" content="fr_FR">
<meta property="og:title" content="{{.Title}}">
<meta property="og:description" content="{{.Description}}">
<meta property="og:url" content="{{.URL}}">
<meta property="og:image" content="{{.Image}}">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="{{.Title}}">
<meta name="twitter:description" content="{{.Description}}">
<meta name="twitter:image" content="{{.Image}}">
<meta http-equiv="refresh" content="0;url={{.URL}}">
</head>
<body><p><a href="{{.URL}}">{{.Title}}</a></p></body>
</html>`))

type ShareHandler struct {
	views   *service.CampaignViewService
	siteURL string
}

// NewShareHandler reads PUBLIC_SITE_URL (the public address of the website,
// for example https://demso.example). In development it defaults to localhost.
func NewShareHandler(views *service.CampaignViewService) *ShareHandler {
	site := strings.TrimRight(os.Getenv("PUBLIC_SITE_URL"), "/")
	if site == "" {
		site = "http://localhost:5173"
	}
	return &ShareHandler{views: views, siteURL: site}
}

// shortText collapses whitespace and cuts the text at max characters.
func shortText(s string, max int) string {
	s = strings.Join(strings.Fields(s), " ")
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	r := []rune(s)
	return strings.TrimSpace(string(r[:max])) + "…"
}

// Page serves a small HTML page for link previews (WhatsApp, Facebook, ...).
func (h *ShareHandler) Page(c *gin.Context) {
	view, err := h.views.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		log.Printf("share page: %v", err)
		c.String(http.StatusInternalServerError, "server error")
		return
	}
	if view == nil {
		c.String(http.StatusNotFound, "not found")
		return
	}

	desc := shortText(view.Description, 160)
	if desc == "" {
		desc = "Soutiens cette cagnotte sur Demso. Ensemble, on va plus loin."
	}

	// A photo flagged as sensitive is never used as a preview image.
	image := h.siteURL + "/og-default.png"
	if view.CoverURL != "" && !view.CoverSensitive {
		image = view.CoverURL
	}

	data := struct{ Title, Description, URL, Image string }{
		Title:       shortText(view.Title, 90),
		Description: desc,
		URL:         h.siteURL + "/c/" + view.URLSlug,
		Image:       image,
	}

	c.Header("Cache-Control", "public, max-age=300")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	if err := shareTemplate.Execute(c.Writer, data); err != nil {
		log.Printf("share page render: %v", err)
	}
}
