var chart1;
var chart2;

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
				type: 'area',
				tooltip: {
					valueSuffix: 'ms'
				}
			},{
				name: "Max response time",
				data: null,
				type: 'area',
				tooltip: {
					valueSuffix: 'ms'
				}
             },
             {
                name: "Min response time",
                data: null,
                type: 'area',
                tooltip: {
                    valueSuffix: 'ms'
                }
             }]
		});

	function wsurl() {
		var l = window.location;
		return ((l.protocol === "https:") ? "wss://" : "ws://") + l.hostname + (((l.port != 80) && (l.port != 443)) ? ":" + l.port : "") + "/ws";
	}

	ws = new WebSocket(wsurl());
	ws.onopen = function () {
		ws.onmessage = function (evt) {
			//console.log(evt.data)
			var data = JSON.parse(evt.data);
			chart1.series[0].addPoint([data.Ts, parseFloat(data.Avg)], true);
			chart1.series[1].addPoint([data.Ts, parseFloat(data.Max)], true);
			chart1.series[2].addPoint([data.Ts, parseFloat(data.Min)], true);
		}
	};
})
