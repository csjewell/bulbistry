package bulbistry

import(
    "context"
    "net/http"

    htpasswd "github.com/tg123/go-htpasswd"
)

type contextKey int
const (
    HTPasswordKey contextKey = iota
    userKey
)

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
	f := ctx.Value(HTPasswordKey);

        if f == nil {
            next.ServeHTTP(w, r)
            return
        } 

	myauth, err := htpasswd.New(f.(string), htpasswd.DefaultSystems, nil)
        if err != nil {
            w.Write(ConfigError(err))
            http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
        username, password, ok := r.BasicAuth()
        if ok {
            ok := myauth.Match(username, password)
		
            if ok {
                ctxWithUser := context.WithValue(ctx, userKey, username)
                rWithUser := r.WithContext(ctxWithUser)
                next.ServeHTTP(w, rWithUser)
            } else {
                w.Header().Set("WWW-Authenticate", `Basic realm="Bulbistry OCI Artifact Registry", charset="UTF-8"`)
                w.Write(InvalidLogin())
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
	    }
        } else {
            w.Header().Set("WWW-Authenticate", `Basic realm="Bulbistry OCI Artifact Registry", charset="UTF-8"`)
            w.Write(NoLogin())
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
        }
    })
}   

func V2Check (w http.ResponseWriter, e *http.Request) {
    w.Header().Set("Content-Type", `text/plain`)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"));
}
