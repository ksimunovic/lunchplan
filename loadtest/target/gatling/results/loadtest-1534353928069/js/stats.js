var stats = {
    type: "GROUP",
name: "Global Information",
path: "",
pathFormatted: "group_missing-name-b06d1",
stats: {
    "name": "Global Information",
    "numberOfRequests": {
        "total": "1226",
        "ok": "1226",
        "ko": "0"
    },
    "minResponseTime": {
        "total": "100",
        "ok": "100",
        "ko": "-"
    },
    "maxResponseTime": {
        "total": "10748",
        "ok": "10748",
        "ko": "-"
    },
    "meanResponseTime": {
        "total": "4100",
        "ok": "4100",
        "ko": "-"
    },
    "standardDeviation": {
        "total": "2361",
        "ok": "2361",
        "ko": "-"
    },
    "percentiles1": {
        "total": "3989",
        "ok": "3989",
        "ko": "-"
    },
    "percentiles2": {
        "total": "6007",
        "ok": "6007",
        "ko": "-"
    },
    "percentiles3": {
        "total": "7903",
        "ok": "7903",
        "ko": "-"
    },
    "percentiles4": {
        "total": "9140",
        "ok": "9140",
        "ko": "-"
    },
    "group1": {
        "name": "t < 800 ms",
        "count": 118,
        "percentage": 10
    },
    "group2": {
        "name": "800 ms < t < 1200 ms",
        "count": 58,
        "percentage": 5
    },
    "group3": {
        "name": "t > 1200 ms",
        "count": 1050,
        "percentage": 86
    },
    "group4": {
        "name": "failed",
        "count": 0,
        "percentage": 0
    },
    "meanNumberOfRequestsPerSecond": {
        "total": "10.255",
        "ok": "10.255",
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
        "total": "1226",
        "ok": "1226",
        "ko": "0"
    },
    "minResponseTime": {
        "total": "100",
        "ok": "100",
        "ko": "-"
    },
    "maxResponseTime": {
        "total": "10748",
        "ok": "10748",
        "ko": "-"
    },
    "meanResponseTime": {
        "total": "4100",
        "ok": "4100",
        "ko": "-"
    },
    "standardDeviation": {
        "total": "2361",
        "ok": "2361",
        "ko": "-"
    },
    "percentiles1": {
        "total": "3989",
        "ok": "3989",
        "ko": "-"
    },
    "percentiles2": {
        "total": "6007",
        "ok": "6007",
        "ko": "-"
    },
    "percentiles3": {
        "total": "7903",
        "ok": "7903",
        "ko": "-"
    },
    "percentiles4": {
        "total": "9140",
        "ok": "9140",
        "ko": "-"
    },
    "group1": {
        "name": "t < 800 ms",
        "count": 118,
        "percentage": 10
    },
    "group2": {
        "name": "800 ms < t < 1200 ms",
        "count": 58,
        "percentage": 5
    },
    "group3": {
        "name": "t > 1200 ms",
        "count": 1050,
        "percentage": 86
    },
    "group4": {
        "name": "failed",
        "count": 0,
        "percentage": 0
    },
    "meanNumberOfRequestsPerSecond": {
        "total": "10.255",
        "ok": "10.255",
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
