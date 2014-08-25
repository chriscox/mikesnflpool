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
      dataService.getGames(function(games) {
        $scope.games = games;
      });
    };

    $scope.deleteGame = function(game) {
      dataService.deleteGame(game, function() {
        $scope.games = _.without($scope.games, game);
      }, function(error) {
        game.hasError = true;
        $timeout(function() {
          game.hasError = false;
        }, 2000);
      });
    };

    $scope.addGame = function() {
      console.log($scope.newGame)
      dataService.addGame($scope.newGame, function(game) {
        $scope.games.push(game);
        $scope.newGame.homeTeam = null;
        $scope.newGame.awayTeam = null;
      }, function(error) {
        $scope.error = true;
        $timeout(function() {
          $scope.error = false;
        }, 5000);
      });
    };

    $scope.updateGame = function(game) {
      dataService.updateGame(game, function(game) {
      }, function(error) {
      });
    };

    $scope.getTeams();
    $scope.getGames();
  });
