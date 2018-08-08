db=db.getSiblingDB("UserService");
db.getCollection("profile").insert({ 
    "_id" : ObjectId("5ac8ac143cd050255c1a6eca"), 
    "firstname" : "Karlo", 
    "lastname" : "Šimunović"
});
db.getCollection("profile").insert({
    "_id" : ObjectId("5adcc4773cd05011d89b29a0"), 
    "firstname" : "Pero", 
    "lastname" : "Slav", 
    "servedby" : ""
});
