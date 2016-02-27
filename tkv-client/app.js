angular.module("app", ["chart.js"])


  // Optional configuration
  .config(['ChartJsProvider', function (ChartJsProvider) {
    // Configure all charts
    ChartJsProvider.setOptions({
      colours: ['#FF5252', '#FF8A80'],
      responsive: false
    });
    // Configure all line charts
    ChartJsProvider.setOptions('Line', {
      datasetFill: false
    });
  }])
  .controller("LineCtrl", ['$scope', '$timeout', '$http', function ($scope, $timeout, $http) {

    var nbInterPolation = 10
    $scope.labels = [];
    for (i = 0; i<=nbInterPolation+1; i++) {
      $scope.labels.push( i )
    }
    $scope.series = [];
    for (i = 0; i<=9; i++) {
      $scope.series.push( i )
    }
    
    var self = this;
    self.breakdowns = [];

    $http.get('http://localhost:8000/stats').then(function(response) {
      $scope.data = [];
      console.log('response ', response.data);
      //
      // create serie from step 0 to max step
      for (j = 0; j<=9; j++) {
      
        var breakdown = [];
        breakdown.push( response.data[0][j]);
        
        for (i = 1; i<=nbInterPolation; i++) {
          var step = Math.floor( i*response.data.length/(nbInterPolation+1));
          breakdown.push( response.data[ step][j]);
        }
        breakdown.push( response.data[response.data.length -1 ][j]);        

        $scope.data.push( breakdown);
      }
      // console.log('scope data', $scope.data);
      // console.log('got response length', response.data.length);
      // console.log('got response [0] length', response.data[0].length);
    }, function(errResponse) {
      console.error('Error while fetching stats');
    });

    $scope.onClick = function (points, evt) {
      console.log(points, evt);
    };

    // Simulate async data update
    $timeout(function () {
    }, 1000);


}])


  ;
