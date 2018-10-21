

var app = angular.module("demoapp", ['leaflet-directive']);

// var hostname = "https://tenktorg.appspot.com/"
var hostname
var protocol
var port
var targetService

app.controller("EventsController", [ '$scope', '$http', function($scope, $http) {

	angular.extend($scope, {
		france: {
			lat: 47.374004,
			lng: 4.890359,
			zoom: 5
		},
		haiti: {
			lat: 19,
			lng: -73,
			zoom: 5
		},
		defaults: {
			scrollWheelZoom: false
		},
		events: {
			map: {
				enable: ['click'],
				logic: 'emit'
			}
		},
		
		villageBorders : {
			data: {
			  "type": "FeatureCollection",
			  "features": [
				{
				  "type": "Feature",
				  "properties": {},
				  "geometry": {
					"type": "Polygon",
					"coordinates": []
				  }
				},
				{
				  "type": "Feature",
				  "properties": {},
				  "geometry": {
					"type": "Polygon",
					"coordinates": []
				  }
				}
				
			  ]
			},
			style: {
				fillColor: "green",
					weight: 2,
					opacity: 1,
					color: 'blue',
					dashArray: '3',
					fillOpacity: 0.1
			}
		}
	});
	
	$scope.markers = {};
	$scope.targetMarkers = {};

	$scope.$on('leafletDirectiveMap.click', function(event, args){

		var jsonLatLng = JSON.stringify( args.leafletEvent.latlng);
		console.log( jsonLatLng);

		console.log("post for translateLatLngInSourceCountryToLatLngInTargetCountry before");
				
		hostname = window.location.hostname
		protocol = window.location.protocol
		port = window.location.port
		targetService = protocol + "//"+ hostname + ":" + port + "/"

		$http.post( targetService +'translateLatLngInSourceCountryToLatLngInTargetCountry', jsonLatLng ).then
		(
			function(response) { // success handler
				console.log(response.status);
				console.log('village translateLatLngInSourceCountryToLatLngInTargetCountry answer', response.data.X, response.data.Y);

				message = "Territory X="+ Math.floor(100*response.data.X)+" Y="+Math.floor(100*response.data.Y);

				parseFloat(response.data.LatClosest);

				lat = parseFloat(response.data.LatClosest);
				lng = parseFloat(response.data.LngClosest);

				latTarget = parseFloat(response.data.LatTarget);
				lngTarget = parseFloat(response.data.LngTarget);

				xSpead = parseFloat(response.data.Xspread);
				ySpead = parseFloat(response.data.Yspread);

				message = "Territory X="+ Math.floor(100*response.data.Xspread)+" Y="+Math.floor(100*response.data.Yspread);
				
				$scope.targetMarkers['targetVillage'] = {
					lat: latTarget,
					lng: lngTarget,
					message: message, 
					focus: false,
					options: {
						noHide: true
					}
				}

				$scope.markers['clickPos'] = {
					lat: args.leafletEvent.latlng.lat,
					lng: args.leafletEvent.latlng.lng,
					message: message, 
					focus: true,
					draggable: false,
					options: {
						noHide: true
					}
				};	

				$http.post(targetService + 'villageTargetBorder', jsonLatLng ).then
				(
					function(response) { // success handler
						console.log(response.status);
						console.log('target village villageCoordinates before ', $scope.villageBorders.data.features[0].geometry.coordinates[0] )
						
						// convert response data field to float
						$scope.villageBorders.data.features[0].geometry.coordinates = [ [ [] ] ];
						$scope.villageBorders.data.features[0].geometry.coordinates[0] = new Array()
						$scope.villageBorders.data.features[0].geometry.coordinates[0].length = response.data[0].length
						
						for (var i = 0; i < response.data[0].length; i++) {
							$scope.villageBorders.data.features[0].geometry.coordinates[0][i] = new Array(2)
							$scope.villageBorders.data.features[0].geometry.coordinates[0][i][0] = parseFloat(response.data[0][i][0]);
							$scope.villageBorders.data.features[0].geometry.coordinates[0][i][1] = parseFloat(response.data[0][i][1]);
						}

						console.log('target village villageCoordinates answer ', $scope.villageBorders.data.features[0].geometry.coordinates[0] )

						$http.post(targetService + 'villageSourceBorder', jsonLatLng ).then
						(
							function(response) { // success handler
								console.log(response.status);
								console.log('source village villageCoordinates before ', $scope.villageBorders.data.features[1].geometry.coordinates[0] )		
										
								// convert response data field to float
								$scope.villageBorders.data.features[1].geometry.coordinates = [ [ [] ] ];
								$scope.villageBorders.data.features[1].geometry.coordinates[0] = new Array()
								$scope.villageBorders.data.features[1].geometry.coordinates[0].length = response.data[0].length
								
								for (var i = 0; i < response.data[0].length; i++) {
									$scope.villageBorders.data.features[1].geometry.coordinates[0][i] = new Array(2)
									$scope.villageBorders.data.features[1].geometry.coordinates[0][i][0] = parseFloat(response.data[0][i][0]);
									$scope.villageBorders.data.features[1].geometry.coordinates[0][i][1] = parseFloat(response.data[0][i][1]);
								}

								console.log('source village villageCoordinates answer ', $scope.villageBorders.data.features[1].geometry.coordinates[0] )

							}, 
							function(errResponse) { // error handler
								console.error('error while posting jsonLatLng');
								console.error(errResponse);
							}
						);
					}, 
					function(errResponse) { // error handler
						console.error('error while posting jsonLatLng');
						console.error(errResponse);
					}
				);
			}, 
			function(errResponse) { // error handler
				console.error('error while posting jsonLatLng');
				console.error(errResponse);
			}
		);
		
		console.log("post for translateLatLngInSourceCountryToLatLngInTargetCountry is over");


	}); // end of click
}]);
