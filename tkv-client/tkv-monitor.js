angular.module('MyApp',['ngMaterial', 'ngMessages'])

.controller('AppCtrl', ['$scope', '$timeout', '$http', function ($scope, $timeout, $http) {

  // get the image
  $http.get('http://localhost:8000/render').then(function(response) {
    $scope.render = response.data
  }, function(errResponse) {
      console.error('Error while fetching render');
  });

  $timeout(function () {
  }, 1000);

}]);
