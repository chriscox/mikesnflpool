/*global _:false */
'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:PlayerPicksCtrl
 * @description
 * # PlayerPicksCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('PlayerPicksCtrl', function ($scope, dataService, utils) {
    $scope.getUsers = function() {
      dataService.getUsers(function(users) {
        _.each(users, function(u) {
          u.wins = 0;
        });
        $scope.users = users;
        $scope.render();
      });
    };

    $scope.getGames = function() {
      dataService.getGames(function(games) {
        $scope.games = games;
      });
    };

    $scope.getAllUserPicks = function() {
      dataService.getAllUserPicks(function(userPicks) {
        $scope.userPicks = userPicks;
        $scope.render();
      }, function(){
      });
    };

    $scope.render = function() {
      // For each user pick, update user wins
      if ($scope.users && $scope.userPicks) {

        _.each($scope.userPicks, function(userPick) {
          userPick.team.teamKey = userPick.teamKey;
          var results = utils.getGameResults(userPick.game, userPick.team);
          var data = {};
          data.abbr = userPick.team.abbr;

          // get this user
          var thisUser = _.find($scope.users, function(u) {
            return (userPick.userKey === u.userKey);
          });

          if (results.spreadWin) {
            data.winner = true;
            thisUser.wins += 1;
          }
          thisUser[userPick.gameKey] = data;
        });

        // sort by firstName then wins
        $scope.users = _.chain($scope.users)
          .sortBy(function(u) { return u.firstName; })
          .sortBy(function(u) { return -u.wins; })
          .value();
          
      }
    };

    $scope.getUsers();
    $scope.getAllUserPicks();
    $scope.getGames();
  });
