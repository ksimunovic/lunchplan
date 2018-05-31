package main

import "net/http"
import "./Controllers"

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
		"Index",
		"GET",
		"/",
		Controller.MealController.Index(),
	},
	Route{
		"Index",
		"GET",
		"/calendar",
		Controller.CalendarController.ShowCalendar(),
	},
	Route{
		"Login",
		"GET",
		"/login",
		Controller.UserController.Login(),
	},

	Route{
		"Login",
		"POST",
		"/login",
		Controller.UserController.ProcessLogin(),
	},
}
