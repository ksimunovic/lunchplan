db=db.getSiblingDB("CalendarService");
db.getCollection("calendar").insert({ 
    "_id" : ObjectId("5ae5e06f3cd0501264f6d380"), 
    "date" : ISODate("2018-05-25T00:00:00.000+0000"), 
    "meal" : {
        "_id" : ObjectId("5ae4a10c239c97ba44011744"), 
        "title" : "WHAaat?2", 
        "description" : "Izmjenjeni content meal-a", 
        "profile" : {
            "_id" : ObjectId("5adcc4773cd05011d89b29a0"), 
            "firstname" : "Pero", 
            "lastname" : "Slav", 
            "servedby" : "DESKTOP-3J11AHK 169.254.17.84"
        }, 
        "timestamp" : NumberInt(1524932876), 
        "servedby" : "DESKTOP-3J11AHK 169.254.17.84"
    }
});
db.getCollection("calendar").insert({
    "_id" : ObjectId("5aeda76e3cd0501738e80b1e"), 
    "date" : ISODate("2018-05-28T00:00:00.000+0000"), 
    "meal" : {
        "_id" : ObjectId("5ae4b3323cd05032c8153057"), 
        "title" : "Saft i tijesto", 
        "description" : "Svinjsko mljeveno meso i tjestenina", 
        "profile" : {
            "_id" : ObjectId("5adcc4773cd05011d89b29a0"), 
            "firstname" : "Pero", 
            "lastname" : "Slav", 
            "servedby" : "DESKTOP-3J11AHK 169.254.17.84"
        }, 
        "timestamp" : NumberInt(1524937522), 
        "servedby" : "DESKTOP-3J11AHK 169.254.17.84"
    }
});
db.getCollection("calendar").insert({
    "_id" : ObjectId("5aedcf9d3cd0502ff4921d20"), 
    "date" : ISODate("2018-05-28T00:00:00.000+0000"), 
    "meal" : {
        "_id" : ObjectId("5ae4b3323cd05032c8153057"), 
        "title" : "Saft i tijesto", 
        "description" : "Svinjsko mljeveno meso i tjestenina", 
        "profile" : {
            "_id" : ObjectId("5adcc4773cd05011d89b29a0"), 
            "firstname" : "Pero", 
            "lastname" : "Slav", 
            "servedby" : "DESKTOP-3J11AHK 169.254.17.84"
        }, 
        "timestamp" : NumberInt(1524937522), 
        "servedby" : "DESKTOP-3J11AHK 169.254.17.84"
    }
});
db.getCollection("calendar").insert({
    "_id" : ObjectId("5af879564c751e15f60f7081"), 
    "date" : ISODate("2018-05-28T00:00:00.000+0000"), 
    "meal" : {
        "_id" : ObjectId("5ae4b3323cd05032c8153057"), 
        "title" : "Saft i tijesto", 
        "description" : "Svinjsko mljeveno meso i tjestenina", 
        "profile" : {
            "_id" : ObjectId("5adcc4773cd05011d89b29a0"), 
            "firstname" : "Pero", 
            "lastname" : "Slav", 
            "servedby" : "DESKTOP-3J11AHK 169.254.17.84"
        }, 
        "timestamp" : NumberInt(1524937522), 
        "servedby" : "Dev-Macbook.local 192.168.1.188"
    }
});
