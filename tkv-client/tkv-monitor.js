angular.module('MyApp',['ngMaterial', 'ngMessages'])

	.controller('RenderImage', ['$scope', '$timeout', '$http', function ($scope, $timeout, $http) {
	
		var self = this;

		// adjust values of x min compared to x max
		// $scope.$watch( function() {
	  
		// 	return $scope.area.x1 > $scope.area.x2;
		// }, function() {	
		// 		$scope.area.x1 = $scope.area.x2-0.01;
		// 	}
		// );	

		// adjust values of x min compared to x max
		// $scope.$watch( function() {
	  
		// 	return $scope.area.y1 > $scope.area.y2;
		// }, function() {	
		// 		$scope.area.y1 = $scope.area.y2-0.01;
		// 	}
		// );	

		$scope.area = {
			centerX: 0.5,
			centerY: 0.5,
			x1: 0.0,
			x2: 1.0,
			y1: 0.0,
			y2: 1.0,
			zoom: 1.0
		}

		$scope.DTpow10 = 0.0;
		$scope.theta = 0.5;

		this.updateDt = function() {

			console.log( $scope.DTpow10);
			var newDt = Math.pow( 10, $scope.DTpow10)/10.0;
			console.log( newDt);
			var jsondt = JSON.stringify( newDt);
			console.log( jsondt);
			
		
			$http.post('http://localhost:8000/dt', jsondt ).then
			(
				function(response) { // success handler
	  				console.log(response.status);
	  				console.log('updating dt');
	  			}, 
  				function(errResponse) { // error handler
      					console.error('error while posting dt');
      					console.error(errResponse);
  				}

  			);
		  	console.log('updateDt called');
	  	};

		this.updateTheta = function() {

			console.log( $scope.theta);
			var jsontheta = JSON.stringify( $scope.theta);
			console.log( jsontheta);
			
			$http.post('http://localhost:8000/theta', jsontheta ).then
			(
				function(response) { // success handler
	  				console.log(response.status);
	  				console.log('updating theta');
	  			}, 
  				function(errResponse) { // error handler
      					console.error('error while posting theta');
      					console.error(errResponse);
  				}

  			);
		  	console.log('updateTheta called');
	  	};


		this.run = function() {
			$http.get('http://localhost:8000/play').then(
				function(response) { // success handler
				},
				function(errResponse) { // error handler
      					console.error('error while request play');
  				}
			)	
		};

		this.oneStep = function() {
			$http.get('http://localhost:8000/oneStep').then(
				function(response) { // success handler
				},
				function(errResponse) { // error handler
      					console.error('error while request one');
  				}
			)	
		};

		this.pause = function() {
			$http.get('http://localhost:8000/pause').then(
				function(response) { // success handler
				},
				function(errResponse) { // error handler
      					console.error('error while request pause');
  				}
			)	
		};

		this.captureConfig = function() {
			$http.get('http://localhost:8000/captureConfig').then( function(response) {},
				function(errResponse) { console.error('error while request captureConfig');})	
		};

		this.updateArea = function() {

			$scope.area.x1 = $scope.area.centerX - 0.5/$scope.area.zoom
			$scope.area.y1 = $scope.area.centerY - 0.5/$scope.area.zoom
			$scope.area.x2 = $scope.area.centerX + 0.5/$scope.area.zoom
			$scope.area.y2 = $scope.area.centerY + 0.5/$scope.area.zoom

			var jsonarea = JSON.stringify( $scope.area);
			console.log( jsonarea);
			
		
			$http.post('http://localhost:8000/area', jsonarea ).then
			(
				function(response) { // success handler
	  				console.log(response.status);
	  				console.log('updating area');
	  				return $http.get('http://localhost:8000/render', '');
  				}).then( function(renderResponse) 
  					{
						$scope.render = 'data:image/gif;base64,' + renderResponse.data
  					}, 
  					function(errResponse) { // error handler
      					console.error('error while posting area');
      					console.error(errResponse);
  					}

  				);

		  	console.log('updateArea called');
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

	  	var pollRender = function() {

	  		$timeout( function() {
				$http.get('http://localhost:8000/render', '').then( function(renderResponse) 
  					{
						$scope.render = 'data:image/gif;base64,' + renderResponse.data
						pollRender()
  					}, 
  					function(errResponse) { // error handler
      					console.error('error while render');
      					console.error(errResponse);
  					}

  					);
	  			} 
	  			,  1000
	  		);
		};

		pollRender();

	  	var pollStatus = function() {

	  		$timeout( function() {
				$http.get('http://localhost:8000/status', '').then( function(StatusResponse) 
  					{
						$scope.status = StatusResponse.data
						pollStatus()
  					}, 
  					function(errResponse) { // error handler
      					console.error('error while Status');
      					console.error(errResponse);
  					}

  					);
	  			} 
	  			,  1000
	  		);
		};

		pollStatus();

	  	$timeout(function () {

	  	}, 1000);

	}]);
