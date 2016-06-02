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
			centerX: 0.492,
			centerY: 0.267,
			x1: 0.0,
			x2: 1.0,
			y1: 0.0,
			y2: 1.0,
			zoom: 1.0
		}
		$scope.zoomPow10 = 0
		
		$scope.ratioBorderBodies = 0.1;

		$scope.DTpow10 = -7.0;
		$scope.theta = 0.5;


		$scope.nbVillagesPerAxe = 100;
		$scope.nbRoutines = 100;

		$scope.toto = "toto";

		$scope.dirConfig = ""

		this.updateDt = function() {

			console.log( $scope.DTpow10);
			var newDt = Math.pow( 10, $scope.DTpow10);
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

		this.updateNbVillagesPerAxe = function() {

			console.log( $scope.nbVillagesPerAxe);
			var json = JSON.stringify( $scope.nbVillagesPerAxe);
			console.log( json);
			
			$http.post('http://localhost:8000/nbVillagesPerAxe', json ).then
			(
				function(response) { // success handler
	  				console.log(response.status);
	  				console.log('updating nbVillagesPerAxe');
	  			}, 
  				function(errResponse) { // error handler
      					console.error('error while posting theta');
      					console.error(errResponse);
  				}

  			);
		  	console.log('updateTheta called');
	  	};

		this.updateNbRoutines = function() {

			console.log( $scope.nbRoutines);
			var json = JSON.stringify( $scope.nbRoutines);
			console.log( json);
			
			$http.post('http://localhost:8000/nbRoutines', json ).then
			(
				function(response) { // success handler
	  				console.log(response.status);
	  				console.log('updating nbRoutines');
	  			}, 
  				function(errResponse) { // error handler
      					console.error('error while posting theta');
      					console.error(errResponse);
  				}

  			);
		  	console.log('updateTheta called');
	  	};


		this.updateRatioBorderBodies = function() {

			console.log( $scope.ratioBorderBodies);
			var json = JSON.stringify( $scope.ratioBorderBodies);
			console.log( json);
			
			$http.post('http://localhost:8000/updateRatioBorderBodies', json ).then
			(
				function(response) { // success handler
	  				console.log(response.status);
	  				console.log('updating updateRatioBorderBodies');
	  			}, 
  				function(errResponse) { // error handler
      					console.error('error while updateRatioBorderBodies');
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

		this.toggleRenderChoice  = function() {
			$http.get('http://localhost:8000/toggleRenderChoice').then( function(response) {},
				function(errResponse) { console.error('error while request toggleRenderChoice');})	
		};

		this.toggleLocalGlobal  = function() {
			$http.get('http://localhost:8000/toggleLocalGlobal').then( function(response) {},
				function(errResponse) { console.error('error while request toggleLocalGlobal');})	
		};

		this.toggleManualAuto  = function() {
			$http.get('http://localhost:8000/toggleManualAuto').then( function(response) {},
				function(errResponse) { console.error('error while request toggleManualAuto');})	
		};

		// fetch coordinates of minimal distance
		this.zoomSpecial = function() {
			$http.get('http://localhost:8000/minDistanceCoord').then( function(response) {
				// get X and Y
				$scope.area.centerX = parseFloat(response.data.X)
				$scope.area.centerY = parseFloat(response.data.Y)
			},
				function(errResponse) { console.error('error while request minimal distance');})	
		};

		this.updateArea = function() {
			
			$scope.area.zoom = Math.pow( 10, $scope.zoomPow10)

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

	  	var pollDensityTenciles = function() {

	  		$timeout( function() {
				$http.get('http://localhost:8000/getDensityTenciles', '').then( function(response) 
  					{
  						console.log( response.data)
						$scope.densityTenciles = response.data
						pollDensityTenciles()
  					}, 
  					function(errResponse) { // error handler
      					console.error('error while Status');
      					console.error(errResponse);
  					}

  					);
	  			} 
	  			,  2000
	  		);
		};

		pollDensityTenciles();



		// function that list the files available
		var pullConfigs = function() {
			$http.get('http://localhost:8000/dirConfig', '').then( function(response) 
				{
					$scope.dirConfig = response.data
					$scope.selected = response.data[0]
				}, 
				function(errResponse) { // error handler
					console.error('error while Status');
					console.error(errResponse);
				}

				);
		}

		pullConfigs();

		$scope.selected = "";

		this.loadConfig = function() {

			$http.get('http://localhost:8000/loadConfig'+'?file='+$scope.selected, '').then( function(response) 
				{
					console.log("file loaded " + $scope.selected);
				}, 
				function(errResponse) { // error handler
					console.error('error while Status');
					console.error(errResponse);
				}

				);
		}

		this.loadConfigOrig = function() {

			$http.get('http://localhost:8000/loadConfigOrig'+'?file='+$scope.selected, '').then( function(response) 
				{
					console.log("file orig loaded " + $scope.selected);
				}, 
				function(errResponse) { // error handler
					console.error('error while Status');
					console.error(errResponse);
				}

				);
		}

	}]);
		