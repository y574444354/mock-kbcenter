package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/zgsm/mock-kbcenter/config"
	"github.com/zgsm/mock-kbcenter/i18n"
	"github.com/zgsm/mock-kbcenter/pkg/httpclient"
)

func Run(cfg *config.Config, workDir string) {
	// Initialize HTTP client
	client, err := httpclient.NewClient(nil)
	if err != nil {
		log.Fatalf("%s", i18n.Translate("proxy.client_init_failed", "", map[string]interface{}{"error": err.Error()}))
	}

	// Create proxy handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, i18n.Translate("proxy.method_not_allowed", "", nil), http.StatusMethodNotAllowed)
			return
		}

		// Get client_id from query params
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			http.Error(w, i18n.Translate("proxy.missing_client_id", "", nil), http.StatusBadRequest)
			return
		}

		// Parse client URL
		clientURL, err := url.Parse(clientID)
		if err != nil {
			http.Error(w, i18n.Translate("proxy.invalid_client_id", "", map[string]interface{}{"error": err.Error()}), http.StatusBadRequest)
			return
		}

		// Create new request
		proxyReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, clientURL.String(), nil)
		if err != nil {
			http.Error(w, i18n.Translate("proxy.request_error", "", map[string]interface{}{"error": err.Error()}), http.StatusInternalServerError)
			return
		}

		// Copy headers
		for name, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}

		// Forward request
		resp, err := client.Get(r.Context(), clientURL.String(), nil)
		if err != nil {
			http.Error(w, i18n.Translate("proxy.forward_error", "", map[string]interface{}{"error": err.Error()}), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// Copy status code and body
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("%s", i18n.Translate("proxy.copy_error", "", map[string]interface{}{"error": err.Error()}))
		}
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Proxy.Port)
	log.Printf("%s", i18n.Translate("proxy.starting", "", map[string]interface{}{"addr": addr}))
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("%s", i18n.Translate("proxy.start_failed", "", map[string]interface{}{"error": err.Error()}))
	}
}
