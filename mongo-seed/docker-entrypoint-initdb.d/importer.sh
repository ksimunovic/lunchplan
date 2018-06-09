#!/bin/bash
mongoimport --host localhost --db UserService --collection account --type json --file docker-entrypoint-initdb.d/account.json 
mongoimport --host localhost --db UserService --collection profile --type json --file docker-entrypoint-initdb.d/profile.json 
mongoimport --host localhost --db UserService --collection session --type json --file docker-entrypoint-initdb.d/session.json 
mongoimport --host localhost --db MealService --collection meal --type json --file docker-entrypoint-initdb.d/meal.json 
mongoimport --host localhost --db TagService --collection tag --type json --file docker-entrypoint-initdb.d/tag.json 
mongoimport --host localhost --db CalendarService --collection calendar --type json --file docker-entrypoint-initdb.d/calendar.json 