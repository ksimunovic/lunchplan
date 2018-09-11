var stats = {
    type: "GROUP",
name: "Global Information",
path: "",
pathFormatted: "group_missing-name-b06d1",
stats: {
    "name": "Global Information",
    "numberOfRequests": {
        "total": "1207",
        "ok": "1207",
        "ko": "0"
    },
    "minResponseTime": {
        "total": "72",
        "ok": "72",
        "ko": "-"
    },
    "maxResponseTime": {
        "total": "11190",
        "ok": "11190",
        "ko": "-"
    },
    "meanResponseTime": {
        "total": "4180",
        "ok": "4180",
        "ko": "-"
    },
    "standardDeviation": {
        "total": "2689",
        "ok": "2689",
        "ko": "-"
    },
    "percentiles1": {
        "total": "3792",
        "ok": "3792",
        "ko": "-"
    },
    "percentiles2": {
        "total": "6373",
        "ok": "6373",
        "ko": "-"
    },
    "percentiles3": {
        "total": "8693",
        "ok": "8693",
        "ko": "-"
    },
    "percentiles4": {
        "total": "10068",
        "ok": "10068",
        "ko": "-"
    },
    "group1": {
        "name": "t < 800 ms",
        "count": 152,
        "percentage": 13
    },
    "group2": {
        "name": "800 ms < t < 1200 ms",
        "count": 44,
        "percentage": 4
    },
    "group3": {
        "name": "t > 1200 ms",
        "count": 1011,
        "percentage": 84
    },
    "group4": {
        "name": "failed",
        "count": 0,
        "percentage": 0
    },
    "meanNumberOfRequestsPerSecond": {
        "total": "10.035",
        "ok": "10.035",
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
        "total": "1207",
        "ok": "1207",
        "ko": "0"
    },
    "minResponseTime": {
        "total": "72",
        "ok": "72",
        "ko": "-"
    },
    "maxResponseTime": {
        "total": "11190",
        "ok": "11190",
        "ko": "-"
    },
    "meanResponseTime": {
        "total": "4180",
        "ok": "4180",
        "ko": "-"
    },
    "standardDeviation": {
        "total": "2689",
        "ok": "2689",
        "ko": "-"
    },
    "percentiles1": {
        "total": "3792",
        "ok": "3792",
        "ko": "-"
    },
    "percentiles2": {
        "total": "6373",
        "ok": "6373",
        "ko": "-"
    },
    "percentiles3": {
        "total": "8693",
        "ok": "8693",
        "ko": "-"
    },
    "percentiles4": {
        "total": "10068",
        "ok": "10068",
        "ko": "-"
    },
    "group1": {
        "name": "t < 800 ms",
        "count": 152,
        "percentage": 13
    },
    "group2": {
        "name": "800 ms < t < 1200 ms",
        "count": 44,
        "percentage": 4
    },
    "group3": {
        "name": "t > 1200 ms",
        "count": 1011,
        "percentage": 84
    },
    "group4": {
        "name": "failed",
        "count": 0,
        "percentage": 0
    },
    "meanNumberOfRequestsPerSecond": {
        "total": "10.035",
        "ok": "10.035",
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
