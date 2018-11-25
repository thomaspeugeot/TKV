var fra = L.map('fra').setView([47, 0], 5);
var hti = L.map('hti').setView([18, -72], 5);

var mapOfMaps = new Map();

mapOfMaps.set( "fra", fra);
mapOfMaps.set( "hti", hti);

L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
	maxZoom: 18,
	attribution: 
		'A propos 10 000 territories; <a href="https://10ktblog.wordpress.com/a-propos/">10 000</a> contributors, ' +
		'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, ' +
		'<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
		'Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
	id: 'mapbox.streets'
}).addTo(fra);

var hostname
var protocol
var port
var targetService

var oReq 

var littleIcon = L.icon({
	iconUrl: '9pixels.png',
	iconSize:     [5, 5], // size of the icon
	iconAnchor:   [0, 0], // point of the icon which will correspond to marker's location
});

function onMapClick(e) {

	hostname = window.location.hostname
	protocol = window.location.protocol
	port = window.location.port
	targetService = protocol + "//"+ hostname + ":" + port + "/"

	messageToServer = { lat: e.latlng.lat , lng: e.latlng.lng, country: this._container.id}

	var messageToServerString = JSON.stringify( messageToServer );
	console.log( messageToServerString);	

	oReq = new XMLHttpRequest();
	// oReq.responseType = 'json';
	oReq.addEventListener("load", reqListener);
	oReq.open("POST", targetService +'translateLatLngInSourceCountryToLatLngInTargetCountry');
	oReq.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
	oReq.send( messageToServerString);				
};

function reqListener( evt) {
	
	var jsonResponse = JSON.parse( this.response)
	
	console.log('village translateLatLngInSourceCountryToLatLngInTargetCountry answer', 
		jsonResponse.X, jsonResponse.Y);

	lat = parseFloat(jsonResponse.LatClosest);
	lng = parseFloat(jsonResponse.LngClosest);

	latTarget = parseFloat(jsonResponse.LatTarget);
	lngTarget = parseFloat(jsonResponse.LngTarget);

	message = "Territory X="+ 
		Math.floor(100*jsonResponse.X)+" Y="+
		Math.floor(100*jsonResponse.Y);
	

		
	L.marker([lat, lng]).addTo( mapOfMaps.get( jsonResponse.Source))
		.bindPopup( message).openPopup();
		
	L.marker([latTarget, lngTarget]).addTo( mapOfMaps.get( jsonResponse.Target))
		.bindPopup( message).openPopup();

	for (var i = 0; i < jsonResponse.SourceBorderPoints[0].length; i++) {

		lng = parseFloat(jsonResponse.SourceBorderPoints[0][i][0]);
		lat = parseFloat(jsonResponse.SourceBorderPoints[0][i][1]);

		marker = new L.marker([lat,lng], {icon: littleIcon, opacity: 0.3} )
			.addTo( mapOfMaps.get( jsonResponse.Source));
	}

	for (var i = 0; i < jsonResponse.TargetBorderPoints[0].length; i++) {

		lng = parseFloat(jsonResponse.TargetBorderPoints[0][i][0]);
		lat = parseFloat(jsonResponse.TargetBorderPoints[0][i][1]);

		marker = new L.marker([lat,lng], {icon: littleIcon, opacity: 0.3})
			.addTo( mapOfMaps.get( jsonResponse.Target));
	}

	// reset zoom & location on target map 
	mapOfMaps.get( jsonResponse.Target).setView( [latTarget, lngTarget], 
		mapOfMaps.get( jsonResponse.Source).getZoom());
};


fra.on('click', onMapClick);

hti.on('click', onMapClick);

L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
	maxZoom: 18,
	attribution: 
		'A propos 10 000 territories; <a href="https://10ktblog.wordpress.com/a-propos/">10 000</a> contributors, ' +
		'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, ' +
		'<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
		'Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
	id: 'mapbox.streets'
}).addTo(hti);
