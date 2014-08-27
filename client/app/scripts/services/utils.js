'use strict';

/**
 * @ngdoc service
 * @name clientApp.utils
 * @description
 * # utils
 * Service in the clientApp.
 */
angular.module('clientApp')
  .service('utils', function utils() {
    var getGameResults = function(game, pickedTeam) {

      var result  = {
        win: false,
        loss: false,
        tie: false,
        spreadWin: false,
        spreadLoss: false
      };
      if (game.ended) {
        result.tie = isGameTie(game);
        if (!result.tie) {
          result.win = isGameWinner(game, pickedTeam);
          result.loss = !result.win;
        }
        result.spreadWin = isSpreadWinner(game, pickedTeam);
        result.spreadLoss = !result.spreadWin;
      }
      return result;
    };

    var isGameTie = function(game) {
      return  (game.awayTeamScore === game.homeTeamScore);
    };

    var isSpreadWinner = function(game, pickedTeam) {
      if ( (game.awayTeamScore - game.awayTeamSpread) > (game.homeTeamScore - game.homeTeamSpread) ){
        return (pickedTeam.teamKey === game.awayTeamKey);
      } else if ( (game.awayTeamScore - game.awayTeamSpread) < (game.homeTeamScore - game.homeTeamSpread) ){
        return (pickedTeam.teamKey === game.homeTeamKey);
      } else {
        // If score minus spread equals other teams, then give losing 
        // team the spread win. Teams must win spread + 1, not be equal.
        var isWinner = isGameWinner(game, game.awayTeam);
        if (!isWinner) {
          return (pickedTeam.teamKey === game.awayTeamKey);
        } else {
          return (pickedTeam.teamKey === game.homeTeamKey);
        }
      }
    };

    var isGameWinner = function(game, pickedTeam) {
      if (game.awayTeamScore > game.homeTeamScore){
        return (pickedTeam.id === game.awayTeamKey);
      } else if (game.awayTeamScore < game.homeTeamScore){
        return (pickedTeam.id === game.homeTeamKey);
      } else {
        return false;
      }
    };

    // Returns 1st, 2nd, 3rd, etc...
    var getOrdinal = function(n) {
      var s = ['th','st','nd','rd'];
      var v =  n%100;
      return n+(s[(v-20)%10]||s[v]||s[0]);
    };

    // Public API

    return {

      getGameResults: getGameResults,
      getOrdinal: getOrdinal

    };
  });
