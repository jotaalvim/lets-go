package main

type contextKey string

// my own key was made to stop key collision
const isAuthenticatedContextKey = contextKey("isAuthenticated")
