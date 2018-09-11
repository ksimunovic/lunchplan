var stats = {
    type: "GROUP",
name: "Global Information",
path: "",
pathFormatted: "group_missing-name-b06d1",
stats: {
    "name": "Global Information",
    "numberOfRequests": {
        "total": "480",
        "ok": "480",
        "ko": "0"
    },
    "minResponseTime": {
        "total": "69",
        "ok": "69",
        "ko": "-"
    },
    "maxResponseTime": {
        "total": "1572",
        "ok": "1572",
        "ko": "-"
    },
    "meanResponseTime": {
        "total": "259",
        "ok": "259",
        "ko": "-"
    },
    "standardDeviation": {
        "total": "179",
        "ok": "179",
        "ko": "-"
    },
    "percentiles1": {
        "total": "219",
        "ok": "219",
        "ko": "-"
    },
    "percentiles2": {
        "total": "353",
        "ok": "353",
        "ko": "-"
    },
    "percentiles3": {
        "total": "570",
        "ok": "570",
        "ko": "-"
    },
    "percentiles4": {
        "total": "749",
        "ok": "749",
        "ko": "-"
    },
    "group1": {
        "name": "t < 800 ms",
        "count": 475,
        "percentage": 99
    },
    "group2": {
        "name": "800 ms < t < 1200 ms",
        "count": 3,
        "percentage": 1
    },
    "group3": {
        "name": "t > 1200 ms",
        "count": 2,
        "percentage": 0
    },
    "group4": {
        "name": "failed",
        "count": 0,
        "percentage": 0
    },
    "meanNumberOfRequestsPerSecond": {
        "total": "4.173",
        "ok": "4.173",
        "ko": "-"
    }
},
contents: {
"req_apilogin-23d40": {
        type: "REQUEST",
        name: "ApiLogin",
path: "ApiLogin",
pathFormatted: "req_apilogin-23d40",
stats: {
    "name": "ApiLogin",
    "numberOfRequests": {
        "total": "480",
        "ok": "480",
        "ko": "0"
    },
    "minResponseTime": {
        "total": "69",
        "ok": "69",
        "ko": "-"
    },
    "maxResponseTime": {
        "total": "1572",
        "ok": "1572",
        "ko": "-"
    },
    "meanResponseTime": {
        "total": "259",
        "ok": "259",
        "ko": "-"
    },
    "standardDeviation": {
        "total": "179",
        "ok": "179",
        "ko": "-"
    },
    "percentiles1": {
        "total": "219",
        "ok": "219",
        "ko": "-"
    },
    "percentiles2": {
        "total": "353",
        "ok": "353",
        "ko": "-"
    },
    "percentiles3": {
        "total": "570",
        "ok": "570",
        "ko": "-"
    },
    "percentiles4": {
        "total": "749",
        "ok": "749",
        "ko": "-"
    },
    "group1": {
        "name": "t < 800 ms",
        "count": 475,
        "percentage": 99
    },
    "group2": {
        "name": "800 ms < t < 1200 ms",
        "count": 3,
        "percentage": 1
    },
    "group3": {
        "name": "t > 1200 ms",
        "count": 2,
        "percentage": 0
    },
    "group4": {
        "name": "failed",
        "count": 0,
        "percentage": 0
    },
    "meanNumberOfRequestsPerSecond": {
        "total": "4.173",
        "ok": "4.173",
        "ko": "-"
    }
}
    }
}

}

function fillStats(stat){
    $("#numberOfRequests").append(stat.numberOfRequests.total);
    $("#numberOfRequestsOK").append(stat.numberOfRequests.ok);
    $("#numberOfRequestsKO").append(stat.numberOfRequests.ko);

    $("#minResponseTime").append(stat.minResponseTime.total);
    $("#minResponseTimeOK").append(stat.minResponseTime.ok);
    $("#minResponseTimeKO").append(stat.minResponseTime.ko);

    $("#maxResponseTime").append(stat.maxResponseTime.total);
    $("#maxResponseTimeOK").append(stat.maxResponseTime.ok);
    $("#maxResponseTimeKO").append(stat.maxResponseTime.ko);

    $("#meanResponseTime").append(stat.meanResponseTime.total);
    $("#meanResponseTimeOK").append(stat.meanResponseTime.ok);
    $("#meanResponseTimeKO").append(stat.meanResponseTime.ko);

    $("#standardDeviation").append(stat.standardDeviation.total);
    $("#standardDeviationOK").append(stat.standardDeviation.ok);
    $("#standardDeviationKO").append(stat.standardDeviation.ko);

    $("#percentiles1").append(stat.percentiles1.total);
    $("#percentiles1OK").append(stat.percentiles1.ok);
    $("#percentiles1KO").append(stat.percentiles1.ko);

    $("#percentiles2").append(stat.percentiles2.total);
    $("#percentiles2OK").append(stat.percentiles2.ok);
    $("#percentiles2KO").append(stat.percentiles2.ko);

    $("#percentiles3").append(stat.percentiles3.total);
    $("#percentiles3OK").append(stat.percentiles3.ok);
    $("#percentiles3KO").append(stat.percentiles3.ko);

    $("#percentiles4").append(stat.percentiles4.total);
    $("#percentiles4OK").append(stat.percentiles4.ok);
    $("#percentiles4KO").append(stat.percentiles4.ko);

    $("#meanNumberOfRequestsPerSecond").append(stat.meanNumberOfRequestsPerSecond.total);
    $("#meanNumberOfRequestsPerSecondOK").append(stat.meanNumberOfRequestsPerSecond.ok);
    $("#meanNumberOfRequestsPerSecondKO").append(stat.meanNumberOfRequestsPerSecond.ko);
}
