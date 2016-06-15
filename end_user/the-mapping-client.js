var map;
var geocoder;
var bounds = new google.maps.LatLngBounds();
var markersArray = [];

// the mapping elements
var parisBounds = new google.maps.LatLngBounds();
var newYorkBounds = new google.maps.LatLngBounds();
var haitiBounds = new google.maps.LatLngBounds();
var senegalBounds = new google.maps.LatLngBounds();

var parisPosition;
var newYorkPosition;
var haitiPosition;
var senegalPosition;



var origin1 = new google.maps.LatLng(55.930, -3.118);
var origin2 = 'Greenwich, England';
var destinationA = 'Stockholm, Sweden';
var destinationB = new google.maps.LatLng(50.087, 14.421);

var destinationIcon = 'https://chart.googleapis.com/chart?chst=d_map_pin_letter&chld=D|FF0000|000000';
var originIcon = 'https://chart.googleapis.com/chart?chst=d_map_pin_letter&chld=O|FFFF00|000000';

function initialize() {
  var opts = {
		  // centrer on paris
    center: new google.maps.LatLng(48.75, 2.3),
    zoom: 10,
    mapTypeControl: true,
    mapTypeControlOptions: {
      style: google.maps.MapTypeControlStyle.DROPDOWN_MENU
    },
    zoomControl: true,
    zoomControlOptions: {
      style: google.maps.ZoomControlStyle.SMALL
    }

  };
  map = new google.maps.Map(document.getElementById('map-canvas'), opts);
  geocoder = new google.maps.Geocoder();
  
  getThemArea();
}


function calculateDistances() {
  var service = new google.maps.DistanceMatrixService();
  service.getDistanceMatrix(
    {
      origins: [origin1, origin2],
      destinations: [destinationA, destinationB],
      travelMode: google.maps.TravelMode.DRIVING,
      unitSystem: google.maps.UnitSystem.METRIC,
      avoidHighways: false,
      avoidTolls: false
    }, callback);
}

function reqListener ( oReq) {

	 var parser=new DOMParser();
	  var xmlDoc=parser.parseFromString(oReq.responseText,"text/xml");
	  var tds = xmlDoc.getElementsByTagName("hello");

	     outputDiv.innerHTML += 'response response tds ' + tds  + '<br>'; 
	     outputDiv.innerHTML += 'response response tds ' + tds.toSource  + '<br>'; 

     outputDiv.innerHTML += 'response response ' + oReq  + '<br>'; 
     outputDiv.innerHTML += 'response text ' + oReq.responseText.string  + '<br>'; 
     outputDiv.innerHTML += 'response text ' + oReq.responseText.toSource  + '<br>'; 
     outputDiv.innerHTML += 'response type ' + oReq.responseType.toSource  + '<br>'; 
     outputDiv.innerHTML += 'reponse XMLHttpRequest complete' + '<br>'; 
    
//      for (var property in oReq) {
//      	 if ( oReq[property] != null ) {
//     	       outputDiv.innerHTML += property + ': ' + oReq[property].toSource +'; ' + '<br>';
//      	 }
     	 
//      	 }

  console.log('reqListener called');
  console.log( oReq.responseText);
  console.log( oReq.responseXML);
  console.log( oReq.response);
}

var ws;

