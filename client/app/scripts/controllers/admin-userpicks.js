'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:AdminUserpicksCtrl
 * @description
 * # AdminUserpicksCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('AdminUserpicksCtrl', function ($scope, dataService) {

    $scope.quickPicks = [
      {type:'home', title:'Home Teams'},
      {type:'away', title:'Away Teams'},
      {type:'favorite', title:'Favorites'},
      {type:'underdog', title:'Underdogs'},
      {type:'random', title:'Random'},
    ];

    $scope.pickCount = 0;

    $scope.getUsers = function() {
      dataService.getUsers(function(users) {
        $scope.users = _.sortBy(users, function(user) {
          return user.firstName;
        });
        $scope.user = users[0];
        $scope.getUserPicks();
      });
    };

    $scope.setUser = function(user) {
      $scope.user = user;
      $scope.getUserPicks();
    };

    $scope.getGames = function() {
      dataService.getGames(function(games) {
        $scope.games = games;
        $scope.render();
      });
    };

    $scope.getUserPicks = function() {
      $scope.userPicks = null;
      dataService.getUserPicks($scope.user, function(userPicks) {
        $scope.userPicks = userPicks;
        $scope.render();
      });
    };

    $scope.render = function() {
      // For each game, update with user picks
      if ($scope.games && $scope.userPicks) {
        var now = new Date();
        _.each($scope.games, function(game) {

          var pick = _.find($scope.userPicks, function(userPick) {
            return userPick.gameKey === game.gameKey;
          });

          if (pick) {
            $scope.pickCount += 1;
            game.awayTeam.selected = (game.awayTeamKey === pick.teamKey) ? true : false;
            game.homeTeam.selected = (game.homeTeamKey === pick.teamKey) ? true : false;
          } else {
            // clear picks
            game.awayTeam.selected = false;
            game.homeTeam.selected = false;
          }
        });
      }
    };

    $scope.selectTeam = function(game, playingAs) {
      game.awayTeam.selected = (playingAs === 'away') ? true : false;
      game.homeTeam.selected = (playingAs === 'home') ? true : false;
      $scope.addUserPick(game);
      $scope.getUserPickCount();
    };

    $scope.addUserPick = function(game) {
      dataService.addUserPick(game, $scope.user, function() {
        
      });
    };

    $scope.getUserPickCount = function() {
      var count = 0;
      _.each($scope.games, function(game) {
        if (game.awayTeam.selected || game.homeTeam.selected) {
          count += 1;
        }
      });
      $scope.pickCount = count;
    };

    $scope.doQuickPick = function(item) {
      _.each($scope.games, function(game) {
        switch (item.type) {
          case 'home':
            $scope.selectTeam(game, 'home');
            break;
          case 'away':
            $scope.selectTeam(game, 'away');
            break;
          case 'favorite':
            if (game.homeTeamSpread >= game.awayTeamSpread) {
              $scope.selectTeam(game, 'home');
            } else {
              $scope.selectTeam(game, 'away');
            }
            break;
          case 'underdog':
            if (game.homeTeamSpread <= game.awayTeamSpread) {
              $scope.selectTeam(game, 'home');
            } else {
              $scope.selectTeam(game, 'away');
            }
            break;
          case 'random':
            var isHome = _.random(0, 1);
            $scope.selectTeam(game, (isHome === 1) ? 'home' : 'away');
            break;
          default:
            break;
        } 
      });
      $scope.getUserPickCount();
    };

    $scope.getUsers();
    $scope.getGames();
  });
