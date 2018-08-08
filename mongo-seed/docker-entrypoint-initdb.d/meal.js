db=db.getSiblingDB("MealService");
db.getCollection("meal").insert({ 
    "_id" : ObjectId("5b0850b74c751edb1c907f98"), 
    "title" : "prvi33", 
    "description" : "prvi", 
    "profile" : {
        "_id" : ObjectId("5ac8ac143cd050255c1a6eca"), 
        "firstname" : "Karlo", 
        "lastname" : "Šimunović"
    }, 
    "timestamp" : NumberInt(1527271607), 
    "tags" : [
        {
            "_id" : ObjectId("5b084a454c751eda2a6c6cad"), 
            "name" : "noviTag", 
            "profile" : {
                "_id" : ObjectId("5ac8ac143cd050255c1a6eca"), 
                "firstname" : "Karlo", 
                "lastname" : "Šimunović"
            }
        }, 
        {
            "_id" : ObjectId("5b0850b74c751eda2a6c6caf"), 
            "name" : "prvi", 
            "profile" : {
                "_id" : ObjectId("5ac8ac143cd050255c1a6eca"), 
                "firstname" : "Karlo", 
                "lastname" : "Šimunović"
            }
        }
    ]
});
db.getCollection("meal").insert({
    "_id" : ObjectId("5b0851044c751edb1c907f9a"), 
    "title" : "Demo", 
    "description" : "demo", 
    "profile" : {
        "_id" : ObjectId("5adcc4773cd05011d89b29a0"), 
        "firstname" : "Pero", 
        "lastname" : "Slav"
    }, 
    "timestamp" : NumberInt(1527271684), 
    "tags" : [
        {
            "_id" : ObjectId("5afc83664c751e47ed695250"), 
            "name" : "naslov taga", 
            "profile" : {
                "_id" : ObjectId("5adcc4773cd05011d89b29a0"), 
                "firstname" : "Pero", 
                "lastname" : "Slav"
            }
        }
    ]
});