function getThemArea() {

	ws = new WebSocket("ws://localhost:8090/");

	ws.onopen = function() {
	    // alert("Opened!");
	    // ws.send("Hello Server");
	};

	ws.onmessage = function (evt) {
		
		
		var arrayResult = evt.data.split(" ");
		
		if( arrayResult[0] == "INIT") 
		{
		var index = 1;
			
		var _parisLatTop = parseFloat( arrayResult[index++] );
		var _parisLngWest = parseFloat( arrayResult[index++]);
		
		var _parisLatBottom = parseFloat( arrayResult[index++]);
		var _parisLngEast = parseFloat( arrayResult[index++]);

		var parisNorthWest = new google.maps.LatLng( _parisLatTop, _parisLngWest);
		var parisSouthEast = new google.maps.LatLng( _parisLatBottom, _parisLngEast);

		parisBounds.extend( parisNorthWest);
		parisBounds.extend( parisSouthEast);
		
		// map.setCenter( {lat: (_parisLatTop + _parisLatBottom) /2.0 , lng: (_parisLngEast + _parisLngWest)/2.0});
		
		var _haitiLatTop = parseFloat( arrayResult[index++] );
		var _haitiLngWest = parseFloat( arrayResult[index++]);
		
		var _haitiLatBottom = parseFloat( arrayResult[index++]);
		var _haitiLngEast = parseFloat( arrayResult[index++]);

		var haitiNorthWest = new google.maps.LatLng( _haitiLatTop, _haitiLngWest);
		var haitiSouthEast = new google.maps.LatLng( _haitiLatBottom, _haitiLngEast);

		haitiBounds.extend( haitiNorthWest);
		haitiBounds.extend( haitiSouthEast);
		
		var _senegalLatTop = parseFloat( arrayResult[index++] );
		var _senegalLngWest = parseFloat( arrayResult[index++]);
		
		var _senegalLatBottom = parseFloat( arrayResult[index++]);
		var _senegalLngEast = parseFloat( arrayResult[index++]);

		var senegalNorthWest = new google.maps.LatLng( _senegalLatTop, _senegalLngWest);
		var senegalSouthEast = new google.maps.LatLng( _senegalLatBottom, _senegalLngEast);

		senegalBounds.extend( senegalNorthWest);
		senegalBounds.extend( senegalSouthEast);
		
		var _newYorkLatTop = parseFloat( arrayResult[index++] );
		var _newYorkLngWest = parseFloat( arrayResult[index++]);
		
		var _newYorkLatBottom = parseFloat( arrayResult[index++]);
		var _newYorkLngEast = parseFloat( arrayResult[index++]);

		var newYorkNorthWest = new google.maps.LatLng( _newYorkLatTop, _newYorkLngWest);
		var newYorkSouthEast = new google.maps.LatLng( _newYorkLatBottom, _newYorkLngEast);

		newYorkBounds.extend( newYorkNorthWest);
		newYorkBounds.extend( newYorkSouthEast);
		
		// map.fitBounds( parisBounds);
		  var rectangleSource = new google.maps.Rectangle({
			    strokeColor: '#FF0000',
			    strokeOpacity: 0.8,
			    strokeWeight: 2,
			    fillColor: '#FF0000',
			    fillOpacity: 0.05,
			    map: map,
			    bounds: parisBounds
			  });

			google.maps.event.addListener( rectangleSource, 'click', function(e) {
			    placeParisMarker(e.latLng, map);
			});

			  var rectangleDestination = new google.maps.Rectangle({
				    strokeColor: '#FF0000',
				    strokeOpacity: 0.8,
				    strokeWeight: 2,
				    fillColor: '#FF0000',
				    fillOpacity: 0.05,
				    map: map,
				    bounds: haitiBounds
				  });

				google.maps.event.addListener( rectangleDestination, 'click', function(e) {
				    placeHaitiMarker(e.latLng, map);
				});

				  var rectangleDestination = new google.maps.Rectangle({
					    strokeColor: '#FF0000',
					    strokeOpacity: 0.8,
					    strokeWeight: 2,
					    fillColor: '#FF0000',
					    fillOpacity: 0.05,
					    map: map,
					    bounds: senegalBounds
					  });

					google.maps.event.addListener( rectangleDestination, 'click', function(e) {
					    placeSenegalMarker(e.latLng, map);
					});

			var rectangleNewYork = new google.maps.Rectangle({
				    strokeColor: '#FF0000',
				    strokeOpacity: 0.8,
				    strokeWeight: 2,
				    fillColor: '#FF0000',
				    fillOpacity: 0.05,
				    map: map,
				    bounds: newYorkBounds
				  });

				google.maps.event.addListener( rectangleNewYork, 'click', function(e) {
				    placeNewYorkMarker(e.latLng, map);
				});

		}
		else if ( arrayResult[0] != "MAP") // this is a position following a translation
			{
			var index = 0;
			var destinationArea = arrayResult[index++];
			
			// get result position
			var _latDest = parseFloat( arrayResult[index++] );
			var _lngDest = parseFloat( arrayResult[index++]);
			
			// get second position
			var _lat2Dest = parseFloat( arrayResult[index++] );
			var _lng2Dest = parseFloat( arrayResult[index++]);
			
			  var marker = new google.maps.Marker({
				    position: {lat: _latDest, lng: _lngDest},
				    map: map,
				    Opacity: 0.35,
				    zIndex: -10
				  });
			  
			  var marker_shape = {coords: [0,0,50,50], type: "rect"}
			  var markerSecond = new google.maps.Marker({
				    position: {lat: _lat2Dest, lng: _lng2Dest},
				    map: map,
				    Opacity: 0.25,
				    zIndex: -20,
				    shape: marker_shape
				  });
				  map.panTo( {lat: _latDest, lng: _lngDest} );
			}
			else { // this is the "PARIS_MAP" answer
				var index = 1;
				
			// we get the number of borders
			var borderNumbers = parseInt( arrayResult[index++]);
				
			for (var borderId = 0; borderId < borderNumbers; borderId++) {
			
			 var borderSizeString_l = arrayResult[index++];
			 var borderSizeInt = parseInt( borderSizeString_l);

			var _latLastBorderPoint = parseFloat( arrayResult[index++] );
			 var _lngLastBorderPoint = parseFloat( arrayResult[index++]); 
			  for (var i = 0; i < borderSizeInt - 1; i++) {


			  	var _latBorderPoint = parseFloat( arrayResult[index++] );
			    var _lngBorderPoint = parseFloat( arrayResult[index++]);
			    
			    var lastBorderPoint = new google.maps.LatLng( _latLastBorderPoint, _lngLastBorderPoint);
			    var borderPoint = new google.maps.LatLng( _latBorderPoint, _lngBorderPoint);
			    
			    var segment = [ lastBorderPoint, borderPoint];
			  var line = new google.maps.Polyline({
				    path: segment,
 				   strokeColor: '#FF0000',
    				strokeOpacity: 1,
    				strokeWeight: 1
				  });
				  _latLastBorderPoint = _latBorderPoint;
				  _lngLastBorderPoint = _lngBorderPoint;
				  line.setMap(map);
			 }

			}
			}
			
			
			
			if (destinationArea == "HAITI") haitiPosition = {lat: _latDest, lng: _lngDest};
			if (destinationArea == "PARIS") parisPosition = {lat: _latDest, lng: _lngDest};
			if (destinationArea == "NEWYORK") newYorkPosition = {lat: _latDest, lng: _lngDest};
		// alert("Message: " + evt.data);
	    
	};



	ws.onclose = function() {
	    alert("Closed!");
	};

	ws.onerror = function(err) {
	    alert("Error: " + err);
	};

//     var oReq = new XMLHttpRequest();
//     oReq.responseType = "string"
//     oReq.onload = reqListener( oReq);
//     oReq.open("get", "http://localhost:8080/com.vogella.jersey.first/rest/hello", true);
//     oReq.send();


    // alert('getThemArea called');
}

