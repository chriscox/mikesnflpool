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

    $scope.addTeam = function() {
      dataService.addTeam($scope.newTeam, function(team) {
        $scope.teams.push(team);
        $scope.newTeam.name = null;
        $scope.newTeam.abbr = null;
      });
    };

    $scope.getTeams();
  });
