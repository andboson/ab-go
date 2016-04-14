var chart1;
var chart2;
var ws;

$(function() {
	var x = new Date();

	Highcharts.setOptions({
		global: {
			timezoneOffset: x.getTimezoneOffset()
		}
	})

		chart1 = new Highcharts.StockChart({
			chart: {
				renderTo: 'container1',
				zoomType: 'x'
			},
			title: {
				text: 'Response time'
			},
			yAxis: {
				title: {
					text: 'Milliseconds'
				}
			},
			scrollbar: {
				enabled: false
			},
			rangeSelector: {
				buttons: [{
					type: 'second',
					count: 5,
					text: '5s'
				}, {
					type: 'second',
					count: 30,
					text: '30s'
				}, {
					type: 'minute',
					count: 1,
					text: '1m'
				}, {
					type: 'all',
					text: 'All'
				}],
				selected: 3
			},
			series: [{
				name: "Average response time",
				data: null,
				color:'blue',
				type: 'line',
				tooltip: {
					valueSuffix: 'ms'
				}
			},{
				name: "Max response time",
				data: null,
				color:'red',
				type: 'line',
				tooltip: {
					valueSuffix: 'ms'
				}
             },
             {
                name: "Min response time",
                data: null,
                color:'green',
                type: 'line',
                tooltip: {
                    valueSuffix: 'ms'
                }
             }]
		});

		chart2 = new Highcharts.StockChart({
        			chart: {
        				renderTo: 'container2',
        				zoomType: 'x'
        			},
        			title: {
        				text: 'Requests per second'
        			},
        			yAxis: {
        				title: {
        					text: 'count'
        				}
        			},
        			scrollbar: {
        				enabled: false
        			},
        			rangeSelector: {
        				buttons: [{
        					type: 'second',
        					count: 5,
        					text: '5s'
        				}, {
        					type: 'second',
        					count: 30,
        					text: '30s'
        				}, {
        					type: 'minute',
        					count: 1,
        					text: '1m'
        				}, {
        					type: 'all',
        					text: 'All'
        				}],
        				selected: 3
        			},
        			series: [{
        				name: "Rps",
        				data: null,
        				type: 'line',
        				tooltip: {
        					valueSuffix: ''
        				}
        			}]
        		});

                //reinit socket, because someone close it
        		makeSocket();
        		intervalID = window.setInterval(function(){
        		    makeSocket();
        		}, 60 * 1000)
})


function wsurl() {
    var l = window.location;
    return ((l.protocol === "https:") ? "wss://" : "ws://") + l.hostname + (((l.port != 80) && (l.port != 443)) ? ":" + l.port : "") + "/ws";
}

function makeSocket(){
	ws = new WebSocket(wsurl());
	ws.onopen = function () {
		ws.onmessage = function (evt) {
			///console.log(evt.data)
			var data = JSON.parse(evt.data);
			chart1.series[0].addPoint([data.Ts, parseFloat(data.Avg)], true);
			chart1.series[1].addPoint([data.Ts, parseFloat(data.Max)], true);
			chart1.series[2].addPoint([data.Ts, parseFloat(data.Min)], true);
			chart2.series[0].addPoint([data.Ts, parseInt(data.Rps)], true);
			$('#container3').html(data.LastResult)
		}
	};
}