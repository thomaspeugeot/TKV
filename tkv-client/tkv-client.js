

var app = angular.module("demoapp", ['leaflet-directive']);

app.controller("EventsController", [ '$scope', '$http', function($scope, $http) {

	angular.extend($scope, {
		center: {
			lat: 52.374004,
			lng: 4.890359,
			zoom: 7
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
	});
	
	$scope.markers = {};

	$scope.$on('leafletDirectiveMap.click', function(event, args){

		var jsonLatLng = JSON.stringify( args.leafletEvent.latlng);
		console.log( jsonLatLng);

		$http.post('http://localhost:8001/villageCoordinates', jsonLatLng ).then
		(
				function(response) { // success handler
					console.log(response.status);
					console.log('village villageCoordinates answer', response.data.X, response.data.Y);

					message = "village "+response.data.X+" "+response.data.Y;

					$scope.markers['user'] = {
						lat: args.leafletEvent.latlng.lat,
						lng: args.leafletEvent.latlng.lng,
						message: message, 
						focus: true,
						draggable: false
					};	
				}, 
  				function(errResponse) { // error handler
  					console.error('error while posting jsonLatLng');
  					console.error(errResponse);
  				}
  				);

	}); // end of click



}]);