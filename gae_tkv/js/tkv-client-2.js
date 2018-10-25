var france = L.map('france').setView([47, 0], 5);

L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
	maxZoom: 18,
	attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, ' +
		'<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
		'Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
	id: 'mapbox.streets'
}).addTo(france);

var hostname
var protocol
var port
var targetService

var oReq 

function onFranceMapClick(e) {

	hostname = window.location.hostname
	protocol = window.location.protocol
	port = window.location.port
	targetService = protocol + "//"+ hostname + ":" + port + "/"


	var jsonLatLng = JSON.stringify( e.latlng);
	console.log( jsonLatLng);

	oReq = new XMLHttpRequest();
	oReq.responseType = 'json';
	oReq.addEventListener("load", reqListener);
	oReq.open("POST", targetService +'translateLatLngInSourceCountryToLatLngInTargetCountry');
	oReq.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
	oReq.send( jsonLatLng);				
};

function reqListener( evt) {
	
	var jsonResponse = this.response
	
	console.log('village translateLatLngInSourceCountryToLatLngInTargetCountry answer', 
		jsonResponse.X, jsonResponse.Y);

	lat = parseFloat(jsonResponse.LatClosest);
	lng = parseFloat(jsonResponse.LngClosest);

	latTarget = parseFloat(jsonResponse.LatTarget);
	lngTarget = parseFloat(jsonResponse.LngTarget);

	xSpead = parseFloat(jsonResponse.Xspread);
	ySpead = parseFloat(jsonResponse.Yspread);

	message = "Territory X="+ 
		Math.floor(100*jsonResponse.Xspread)+" Y="+
		Math.floor(100*jsonResponse.Yspread);
			
	L.marker([lat, lng]).addTo(france)
		.bindPopup( message).openPopup();
		
	L.marker([latTarget, lngTarget]).addTo(haiti)
		.bindPopup( message).openPopup();

};


france.on('click', onFranceMapClick);

var haiti = L.map('haiti').setView([18, -72], 5);

L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
	maxZoom: 18,
	attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, ' +
		'<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
		'Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
	id: 'mapbox.streets'
}).addTo(haiti);
