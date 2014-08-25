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

    $scope.getTournaments = function() {
      dataService.getTournaments(function(tournaments) {
        $scope.tournaments = tournaments;
      });
    };

    $scope.addTournament = function() {
      var name = $scope.tournamentName;
      var season = $scope.tournamentSeason;
      dataService.addTournament(name, season, function(tournament) {
        $scope.tournaments = tournament;
      });
    };

    $scope.getTournaments();
  });
