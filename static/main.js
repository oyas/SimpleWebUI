
function clear_connection_status_text(){
	$("#connection_status").html("");
}

window.onload = function () {
	// add class for bootstrap
	$('table').addClass('table table-stlipe')
	$('pre').addClass('pre')

	var conn;
	
	var vue_data = {
		message: 'message',
		timer: '',
		d: {
			e: 'ok'
		}
	};

	var el = document.getElementById('main_container');
	var compiled = Vue.compile('<div>' + el.innerHTML + '</div>');
	var methods = {
		submit: function(){
			//var raw_form_data = new FormData(document.getElementById('main_form').cloneNode(true));
			//var form_data = Object.assign(...Array.from(raw_form_data, ([x,y]) => ({[x]:y})));
			var form_data = $('#main_form').serializeArray();
			conn.send( JSON.stringify(form_data) );
			console.log('sent data:' + JSON.stringify(form_data));
		},
	};

	var vm = new Vue({
		el: '#main_container',
		data: vue_data,
		methods: methods,
	});
	var root_keys = Object.keys(vue_data).toString();

	var make_connection = function(){
		console.log("Connecting...");
		conn = new WebSocket("ws://" + document.location.host + "/ws");
		conn.onopen = function (evt) {
			$("#connection_status").html("<span style=\"color: green;\">Connected.</span>");
			setTimeout(clear_connection_status_text, 2000);
			console.log("Connected.");
		};
		conn.onclose = function (evt) {
			var item = document.createElement("div");
			item.innerHTML = "<b>Connection closed.</b>";
			$("#connection_status").html("<span style=\"color: red;\">Connection closed.</span>");
			setTimeout(clear_connection_status_text, 5000);
			console.log("Connection closed.");
			setTimeout(make_connection, 5000);    // 再接続
		};
		conn.onmessage = function (evt) {
			var messages = evt.data.split('\n');
			data = "";
			while( data == "" ) data = messages.pop();
			try {
				data = JSON.parse(data);
				Object.assign(vue_data, data);
				keys = Object.keys(vue_data).toString();
				if( keys != root_keys ){
					vm.$destroy();
					vm = new Vue({
						data: vue_data,
						methods: methods,
						render: compiled.render,
						staticRenderFns: compiled.staticRenderFns
					});
					vm.$mount('#main_container');
					root_keys = keys;
					console.log('rebuilded');
				}
				vm.$forceUpdate();
			} catch(e) {
				console.log(e);
			}
		};
	};

	if (window["WebSocket"]) {
		make_connection();
	} else {
		$("#connection_status").html("Your browser does not support WebSockets.");
	}
};
