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
		ServiceCall("Register", LoadConfiguration().UserService.Host),
	}, Route{
		"Login",
		"POST",
		"/login",
		ServiceCall("Login", LoadConfiguration().UserService.Host),
	}, Route{
		"Account",
		"GET",
		"/account",
		ServiceCall("GetAccount", LoadConfiguration().UserService.Host),
	},
	Route{
		"Create",
		"POST",
		"/meal",
		ServiceCall("Create", LoadConfiguration().MealService.Host),
	},
	Route{
		"GetAllUserMeals",
		"GET",
		"/meal/all",
		ServiceCall("GetAllUserMeals", LoadConfiguration().MealService.Host),
	},
	Route{
		"Suggest",
		"GET",
		"/meal/suggest",
		ServiceCall("Suggest", LoadConfiguration().MealService.Host),
	},
	Route{
		"Read",
		"GET",
		"/meal/{id}",
		ServiceCall("Read", LoadConfiguration().MealService.Host),
	},
	Route{
		"Update",
		"POST",
		"/meal/{id}",
		ServiceCall("Update", LoadConfiguration().MealService.Host),
	},
	Route{
		"Delete",
		"DELETE",
		"/meal/{id}",
		ServiceCall("Delete", LoadConfiguration().MealService.Host),
	},
	Route{
		"Create",
		"POST",
		"/tag",
		ServiceCall("Create", LoadConfiguration().TagService.Host),
	},
	Route{
		"GetAllUserTags",
		"GET",
		"/tag/all",
		ServiceCall("GetAllUserTags", LoadConfiguration().TagService.Host),
	},
	Route{
		"Read",
		"GET",
		"/tag/{id}",
		ServiceCall("Read", LoadConfiguration().TagService.Host),
	},
	Route{
		"Update",
		"POST",
		"/tag/{id}",
		ServiceCall("Update", LoadConfiguration().TagService.Host),
	},
	Route{
		"Delete",
		"DELETE",
		"/tag/{id}",
		ServiceCall("Delete", LoadConfiguration().TagService.Host),
	},
	Route{
		"Create",
		"POST",
		"/calendar",
		ServiceCall("Create", LoadConfiguration().CalendarService.Host),
	},
	Route{
		"GetAllUserCalendars",
		"GET",
		"/calendar/all",
		ServiceCall("GetAllUserCalendars", LoadConfiguration().CalendarService.Host),
	},
	Route{
		"Read",
		"GET",
		"/calendar/{id}",
		ServiceCall("Read", LoadConfiguration().CalendarService.Host),
	},
	Route{
		"Update",
		"POST",
		"/calendar/{id}",
		ServiceCall("Update", LoadConfiguration().CalendarService.Host),
	},
	Route{
		"Delete",
		"DELETE",
		"/calendar/{id}",
		ServiceCall("Delete", LoadConfiguration().CalendarService.Host),
	},
}
