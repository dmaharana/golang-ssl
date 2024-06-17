package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	certPath = "./sslcerts/nginx.crt"
	keyPath  = "./sslcerts/nginx.key"
	buildDir = "build"
)

//go:embed build
var buildFS embed.FS

func main() {
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	log.Printf("Loading static files from %s", buildDir)

	// create a subdirectory to serve the embedded files
	staticFilesSubDir, err := fs.Sub(buildFS, buildDir)
	if err != nil {
		log.Fatalf("Error creating subdirectory: %v", err)
	}
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(staticFilesSubDir))))

	// Define a route to handle the POST request
	r.Get("/hw", handleH1Tag)

	// load UI from build directory
	// fs := http.FileServer(http.Dir(buildDir))
	// r.Handle("/*", http.StripPrefix("/", fs))

	// Load the PEM certificate and key
	cert, err := loadPEMCertificate(certPath, keyPath)
	if err != nil {
		log.Println("Error loading PEM certificate and key:", err)
		startHTTPServer(r)
		return
	}

	// Start the HTTPS server
	err = startHTTPSServer(r, cert)
	if err != nil {
		log.Println("Error starting HTTPS server:", err)
		return
	}

}

func startHTTPSServer(r *chi.Mux, cert *tls.Certificate) error {
	// Create the HTTPS server
	server := &http.Server{
		Addr:    ":8443",
		Handler: r,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*cert},
		},
	}

	// Start the HTTPS server
	log.Println("Starting HTTPS server on :8443")
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		log.Println("Error starting HTTPS server:", err)
	}

	return err
}

func startHTTPServer(r *chi.Mux) error {
	// Create the HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start the HTTP server
	log.Println("Starting HTTP server on :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Println("Error starting HTTP server:", err)
	}

	return err
}

func handleH1Tag(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `
		<html>
			<head>
				<title>H1 Tag</title>
			</head>
			<body>
				<h1>H1 Tag</h1>
			</body>
		</html>
		`
	// Read the request body
	h1Tag := r.FormValue("h1_tag")

	// Print the received h1 tag
	fmt.Printf("Received h1 tag: %s\n", h1Tag)
	log.Printf("Request was sent from %s", r.Host)

	// Write a response
	// fmt.Fprintf(w, "Received h1 tag: %s", h1Tag)
	fmt.Fprintf(w, "%s", htmlTemplate)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
