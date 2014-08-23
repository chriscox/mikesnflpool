'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:AdminTeamsCtrl
 * @description
 * # AdminTeamsCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('AdminTeamsCtrl', function ($scope, dataService) {

    $scope.getTeams = function() {
      dataService.getTeams(function(teams) {
        $scope.teams = teams;
      });
    };

    $scope.addTeams = function() {
      dataService.addTeams(function(teams) {
        $scope.teams = teams;
      });
    };

    $scope.getTeams();
  });
