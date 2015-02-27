/*global _:false */
'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:LeaderboardCtrl
 * @description
 * # LeaderboardCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('LeaderboardCtrl', function ($scope, dataService, utils) {

    var _users = null;
    var _realUsers = null; // without robots

    $scope.getUsers = function() {
      dataService.getUsers(function(users) {
        _users = users;
        $scope.render();
      });
    };

    $scope.getUserStats = function() {
      dataService.getUserStats(function(userStats) {
        $scope.userStats = userStats;
        $scope.render();
      });
    };

    $scope.render = function() {
      // For each user, update stats
      if (_users && $scope.userStats) {

        _.each(_users, function(user) {

          user.robot = (user.firstName.substring(0, 2) === '[ ');
          user.moneyWins = [];
          user.moneyTotal = 0;
          user.total = 0;

          for (var i=0; i<17; i++) {
            var wins = $scope.userStats.stats[user.userKey][i+1];
            if (wins) {
              user[i] = wins;
              user.total += wins;
            } else {
              user[i] = null;
            }
          }
        });

        // sort by firstName then wins
        _users = _.chain(_users)
          .sortBy(function(u) { return u.firstName; })
          .sortBy(function(u) { return -u.total; })
          .value();

        // determine place
        var userTotal = 0;
        var userPlace = 1;
        _.each(_users, function(user) {
          if (user.total > userTotal) {
            user.place = utils.getOrdinal(userPlace);
          } else if (user.total < userTotal) {
            userPlace += 1;
            user.place = utils.getOrdinal(userPlace);
          } else {
            user.place = utils.getOrdinal(userPlace);
          }
          userTotal = user.total;
        });

        $scope.users = _users;

        // determine money winners
        var pot = 0;
        _realUsers = _.filter($scope.users, function(u) {
          return !u.robot;
        });

        for (var i=0; i<17; i++) {
          if (i >= dataService.getCurrentWeek() - 1) {
            break;
          }
          var sortedUsers = _.sortBy(_realUsers, function(u) {
            return -u[i];
          });
          if (sortedUsers[0][i] > sortedUsers[1][i]) {
            sortedUsers[0].moneyWins[i] = 16 + pot;
            sortedUsers[0].moneyTotal += 16 + pot;
            pot = 0;
          } else {
            // increase pot
            pot += 16;
          }
        }

        // sort by firstName then money wins
        $scope.realUsers = _.chain(_realUsers)
          .sortBy(function(u) { return u.firstName; })
          .sortBy(function(u) { return -u.moneyTotal; })
          .value();
      }

    };

    $scope.getUsers();
    $scope.getUserStats();
  });
