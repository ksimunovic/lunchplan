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
		ServiceCall("Register", LoadConfiguration().UserService.Port),
	}, Route{
		"Login",
		"POST",
		"/login",
		ServiceCall("Login", LoadConfiguration().UserService.Port),
	}, Route{
		"Account",
		"GET",
		"/account",
		ServiceCall("GetAccount", LoadConfiguration().UserService.Port),
	},
	Route{
		"Create",
		"POST",
		"/meal",
		ServiceCall("Create", LoadConfiguration().MealService.Port),
	},
	Route{
		"GetAllUserMeals",
		"GET",
		"/meal/all",
		ServiceCall("GetAllUserMeals", LoadConfiguration().MealService.Port),
	},
	Route{
		"Suggest",
		"GET",
		"/meal/suggest",
		ServiceCall("Suggest", LoadConfiguration().MealService.Port),
	},
	Route{
		"Read",
		"GET",
		"/meal/{id}",
		ServiceCall("Read", LoadConfiguration().MealService.Port),
	},
	Route{
		"Update",
		"POST",
		"/meal/{id}",
		ServiceCall("Update", LoadConfiguration().MealService.Port),
	},
	Route{
		"Delete",
		"DELETE",
		"/meal/{id}",
		ServiceCall("Delete", LoadConfiguration().MealService.Port),
	},
	Route{
		"Create",
		"POST",
		"/tag",
		ServiceCall("Create", LoadConfiguration().TagService.Port),
	},
	Route{
		"GetAllUserTags",
		"GET",
		"/tag/all",
		ServiceCall("GetAllUserTags", LoadConfiguration().TagService.Port),
	},
	Route{
		"Read",
		"GET",
		"/tag/{id}",
		ServiceCall("Read", LoadConfiguration().TagService.Port),
	},
	Route{
		"Update",
		"POST",
		"/tag/{id}",
		ServiceCall("Update", LoadConfiguration().TagService.Port),
	},
	Route{
		"Delete",
		"DELETE",
		"/tag/{id}",
		ServiceCall("Delete", LoadConfiguration().TagService.Port),
	},
	Route{
		"Create",
		"POST",
		"/calendar",
		ServiceCall("Create", LoadConfiguration().CalendarService.Port),
	},
	Route{
		"GetAllUserCalendars",
		"GET",
		"/calendar/all",
		ServiceCall("GetAllUserCalendars", LoadConfiguration().CalendarService.Port),
	},
	Route{
		"Read",
		"GET",
		"/calendar/{id}",
		ServiceCall("Read", LoadConfiguration().CalendarService.Port),
	},
	Route{
		"Update",
		"POST",
		"/calendar/{id}",
		ServiceCall("Update", LoadConfiguration().CalendarService.Port),
	},
	Route{
		"Delete",
		"DELETE",
		"/calendar/{id}",
		ServiceCall("Delete", LoadConfiguration().CalendarService.Port),
	},
}
