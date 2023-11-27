package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"films-api.rdelgado.es/src/internals/models"
)

type (
	// struct for holding response details
	responseData struct {
		status int
		size   int
	}

	// http.ResponseWriter to log request and responses
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract and parse token to get the user_id
		// If token does not exists => continue
		// If token exists => check user in BD and add user_id to context

		tokenString, err := app.tokens.ExtractToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		id, err := app.tokens.VerifyToken(tokenString)
		if err != nil {

			if errors.Is(err, models.ErrInvalidToken) {
				app.clientError(w, http.StatusBadRequest, err)
			} else {
				app.clientError(w, http.StatusUnauthorized, err)
			}

			return
		}
		// Otherwise, we check to see if a user with that ID exists in our database.
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// Save user_id to context
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			ctx = context.WithValue(ctx, userIdContextKey, id)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			app.clientError(w, http.StatusUnauthorized, errors.New("not authenticated"))
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

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

		// Log request received
		app.logger.Info("CLIENT -> API",
			"addr", ip,
			"proto", proto,
			"method", method,
			"uri", uri)

		// Serve HTTP request to rest of middleware and handlers
		next.ServeHTTP(w, r)
	})
}

func (app *application) logResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
			start  = time.Now()
		)

		logResponseWriter := &loggingResponseWriter{
			ResponseWriter: w,
			responseData: &responseData{
				status: 0,
				size:   0,
			},
		}

		// Serve HTTP request to rest of middleware and handlers
		next.ServeHTTP(logResponseWriter, r)

		// After serving the request, log the response to be sent
		userId := r.Context().Value(userIdContextKey)
		duration := time.Since(start)
		app.logger.Info("CLIENT <- API",
			"addr", ip,
			"proto", proto,
			"method", method,
			"uri", uri,
			"size", logResponseWriter.responseData.size,
			"status", logResponseWriter.responseData.status,
			"duration", duration.String(),
			"userId", userId)
	})
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.responseData.size += size            // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.responseData.status = statusCode       // capture status code
}