function placeParisMarker(position, map) {
	  var marker = new google.maps.Marker({
	    position: position,
	    map: map
	  });
	  parisPosition = position;
	  map.panTo(position);
	  
	  // var messagePosition = "PARIS " + parisPosition.lat().toString() + " " + parisPosition.lng().toString();
	  // ws.send( messagePosition);
	  // alert("Message sent: " + messagePosition);
}

function placeNewYorkMarker(position, map) {
	  var marker = new google.maps.Marker({
	    position: position,
	    map: map
	  });
	  newYorkPosition = position;
	  map.panTo(position);
	  
		// alert("Message sent: " + messagePosition);
}

function placeHaitiMarker(position, map) {
	  var marker = new google.maps.Marker({
	    position: position,
	    map: map
	  });
	  haitiPosition = position;
	  map.panTo(position);
	  
		// alert("Message sent: " + messagePosition);
}

function placeSenegalMarker(position, map) {
	  var marker = new google.maps.Marker({
	    position: position,
	    map: map
	  });
	  senegalPosition = position;
	  map.panTo(position);
	  
		// alert("Message sent: " + messagePosition);
}

function getParisMap() {

	  var messagePosition = "PARIS_MAP " + parisPosition.lat().toString() + " " + parisPosition.lng().toString();
	  ws.send( messagePosition);
}

