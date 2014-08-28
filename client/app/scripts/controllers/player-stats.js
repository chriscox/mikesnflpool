/*global _:false */
'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:PlayerStatsCtrl
 * @description
 * # PlayerStatsCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('PlayerStatsCtrl', function ($scope, dataService, utils) {

    $scope.winsPerSeason = 0;
    $scope.totalFinalGames = 0;
    $scope.pickStats = { home:0, away:0, favored:0, underdogs:0 };

    $scope.getUsers = function() {
      dataService.getUsers(function(users) {
        $scope.users = _.sortBy(users, function(user) {
          return user.firstName;
        });
        $scope.user = _.find(users, function(user) {
          return ($scope.activeUserKey == user.userKey);
        });
      });
    };

    $scope.getAllGames = function() {
      dataService.getAllGames(function(games) {
        $scope.games = _.chain(games)
          .groupBy(function(game) { return game.week; })
          .sortBy(function(game) { return game.week; })
          .value();
        $scope.totalFinalGames = _.where(games, {ended:true}).length;
        $scope.gamesById = _.indexBy(games, 'id');
        $scope.render();
      });
    };

    $scope.getUserPicks = function() {
      dataService.getUserPicks($scope.activeUserKey, function(userPicks) {
        $scope.userPicks = userPicks;
        $scope.render();
      });
    };

    $scope.render = function() {
      if ($scope.gamesById && $scope.userPicks) {

        var winsPerSeason = 0;
        var pickStats = { home:0, favored:0 };

        _.each($scope.userPicks, function(userPick) {
          var game = $scope.gamesById[userPick.game.id];
          var results = utils.getGameResults(game, userPick.team);
          userPick.results = results;
          game.pick = userPick;

          // Subtotal wins
          if (results.spreadWin) {
            var index = game.week - 1;
            if (!$scope.games[index].spreadWins) {
              $scope.games[index].spreadWins = 0;
            }
            $scope.games[index].spreadWins += 1;
            winsPerSeason += 1;
          }

          // Pick stats
          if (game.awayTeam.id === userPick.team.id) {
            $scope.pickStats.away += 1;
            if (game.awaySpread > game.homeSpread) {
              $scope.pickStats.favored += 1;
            } else {
              $scope.pickStats.underdogs += 1;
            }
          } else {
            $scope.pickStats.home += 1;
            if (game.homeSpread > game.awaySpread) {
              $scope.pickStats.favored += 1;
            } else {
              $scope.pickStats.underdogs += 1;
            }
          }
        });

        $scope.winsPerSeason = winsPerSeason;
        $scope.buildChart();
      }
    };

    $scope.setUser = function(user) {
      $scope.user = user;
      dataService.setRouteParam({userKey:$scope.user.userKey});
    };

    $scope.activeUserKey = dataService.getRouteParam('userKey', 'ahBkZXZ-bWlrZXNuZmxwb29schoLEgRVc2VyIhBjb3guY2hyaXNAbWUuY29tDA');
    $scope.getUsers();
    $scope.getAllGames();
    $scope.getUserPicks();

    $scope.xFunction = function(){
      return function(d) { return d.key; };
    };

    $scope.yFunction = function(){
      return function(d) { return d.y; };
    };

    $scope.buildChart = function() {
      $scope.winsData = [
        { key: "Wins", y: $scope.winsPerSeason },
        { key: "Losses", y: $scope.totalFinalGames - $scope.winsPerSeason }
      ];

      $scope.homeData = [
        { key: "Home Picks", y: $scope.pickStats.home },
        { key: "Away Picks", y: $scope.pickStats.away }
      ];

      $scope.favoredData = [
        { key: "Favored Picks", y: $scope.pickStats.favored },
        { key: "Underdogs Picks", y: $scope.pickStats.underdogs }
      ];
    }
  });
