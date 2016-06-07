angular.module('MyApp',['ngMaterial', 'ngMessages'])

	.controller('RenderImage', ['$scope', '$timeout', '$http', function ($scope, $timeout, $http) {
	
		var vm = this;

		// adjust values of x min compared to x max
		$scope.$watch( 'vm.area', function() {	
				vm.updateArea();
			}, true // http://stackoverflow.com/questions/19455501/watch-an-object
		);	

		// adjust values of x min compared to x max
		// $scope.$watch( function() {
	  
		// 	return $scope.area.y1 > $scope.area.y2;
		// }, function() {	
		// 		$scope.area.y1 = $scope.area.y2-0.01;
		// 	}
		// );	

		vm.area = {
			centerX: 0.492,
			centerY: 0.267,
			x1: 0.0,
			x2: 1.0,
			y1: 0.0,
			y2: 1.0,
			zoom: 1.0,
			gridNb: 1,
			zoomPow10: 0
		}
		vm.GridNbPow10 = 0
		
		vm.ratioBorderBodies = 0.1;

		vm.DTpow10 = -7.0;
		vm.theta = 0.5;


		vm.nbVillagesPerAxe = 100;
		vm.nbRoutines = 100;

		vm.toto = "toto";

		vm.dirConfig = ""

		vm.newImageCenter = function($event) {
			console.log( "newImageCenter ", $event, $event.target.x, $event.target.y );
			
		}
		
		vm.updateDt = function() {

			console.log( vm.DTpow10);
			var newDt = Math.pow( 10, vm.DTpow10);
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

		vm.updateGridNbPow10 = function() {

			console.log( vm.GridNbPow10);
			vm.gridNb = Math.pow( 10, vm.GridNbPow10);
			console.log( vm.gridNb);
			var jsonGridNb = JSON.stringify( vm.gridNb);
			console.log( jsonGridNb);
			
		
			$http.post('http://localhost:8000/fieldGridNb', jsonGridNb ).then
			(
				function(response) { // success handler
	  				console.log(response.status);
	  				console.log('updating dt');
	  			}, 
  				function(errResponse) { // error handler
      					console.error('error while posting jsonGridNb');
      					console.error(errResponse);
  				}

  			);
		  	console.log('updateGridNbPow10 called');
	  	};

		vm.updateTheta = function() {

			console.log( vm.theta);
			var jsontheta = JSON.stringify( vm.theta);
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

		vm.updateNbVillagesPerAxe = function() {

			console.log( vm.nbVillagesPerAxe);
			var json = JSON.stringify( vm.nbVillagesPerAxe);
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

		vm.updateNbRoutines = function() {

			console.log( vm.nbRoutines);
			var json = JSON.stringify( vm.nbRoutines);
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


		vm.updateRatioBorderBodies = function() {

			console.log( vm.ratioBorderBodies);
			var json = JSON.stringify( vm.ratioBorderBodies);
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


		vm.run = function() {
			$http.get('http://localhost:8000/play').then(
				function(response) { // success handler
				},
				function(errResponse) { // error handler
      					console.error('error while request play');
  				}
			)	
		};

		vm.oneStep = function() {
			$http.get('http://localhost:8000/oneStep').then(
				function(response) { // success handler
				},
				function(errResponse) { // error handler
      					console.error('error while request one');
  				}
			)	
		};

		vm.pause = function() {
			$http.get('http://localhost:8000/pause').then(
				function(response) { // success handler
				},
				function(errResponse) { // error handler
      					console.error('error while request pause');
  				}
			)	
		};

		vm.captureConfig = function() {
			$http.get('http://localhost:8000/captureConfig').then( function(response) {},
				function(errResponse) { console.error('error while request captureConfig');})	
		};

		vm.toggleRenderChoice  = function() {
			$http.get('http://localhost:8000/toggleRenderChoice').then( function(response) {},
				function(errResponse) { console.error('error while request toggleRenderChoice');})	
		};

		vm.toggleFieldRendering  = function() {
			$http.get('http://localhost:8000/toggleFieldRendering').then( function(response) {},
				function(errResponse) { console.error('error while request toggleFieldRendering');})	
		};

		vm.toggleLocalGlobal  = function() {
			$http.get('http://localhost:8000/toggleLocalGlobal').then( function(response) {},
				function(errResponse) { console.error('error while request toggleLocalGlobal');})	
		};

		vm.toggleManualAuto  = function() {
			$http.get('http://localhost:8000/toggleManualAuto').then( function(response) {},
				function(errResponse) { console.error('error while request toggleManualAuto');})	
		};

		// fetch coordinates of minimal distance
		vm.zoomSpecial = function() {
			$http.get('http://localhost:8000/minDistanceCoord').then( function(response) {
				// get X and Y
				vm.area.centerX = parseFloat(response.data.X);
				vm.area.centerY = parseFloat(response.data.Y);
				// vm.updateArea();
				},
				function(errResponse) { console.error('error while request minimal distance');})	
		};

		vm.centerArea = function() {
			vm.area.centerX = 0.5
			vm.area.centerY = 0.5
			// vm.updateArea()
		};


		vm.updateArea = function() {
			
			vm.area.zoom = Math.pow( 10, vm.area.zoomPow10)

			vm.area.x1 = vm.area.centerX - 0.5/vm.area.zoom
			vm.area.y1 = vm.area.centerY - 0.5/vm.area.zoom
			vm.area.x2 = vm.area.centerX + 0.5/vm.area.zoom
			vm.area.y2 = vm.area.centerY + 0.5/vm.area.zoom

			var jsonarea = JSON.stringify( vm.area);
			console.log( jsonarea);
			
		
			$http.post('http://localhost:8000/area', jsonarea ).then
			(
				function(response) { // success handler
	  				console.log(response.status);
	  				console.log('updating area');
	  				return $http.get('http://localhost:8000/render', '');
  				}).then( function(renderResponse) 
  					{
						vm.render = 'data:image/gif;base64,' + renderResponse.data
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
						vm.render = 'data:image/gif;base64,' + renderResponse.data
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
						vm.status = StatusResponse.data
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
  						// console.log( response.data)
						vm.densityTenciles = response.data
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
					vm.dirConfig = response.data
					vm.selected = response.data[0]
				}, 
				function(errResponse) { // error handler
					console.error('error while Status');
					console.error(errResponse);
				}

				);
		}

		pullConfigs();

		vm.selected = "";

		vm.loadConfig = function() {

			$http.get('http://localhost:8000/loadConfig'+'?file='+vm.selected, '').then( function(response) 
				{
					console.log("file loaded " + vm.selected);
				}, 
				function(errResponse) { // error handler
					console.error('error while Status');
					console.error(errResponse);
				}

				);
		}

		vm.loadConfigOrig = function() {

			$http.get('http://localhost:8000/loadConfigOrig'+'?file='+vm.selected, '').then( function(response) 
				{
					console.log("file orig loaded " + vm.selected);
				}, 
				function(errResponse) { // error handler
					console.error('error while Status');
					console.error(errResponse);
				}

				);
		}

	}]);
		