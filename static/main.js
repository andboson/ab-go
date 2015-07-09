var chart1;
var chart2;

$(function() {
	var x = new Date();

	Highcharts.setOptions({
		global: {
			timezoneOffset: x.getTimezoneOffset()
		}
	})

	$.getJSON('/data?callback=?', function(data) {
		chart1 = new Highcharts.StockChart({
			chart: {
				renderTo: 'container1',
				zoomType: 'x'
			},
			title: {
				text: 'GC pauses'
			},
			yAxis: {
				title: {
					text: 'Nanoseconds'
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
				name: "GC pauses",
				data: data.GcPauses,
				type: 'area',
				tooltip: {
					valueSuffix: 'ns'
				}
			}]
		});

	});

	function wsurl() {
		var l = window.location;
		return ((l.protocol === "https:") ? "wss://" : "ws://") + l.hostname + (((l.port != 80) && (l.port != 443)) ? ":" + l.port : "") + "/ws";
	}

	ws = new WebSocket(wsurl());
	ws.onopen = function () {
		ws.onmessage = function (evt) {
			console.log(evt.data)
		//	var data = JSON.parse(evt.data);
		//	chart1.series[0].addPoint([data.Ts, data.GcPause], true);
		}
	};
})
