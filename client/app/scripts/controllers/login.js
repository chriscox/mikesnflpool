'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:LoginCtrl
 * @description
 * # LoginCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('LoginCtrl', function ($scope, dataService) {

    $scope.remember = true;

    $scope.login = function() {
      dataService.login($scope.email, $scope.password, function() {
        $scope.hasError = false;
        $scope.hasSuccess = true;
      }, function(error) {
        $scope.hasError = true;
        $scope.hasSuccess = false;
      });
    };
  });
