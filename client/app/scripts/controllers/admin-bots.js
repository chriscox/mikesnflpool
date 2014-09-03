'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:AdminBotsCtrl
 * @description
 * # AdminBotsCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('AdminBotsCtrl', function ($scope, dataService) {

    $scope.botTypes = [
      {type:'home', title:'Home Teams'},
      {type:'away', title:'Away Teams'},
      {type:'favorite', title:'Favorites'},
      {type:'underdog', title:'Underdogs'},
      {type:'random', title:'Random'},
    ];

    $scope.getTournaments = function() {
      dataService.getTournaments(function(tournaments) {
        // TODO: Update this to handle more than a single single tournament. 
        $scope.tournamentKey = tournaments[0].tournamentKey;
      });
    };

    $scope.getGames = function() {
      dataService.getGames(function(games) {
        $scope.games = games;
      });
    };

    $scope.getBots = function() {
      dataService.getBots(function(bots) {
        $scope.botCounter = 0;
        $scope.bots = bots;
      });
    };

    $scope.addBot = function() {
      if ($scope.botName == undefined) {
        return;
      }
      dataService.registerBot(
        $scope.tournamentKey,
        $scope.botName,
        $scope.botType,
      function() {
        $scope.getBots();
      }, function(error) {
      });
    };

    $scope.selectTeam = function(game, playingAs, bot) {
      game.awayTeam.selected = (playingAs === 'away') ? true : false;
      game.homeTeam.selected = (playingAs === 'home') ? true : false;
      $scope.addUserPick(game, bot);
    };

    $scope.addUserPick = function(game, bot) {
      dataService.addUserPick(game, bot, function(pick) {
        if ($scope.botCounter < $scope.bots.length - 1) {
          $scope.botCounter ++;
          $scope.makeBotPicks();
        } else {
          $scope.buttonDisabled = false;
        }
      });
    };

    $scope.makeBotPicks = function() {
      $scope.buttonDisabled = true;
      var bot = $scope.bots[$scope.botCounter];
      _.each($scope.games, function(game) {
        switch (bot.botType) {
          case 'home':
            $scope.selectTeam(game, 'home', bot);
            break;
          case 'away':
            $scope.selectTeam(game, 'away', bot);
            break;
          case 'favorite':
            if (game.homeTeamSpread >= game.awayTeamSpread) {
              $scope.selectTeam(game, 'home', bot);
            } else {
              $scope.selectTeam(game, 'away', bot);
            }
            break;
          case 'underdog':
            if (game.homeTeamSpread <= game.awayTeamSpread) {
              $scope.selectTeam(game, 'home', bot);
            } else {
              $scope.selectTeam(game, 'away', bot);
            }
            break;
          case 'random':
            var isHome = _.random(0, 1);
            $scope.selectTeam(game, (isHome === 1) ? 'home' : 'away', bot);
            break;
          default:
            break;
        }
      });
    };

    $scope.getTournaments();
    $scope.getGames();
    $scope.getBots();
  });
