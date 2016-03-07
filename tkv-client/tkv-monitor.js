angular.module('MyApp',['ngMaterial', 'ngMessages'])

	.controller('RenderImage', ['$scope', '$timeout', '$http', function ($scope, $timeout, $http) {

		var self = this;

		// adjust values of x min compared to x max
		$scope.$watch( function() {
	  
			return $scope.area.x1 > $scope.area.x2;
		}, function() {	
				$scope.area.x1 = $scope.area.x2;
			}
		);	

		// adjust values of x min compared to x max
		$scope.$watch( function() {
	  
			return $scope.area.y1 > $scope.area.y2;
		}, function() {	
				$scope.area.y1 = $scope.area.y2;
			}
		);	

		$scope.area = {
			x1: 0.0,
			x2: 1.0,
			y1: 0.0,
			y2: 1.0
		}

		this.updateArea = function() {
		  	console.log('updateArea called');
	  	};

		this.submit = function() {
		  	console.log('Submit called');
	  	};

	  	// get the image
	  	$http.get('http://localhost:8000/render', '').then(function(response) {
		    
			$scope.render = 'data:image/gif;base64,' + response.data

		  	}, function(errResponse) {
	      	console.error('Error while fetching render');
	  	});

	  	$http.get('http://localhost:8000/status', '').then(function(response) {
		    
			$scope.status = response.data

		  	}, function(errResponse) {
	      	console.error('Error while fetching render');
	  	});

	  	$timeout(function () {

	  	}, 1000);

	}]);
