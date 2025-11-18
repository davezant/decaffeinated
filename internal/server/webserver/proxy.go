package webserver

import (
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var blocked = []string{
	"facebook.com",
	"pornhub.com",
	"instagram.com",
}

func isBlocked(host string) bool {
	host = strings.ToLower(host)
	for _, b := range blocked {
		if strings.Contains(host, b) {
			return true
		}
	}
	return false
}

func proxyHandler(w http.ResponseWriter, req *http.Request) {
	// ---- HTTPS (CONNECT) ----
	if req.Method == http.MethodConnect {
		host := req.Host
		log.Println("[HTTPS Intercept]", host)

		if isBlocked(host) {
			log.Println("[BLOCKED HTTPS]", host)
			http.Error(w, "Site bloqueado", http.StatusForbidden)
			return
		}

		destConn, err := net.DialTimeout("tcp", host, 5*time.Second)
		if err != nil {
			log.Println("[ERROR] Falha ao conectar destino:", err)
			http.Error(w, "Falha ao conectar", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(200)
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Hijack não suportado", http.StatusInternalServerError)
			return
		}

		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			destConn.Close()
			http.Error(w, "Erro ao hijack", http.StatusInternalServerError)
			return
		}

		// Logando início da conexão
		log.Println("[TUNNEL OPEN]", host)

		go func() {
			io.Copy(destConn, clientConn)
			destConn.Close()
			clientConn.Close()
			log.Println("[TUNNEL CLOSED WRITE]", host)
		}()

		go func() {
			io.Copy(clientConn, destConn)
			destConn.Close()
			clientConn.Close()
			log.Println("[TUNNEL CLOSED READ]", host)
		}()

		return
	}

	// ---- HTTP normal ----
	host := req.URL.Host
	log.Println("[HTTP Intercept]", host)

	if isBlocked(host) {
		log.Println("[BLOCKED HTTP]", host)
		http.Error(w, "Site bloqueado", http.StatusForbidden)
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println("[ERROR] RoundTrip:", err)
		http.Error(w, "Erro ao acessar destino", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	log.Println("[HTTP OK]", host)
}

func StartProxyServer() {
	log.Println("[INFO] proxy: Rodando em :8080")
	log.Println("[INFO] proxy: Sete no navegador: HTTP & HTTPS -> localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(proxyHandler)))
}
