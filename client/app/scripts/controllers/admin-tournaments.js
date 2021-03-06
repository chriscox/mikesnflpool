'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:AdminTournamentsCtrl
 * @description
 * # AdminTournamentsCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('AdminTournamentsCtrl', function ($scope, dataService) {

    $scope.getAllTournaments = function() {
      dataService.getAllTournaments(function(tournaments) {
        $scope.tournaments = tournaments;
      });
    };

    $scope.addTournament = function() {
      var name = $scope.tournamentName;
      var season = $scope.tournamentSeason;
      dataService.addTournament(name, season, function(tournament) {
        $scope.getAllTournaments();
      });
    };

    $scope.test = function() {
      console.log('test')
    };

    $scope.getAllTournaments();
  });
