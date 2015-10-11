'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:AdminTournamentPlayersCtrl
 * @description
 * # AdminTournamentPlayersCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('AdminTournamentPlayersCtrl', function ($scope, dataService) {

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

    $scope.getAllTournaments();
  });
