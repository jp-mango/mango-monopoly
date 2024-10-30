package main

import "net/http"

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//used to restrict where the resources for web page can be loaded from.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		// used to control what information is included in a Referer header when a user navigates away from your web page.
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// instructs browsers to not MIME-type sniff the content-type of the response, which in turn helps to prevent content-sniffing attacks.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// used to help prevent clickjacking attacks in older browsers that don’t support CSP headers.
		w.Header().Set("X-Frame-Options", "deny")

		// used to disable the blocking of cross-site scripting attacks.
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}
