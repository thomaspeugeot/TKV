

var app = angular.module("demoapp", ['leaflet-directive']);

app.controller("EventsController", [ '$scope', function($scope) {
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
        }
    });


	$scope.$on('leafletDirectiveMap.click', function(event, args){
    	console.log(args.leafletEvent.latlng);
	});
}]);