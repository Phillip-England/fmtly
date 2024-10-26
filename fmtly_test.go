package fmtly

import (
	"strings"
	"testing"
)

func Doc(head string, body string) string {
	return `
		<html>
			` + head + `
			` + body + `
		</html>
	`
}

func Root(children ...string) string {
	return Doc(`
		<head>
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<script src="/static/js/index.js"></script>
			<link rel='stylesheet' href="/static/css/output.css"></link>
		</head>
	`, `
		<body>
			<div id='root'>`+strings.Join(children, "")+`</div>
		</body>
	`)
}

func TestMain(t *testing.T) {
	// mux, gCtx := vbf.VeryBestFramework()

	// vbf.HandleFavicon(mux)
	// vbf.HandleStaticFiles(mux)

	// vbf.AddRoute("GET /", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {
	// 	vbf.WriteHTML(w, Root())
	// }, vbf.MwLogger)

	// err := vbf.Serve(mux, "8080")
	// if err != nil {
	// 	panic(err)
	// }




}
