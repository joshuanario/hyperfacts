package web

import (
	"bytes"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"joshuanario.com/hyperfacts/internal/cache"
)

func Catalog() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		vars := mux.Vars(r)
		id := vars["id"]
		if len(id) <= 0 {
			http.Error(w, "Error: invalid id '"+id+"'", http.StatusBadRequest)
			return
		}
		data, err := cache.MarshalCatalog(id)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func Renderable() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		vars := mux.Vars(r)
		id := vars["id"]
		if len(id) <= 0 {
			http.Error(w, "Error: invalid id '"+id+"'", http.StatusBadRequest)
			return
		}
		hash := vars["hash"]
		if len(hash) <= 0 {
			http.Error(w, "Error: invalid hash", http.StatusBadRequest)
			return
		}

		data, err := cache.MarshalRenderable(id, hash)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func Expressable() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		vars := mux.Vars(r)
		id := vars["id"]
		if len(id) <= 0 {
			http.Error(w, "Error: invalid id '"+id+"'", http.StatusBadRequest)
			return
		}
		parsedquery, err := neturl.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		name, err := neturl.QueryUnescape(parsedquery.Get("name"))
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		documentName, err := neturl.QueryUnescape(parsedquery.Get("documentName"))
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		contextref, err := neturl.QueryUnescape(parsedquery.Get("contextref"))
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		data, err := cache.MarshalExpressable(id, documentName, name, contextref)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func ExpressableDocument() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		vars := mux.Vars(r)
		id := vars["id"]
		if len(id) <= 0 {
			http.Error(w, "Error: invalid id '"+id+"'", http.StatusBadRequest)
			return
		}
		parsedquery, err := neturl.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		documentName, err := neturl.QueryUnescape(parsedquery.Get("documentname"))
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		data, err := cache.MarshalExpressableDocument(id, documentName)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	}
}
func NewRouter() http.Handler {
	r := mux.NewRouter()
	foldersRoute := r.PathPrefix("/folders").Subrouter()
	foldersRoute.HandleFunc("/{id}", Catalog()).Methods("GET")
	projectIDRoute := foldersRoute.PathPrefix("/{id}").Subrouter()
	projectIDRoute.HandleFunc("/facts", Expressable()).Methods("GET")
	projectIDRoute.HandleFunc("/documents", ExpressableDocument()).Methods("GET")
	projectIDRoute.HandleFunc("/{hash}", Renderable()).Methods("GET")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	browserDir := filepath.Join(wd, "king-of-roses")
	// Register "/" before the catch-all file server; mux uses the first match.
	r.Path("/").HandlerFunc(OfflineConceptNetworkBrowser(browserDir)).Methods("GET")
	return r
}

func OfflineConceptNetworkBrowser(browserDir string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Error: incorrect verb, "+r.Method, http.StatusInternalServerError)
			return
		}
		htmlBytes, err := os.ReadFile(filepath.Join(browserDir, "index.html"))
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.URL.Query().Get("id")
		if len(id) > 0 {
			data, err := cache.MarshalSuperRenderable(id)
			if err != nil {
				http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			htmlBytes = injectData(htmlBytes, data)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(htmlBytes)
	}
}

// injectData inserts a script defining globalThis.SUPER_RENDERABLE just before </head>
// (or at the start of the document if there is no </head>). data must come
// from json.Marshal, which escapes <, > and & so it cannot close the
// script tag early.
func injectData(htmlBytes, data []byte) []byte {
	var script bytes.Buffer
	script.WriteString("<script>\nglobalThis.SUPER_RENDERABLE = ")
	script.Write(data)
	script.WriteString(";\n</script>\n")

	idx := bytes.Index(htmlBytes, []byte("</head>"))
	if idx < 0 {
		idx = 0
	}
	out := make([]byte, 0, len(htmlBytes)+script.Len())
	out = append(out, htmlBytes[:idx]...)
	out = append(out, script.Bytes()...)
	out = append(out, htmlBytes[idx:]...)
	return out
}
