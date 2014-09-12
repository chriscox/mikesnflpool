'use strict';

angular.module('clientApp')
  .controller('PasswordResetCtrl', function ($scope, dataService) {

  	$scope.submit = function() {
      dataService.passwordReset($scope.email, function() {      	
        $scope.hasError = false;
        $scope.hasSuccess = true;
      }, function(error) {
      	console.log(error)
        $scope.hasError = true;
        $scope.hasSuccess = false;
      });
    };

  });
