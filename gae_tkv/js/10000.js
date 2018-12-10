var mapOfMapViews = new Map();

mapOfMapViews.set( 'fra', [47, 0]);
mapOfMapViews.set( 'hti', [18, -72]);
mapOfMapViews.set( 'usa', [39, -100]);

// dynamicaly allocate topMap and bottomMap map
function TwoWayMap(map){
	this.map = map;
	this.reverseMap = {};
	for(var key in map){
	   var value = map[key];
	   this.reverseMap[value] = key;   
	}
 }
 TwoWayMap.prototype.set = function(key, value){ this.map[key] = value; this.reverseMap[value] = key;};
 TwoWayMap.prototype.get = function(key){ return this.map[key]; };
 TwoWayMap.prototype.revGet = function(key){ return this.reverseMap[key]; };

var mapOfMapNames = new TwoWayMap( 
	{
		'topMap': 'usa',
		'bottomMap': 'hti'
	}
);

var topMapCenter = mapOfMapViews.get( mapOfMapNames.get( 'topMap'));
var bottomMapCenter = mapOfMapViews.get( mapOfMapNames.get( 'bottomMap'));

var topMap = L.map('topMap').setView( topMapCenter, 4);
var bottomMap = L.map('bottomMap').setView( bottomMapCenter, 4);

var mapOfMaps = new Map();

mapOfMaps.set( 'topMap', topMap);
mapOfMaps.set( 'bottomMap', bottomMap);

// return the country on the other side
function otherSideCountry( country) {
	if (mapOfMapNames.get( 'topMap') == country) {
		return mapOfMapNames.get( 'bottomMap')
	} else {
		return mapOfMapNames.get('topMap')
	}
}

L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
	maxZoom: 18,
	id: 'mapbox.streets'
}).addTo(topMap);

L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
	maxZoom: 18,
	attribution: 
		'<a href="https://10ktblog.wordpress.com/a-propos/">10 000</a> ' +
		'<a href="https://www.openstreetmap.org/">OpenStreetMap</a>  ' +
		'<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
		'<a href="https://www.mapbox.com/">Mapbox</a>',
	id: 'mapbox.streets'
}).addTo(bottomMap);

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

L.Control.SwitchMap = L.Control.extend({
    onAdd: function(map) {
        var img = L.DomUtil.create('img');

		img.type = "button";
		img.value = "fra <> usa";
		img.style.backgroundSize = "30px 30px";
		img.style.width = '30px';
		img.style.height = '30px';
		img.style.backgroundColor = 'white';

		img.onclick = function() {
			console.log('buttonClicked');
			currentTopMap = mapOfMapNames.get('topMap')
			if ( 'fra' == currentTopMap ) {
				mapOfMapNames.set('topMap', 'usa');
 			} else {
				mapOfMapNames.set('topMap', 'fra');
			}
			topMapCenter = mapOfMapViews.get( mapOfMapNames.get( 'topMap'));
			topMap.setView( topMapCenter, 4);
			topMap.removeLayer( Markers);	
		}
		return img;
    },
});

L.control.switchMap = function(opts) {
    return new L.Control.SwitchMap(opts);
}

L.control.switchMap({ position: 'bottomleft' }).addTo(topMap);

function onMapClick(e) {

	hostname = window.location.hostname
	protocol = window.location.protocol
	port = window.location.port
	targetService = protocol + "//"+ hostname + ":" + port + "/"

    var sideOfMap = this._container.id
    var sourceCountry = mapOfMapNames.get( sideOfMap)

	var targetCountry = otherSideCountry( sourceCountry)

	messageToServer = { lat: e.latlng.lat , lng: e.latlng.lng, 
		sourceCountry: sourceCountry, targetCountry: targetCountry }

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
	
	L.marker([lat, lng]).addTo( mapOfMaps.get( mapOfMapNames.revGet( jsonResponse.Source)))
		.bindPopup( message).openPopup();
		
	L.marker([latTarget, lngTarget]).addTo( mapOfMaps.get( mapOfMapNames.revGet(jsonResponse.Target)))
		.bindPopup( message).openPopup();

	for (var i = 0; i < jsonResponse.SourceBorderPoints[0].length; i++) {

		lng = parseFloat(jsonResponse.SourceBorderPoints[0][i][0]);
		lat = parseFloat(jsonResponse.SourceBorderPoints[0][i][1]);

		marker = new L.marker([lat,lng], {icon: littleIcon, opacity: 0.3} )
			.addTo( mapOfMaps.get( mapOfMapNames.revGet( jsonResponse.Source)));
	}

	for (var i = 0; i < jsonResponse.TargetBorderPoints[0].length; i++) {

		lng = parseFloat(jsonResponse.TargetBorderPoints[0][i][0]);
		lat = parseFloat(jsonResponse.TargetBorderPoints[0][i][1]);

		marker = new L.marker([lat,lng], {icon: littleIcon, opacity: 0.3})
			.addTo( mapOfMaps.get( mapOfMapNames.revGet( jsonResponse.Target)));
	}

	// reset zoom & location on target map 
	mapOfMaps.get( mapOfMapNames.revGet(jsonResponse.Target)).setView( [latTarget, lngTarget], 
		mapOfMaps.get( mapOfMapNames.revGet( jsonResponse.Source)).getZoom());
};


topMap.on('click', onMapClick);

bottomMap.on('click', onMapClick);
