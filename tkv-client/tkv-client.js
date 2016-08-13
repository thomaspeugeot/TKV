

var app = angular.module("demoapp", ['leaflet-directive']);

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

					message = "village "+response.data.X+" "+response.data.Y+" "+response.data.Distance+" "+response.data.LatClosest+" "+response.data.LngClosest;

					parseFloat(response.data.LatClosest);

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

					lat = 	parseFloat(response.data.LatClosest);
					lng =  parseFloat(response.data.LngClosest)

					$scope.markers['closestVillage'] = {
						lat: lat,
						lng: lng,
						message: 'closest village'
					}

				}, 
  				function(errResponse) { // error handler
  					console.error('error while posting jsonLatLng');
  					console.error(errResponse);
  				}
  				);

	}); // end of click
}]);
