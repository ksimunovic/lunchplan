package main

import "net/http"

// Defines a single route, e.g. a human readable name, HTTP method and the
// pattern the function that will execute when the route is called.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Defines the type Routes which is just an array (slice) of Route structs.
type Routes []Route

// Initialize our routes
var routes = Routes{
	Route{
		"Register",
		"POST",
		"/register",
		UserService("Register"),
	}, Route{
		"Login",
		"POST",
		"/login",
		UserService("Login"),
	}, Route{
		"Account",
		"GET",
		"/account",
		UserService("GetAccount"),
	}, Route{
		"Blogs",
		"GET",
		"/blogs",
		UserService("Blogs"),
	}, Route{
		"Blog",
		"POST",
		"/blog",
		UserService("Blog"),
	},
}
