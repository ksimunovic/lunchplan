db.createUser({user: "root", pwd: "root", roles: [{
            role : "__queryableBackup", 
            db : "admin"
        }, 
        {
            role : "__system", 
            db : "admin"
        }, 
        {
            role : "backup", 
            db : "admin"
        }, 
        {
            role : "clusterAdmin", 
            db : "admin"
        }, 
        {
            role : "clusterManager", 
            db : "admin"
        }, 
        {
            role : "clusterMonitor", 
            db : "admin"
        }, 
        {
            role : "dbAdmin", 
            db : "admin"
        }, 
        {
            role : "dbAdminAnyDatabase", 
            db : "admin"
        }, 
        {
            role : "dbOwner", 
            db : "admin"
        }, 
        {
            role : "enableSharding", 
            db : "admin"
        }, 
        {
            role : "hostManager", 
            db : "admin"
        }, 
        {
            role : "read", 
            db : "admin"
        }, 
        {
            role : "readAnyDatabase", 
            db : "admin"
        }, 
        {
            role : "readWrite", 
            db : "admin"
        }, 
        {
            role : "readWriteAnyDatabase", 
            db : "admin"
        }, 
        {
            role : "restore", 
            db : "admin"
        }, 
        {
            role : "root", 
            db : "admin"
        }, 
        {
            role : "userAdmin", 
            db : "admin"
        }, 
        {
            role : "userAdminAnyDatabase", 
            db : "admin"
        }]});