db=db.getSiblingDB("UserService");
db.getCollection("account").insert({ 
    "_id" : ObjectId("5ac8ac14b17075f123dab560"), 
    "profile" : {
        "_id" : ObjectId("5ac8ac143cd050255c1a6eca"), 
        "firstname" : "Karlo", 
        "lastname" : "Šimunović"
    }, 
    "email" : "probniEmail", 
    "password" : "$2a$10$/ob4brtUOX0kPfnohVAcf.O4o9puVxFj7H/eh9L.wku14J8X7WnpW"
});
db.getCollection("account").insert({
    "_id" : ObjectId("5adcc478239c97ba4400d517"), 
    "profile" : {
        "_id" : ObjectId("5adcc4773cd05011d89b29a0"), 
        "firstname" : "Pero", 
        "lastname" : "Slav", 
        "servedby" : ""
    }, 
    "email" : "peroslav@gmail.com", 
    "password" : "$2a$10$bNBAdJNt1RoPmhfluacF3OuX.Xbp5dC4NEvMJ8kLK/wAF4.aFwR1m"
});
db.getCollection("account").insert({
    "_id" : ObjectId("5ae21687239c97ba4400fbbc"), 
    "profile" : {
        "_id" : ObjectId("5ae216873cd050284811e0c5"), 
        "firstname" : "Pero", 
        "lastname" : "Slav", 
        "servedby" : ""
    }, 
    "email" : "peroslav@gmail.com", 
    "password" : "$2a$10$tIOEhDWgQpGqaohbjXj05uwHFFgvVS/Gw5.UPjK.Lvt0GCeUFBj3i"
});
db.getCollection("account").insert({
    "_id" : ObjectId("5ae4bdbe239c97ba440123d4"), 
    "profile" : {
        "_id" : ObjectId("5ae4bdbe3cd0501f8866aa6d"), 
        "firstname" : "Pero", 
        "lastname" : "Slav", 
        "servedby" : ""
    }, 
    "email" : "peroslav@gmail.com", 
    "password" : "$2a$10$RJUAtGNBvY1Z00MAxAFg.OTj4LoMYyOpUYRc7dKtIMw4jx9dwU29m"
});
db.getCollection("account").insert({
    "_id" : ObjectId("5aed709597de6d2a6b7ca5a9"), 
    "profile" : {
        "_id" : ObjectId("5aed70953cd0502108044ab7"), 
        "firstname" : "Pero", 
        "lastname" : "Slav", 
        "servedby" : ""
    }, 
    "email" : "peroslav@gmail.com", 
    "password" : "$2a$10$qvLHvo1p4xWPgqrhDPoTRea2JQeyVJol8SXAadCXmjoSrequdyP4y"
});
db.getCollection("account").insert({
    "_id" : ObjectId("5aed865997de6d2a6b7caf56"), 
    "profile" : {
        "_id" : ObjectId("5aed86593cd0502108044abb"), 
        "firstname" : "Karlo", 
        "lastname" : "Šimunović", 
        "servedby" : ""
    }, 
    "email" : "probniEmail", 
    "password" : "$2a$10$oxkJGzmWA.zJQLxQFAmjpO0LyX48U7.Lrn3y.vRgw66IxQsMIm3Jm"
});