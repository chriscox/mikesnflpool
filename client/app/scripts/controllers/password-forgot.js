'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:PasswordForgotCtrl
 * @description
 * # PasswordForgotCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('PasswordForgotCtrl', function ($scope, dataService) {
    
    $scope.submit = function() {
      dataService.passwordForgot($scope.email, function() {        
        $scope.hasError = false;
        $scope.hasSuccess = true;
      }, function(error) {
        $scope.hasError = true;
        $scope.hasSuccess = false;
      });
    };
  });