function getHaitiMap() {

	  var messagePosition = "HAITI_MAP " + haitiPosition.lat().toString() + " " + haitiPosition.lng().toString();
	  ws.send( messagePosition);
	  
}

function getNewYorkMap() {

	  var messagePosition = "NEW_YORK_MAP " + newYorkPosition.lat().toString() + " " + newYorkPosition.lng().toString();
	  ws.send( messagePosition);
	  
}

function translateParisHaiti() {

	  var messagePosition = "PARIS_HAITI " + parisPosition.lat().toString() + " " + parisPosition.lng().toString();
	  ws.send( messagePosition);
	  
}

function translateHaitiParis() {

	  var messagePosition = "HAITI_PARIS " + haitiPosition.lat().toString() + " " + haitiPosition.lng().toString();
	  ws.send( messagePosition);
	  
}

function translateParisSenegal() {

	  var messagePosition = "PARIS_SENEGAL " + parisPosition.lat().toString() + " " + parisPosition.lng().toString();
	  ws.send( messagePosition);
	  
}

function translateSenegalParis() {

	  var messagePosition = "SENEGAL_PARIS " + senegalPosition.lat().toString() + " " + senegalPosition.lng().toString();
	  ws.send( messagePosition);
	  
}

function translateNewYorkHaiti() {

	  var messagePosition = "NEWYORK_HAITI " + newYorkPosition.lat().toString() + " " + newYorkPosition.lng().toString();
	  ws.send( messagePosition);
	  
}

function translateHaitiNewYork() {

	  var messagePosition = "HAITI_NEWYORK " + haitiPosition.lat().toString() + " " + haitiPosition.lng().toString();
	  ws.send( messagePosition);
	  
}

function paris3Neighbors() {
	
	  var messagePosition = "PARIS_NEIGHBORS " + haitiPosition.lat().toString() + " " + haitiPosition.lng().toString();
	  ws.send( messagePosition);
	
}

function callback(response, status) {
  if (status != google.maps.DistanceMatrixStatus.OK) {
    alert('Error was: ' + status);
  } else {
    var origins = response.originAddresses;
    var destinations = response.destinationAddresses;
    var outputDiv = document.getElementById('outputDiv');
    outputDiv.innerHTML = '';
    deleteOverlays();

    for (var i = 0; i < origins.length; i++) {
      var results = response.rows[i].elements;
      addMarker(origins[i], false);
      for (var j = 0; j < results.length; j++) {
        addMarker(destinations[j], true);
        outputDiv.innerHTML += origins[i] + ' to ' + destinations[j]
            + ': ' + results[j].distance.text + ' in '
            + results[j].duration.text + '<br>';
      }
    }
  }
}

function addMarker(location, isDestination) {
  var icon;
  if (isDestination) {
    icon = destinationIcon;
  } else {
    icon = originIcon;
  }
  geocoder.geocode({'address': location}, function(results, status) {
    if (status == google.maps.GeocoderStatus.OK) {
      bounds.extend(results[0].geometry.location);
      map.fitBounds(bounds);
      var marker = new google.maps.Marker({
        map: map,
        position: results[0].geometry.location,
        icon: icon
      });
      markersArray.push(marker);
    } else {
      alert('Geocode was not successful for the following reason: '
        + status);
    }
  });
}

function deleteOverlays() {
  for (var i = 0; i < markersArray.length; i++) {
    markersArray[i].setMap(null);
  }
  markersArray = [];
}

google.maps.event.addDomListener(window, 'load', initialize);