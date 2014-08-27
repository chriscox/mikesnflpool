/*global _:false */
'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:AdminGamesCtrl
 * @description
 * # AdminGamesCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('AdminGamesCtrl', function ($scope, $timeout, dataService) {
    
    $scope.getTeams = function() {
      dataService.getTeams(function(teams) {
        $scope.teams = teams;
      });
    };

    $scope.getGames = function() {
      $scope.games =[];
      dataService.getGames(function(games) {
        $scope.games = games;
      });
    };

    $scope.deleteGame = function(game) {
      dataService.deleteGame(game, function() {
        $scope.games = _.without($scope.games, game);
      }, function() {
        game.hasError = true;
        $timeout(function() {
          game.hasError = false;
        }, 2000);
      });
    };

    $scope.addGame = function() {
      dataService.addGame($scope.newGame, function(game) {
        $scope.games.push(game);
        $scope.newGame.homeTeam = null;
        $scope.newGame.awayTeam = null;
      }, function() {
        $scope.error = true;
        $timeout(function() {
          $scope.error = false;
        }, 5000);
      });
    };

    $scope.updateGame = function(game) {
      dataService.updateGame(game, function() {
      }, function() {
      });
    };

    $scope.getTeams();
    $scope.getGames();
  });
