'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:AboutCtrl
 * @description
 * # AboutCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('AboutCtrl', function ($scope, Restangular) {
    $scope.teams = [
      {Name:'123'},
      {Name:'456'}
    ];

    Restangular.all('teams').getList({
        }).then(function(teams) {
          $scope.teams = teams;
        }, function(error) {

        });
  });
