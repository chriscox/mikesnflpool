/*global _:false */
'use strict';

/**
 * @ngdoc service
 * @name clientApp.data
 * @description
 * # data
 * Service in the clientApp.
 */
angular.module('clientApp')
  .service('dataService', function dataService($location, $http, Restangular, $cookieStore, $rootScope) {

    var _currentSeason = 2015;
    var _currentWeek = null;

    var _activeSeason = _currentSeason;
    var _activeWeek = null;
    var _users;

    var getActiveRoute = function() {
      return $location.path();
    };

    var setActiveRoute = function(path) {
      $location.path(path);
    };

    var getRouteParam = function(key, fallback) {
      var value = $location.search()[key];
      if (!value && fallback) {
        var param = {};
        param[key] = fallback;
        setRouteParam(param);
      }
      return (value) ? value : fallback;
    };

    var setRouteParam = function(param) {
      var params = $location.search();
      for (var key in param) {
        if (param.hasOwnProperty(key)){
          params[key] = param[key];
        }
      }
      $location.search(params);
    };

    var getActiveSeason = function() {
      return _activeSeason;
    };

    var setActiveSeason = function(season) {
      _activeSeason = season;
    };

    var getCurrentWeek = function() {
      var now = new Date();
      if (now <= new Date("9/14/2015")) {
        return 1;
      } else if (now <= new Date("9/21/2015")) {
        return 2;
      } else if (now <= new Date("9/28/2015")) {
        return 3;
      } else if (now <= new Date("10/5/2015")) {
        return 4;
      } else if (now <= new Date("10/12/2015")) {
        return 5;
      } else if (now <= new Date("10/19/2015")) {
        return 6;
      } else if (now <= new Date("10/26/2015")) {
        return 7;
      } else if (now <= new Date("11/2/2015")) {
        return 8;
      } else if (now <= new Date("11/9/2015")) {
        return 9;
      } else if (now <= new Date("11/16/2015")) {
        return 10;
      } else if (now <= new Date("11/23/2015")) {
        return 11;
      } else if (now <= new Date("11/30/2015")) {
        return 12;
      } else if (now <= new Date("12/7/2015")) {
        return 13;
      } else if (now <= new Date("12/14/2015")) {
        return 14;
      } else if (now <= new Date("12/21/2015")) {
        return 15;
      } else if (now <= new Date("12/28/2015")) {
        return 16;
      } else {
        return 17;
      }
    };

    var getActiveWeek = function() {
      var week =  $location.search().week;
      if (!week) {
        _currentWeek = getCurrentWeek();
        if (_activeWeek === null) {
          _activeWeek = _currentWeek;
        }
        week = (_currentWeek !== _activeWeek) ? _activeWeek : _currentWeek;
        setRouteParam({week:week});
      }
      week = parseInt(week);
      _activeWeek = week;
      return week;
    };

    // Errors

    var parseError = function(error) {
      if (error.status === 401) {
        $location.path('/login');
      }
    };

    // Bots

    var registerBot = function(tournamentKey, firstName, botType, callback, onError) {
      Restangular.all('auth').post({
        tournamentKey:tournamentKey,
        firstName:firstName,
        email:botType,
        botType:botType,
        bot:true,
      }).then(function(user) {
        _users = null;
        callback(user);
      }, function(error) {
        onError(error);
      });
    };

    var getBots = function(callback) {
      Restangular.one('tournament', authenticatedUser().tournamentKey)
      .all('users').getList().then(function(users) {
        var botUsers = _.filter(users, function(u) {
          return u.bot;
        });
        return callback(botUsers);
      }, function(error) {
        parseError(error);
      });
    };

    // Register

    var register = function(tournamentKey, firstName, lastName, email, password, callback, onError) {
      Restangular.all('auth').post({
        tournamentKey:tournamentKey,
        firstName:firstName,
        lastName:lastName,
        email:email,
        password:password
      }).then(function(user) {
        setUserCookie(user);
        $rootScope.$broadcast('isAuthenticated', true);
        callback(user);
      }, function(error) {
        onError(error);
      });
    };

    // login/logout

    var setUserCookie = function(user) {
      $cookieStore.put('user', user);
    };

    var getUserCookie = function() {
      return $cookieStore.get('user');
    };

    var clearCookies = function() {
      $cookieStore.remove('user');
    };

    var passwordForgot = function(email, callback, onError) {
      Restangular.all('passwordforgot').post({
        email:email,
      }).then(function() {
        callback();
      }, function(error) {
        onError(error);
      });
    };

    var passwordReset = function(password, callback, onError) {
      Restangular.all('passwordreset').post({
        token:getRouteParam('token'),
        password:password,
      }).then(function() {
        callback();
      }, function(error) {
        onError(error);
      });
    };

    var login = function(email, password, callback, onError) {
      Restangular.all('login').post({
        email:email,
        password:password
      }).then(function(user) {
        setUserCookie(user);
        $rootScope.$broadcast('isAuthenticated', true);
        $rootScope.$broadcast('isAdmin', user.admin);
        callback(user);
      }, function(error) {
        onError(error);
      });
    };

    var logout = function() {
      clearCookies();
      $rootScope.$broadcast('isAuthenticated', false);
      $rootScope.$broadcast('isAdmin', false);
      $location.path('/login');
    };

    var isAuthenticated = function() {
      return (getUserCookie() != null && getUserCookie() != undefined);
    };

    var authenticatedUser = function() {
      if (isAuthenticated) {
        return getUserCookie();
      } else {
        return null;
      }
    };

    var adminUser = function() {
      if (isAuthenticated()) {
        var cookie = getUserCookie();
        if ('admin' in cookie) {
          return cookie.admin;
        }
      }
      return false;
    };

    // Users

    var getUsers = function(callback) {
      if (_users) {
        return callback(_users);
      } else {
        Restangular.one('tournament', authenticatedUser().tournamentKey)
        .all('users').getList().then(function(users) {
          _users = users;
          return callback(users);
        }, function(error) {
          parseError(error);
        });
      }
    };

    // UserPicks

    var getUserPicks = function(user, callback, onError) {
      Restangular.one('tournament', authenticatedUser().tournamentKey)
      .one('season', getActiveSeason())
      .one('user', user.userKey)
      .all('userpicks').getList({
        week: getActiveWeek(),
      }).then(function(userPicks) {
        callback(userPicks);
      }, function(error) {
        onError(error);
      });
    };

    var getAllUserPicks = function(callback, onError) {
      Restangular.one('tournament', authenticatedUser().tournamentKey)
        .one('season', getActiveSeason())
        .all('userpicks').getList({
          week: getActiveWeek()
        }).then(function(userPicks) {
          callback(userPicks);
        }, function(error) {
          onError(error);
        });
    };

    var addUserPick = function(game, user, callback, onError) {
      Restangular.all('userpicks').post({
        tournamentKey:user.tournamentKey,
        userKey:user.userKey,
        game:game
      }).then(function(pick) {
        callback(pick);
      }, function(error) {
        onError(error);
      });
    };

    var getUserStats = function(callback) {
      Restangular.one('tournament', authenticatedUser().tournamentKey)
        .one('season', getActiveSeason())
        .one('userstats').get()
        .then(function(userStats) {
          callback(userStats);
        }, function() {
        });
    };

    // Teams

    var getTeams = function(callback) {
      Restangular.all('teams').getList().then(function(teams) {
        callback(teams);
      });
    };

    var getTeamSchedule = function(callback) {
      Restangular.one('season', getActiveSeason())
        .one('teams', getRouteParam('team', 'sf'))
        .all('schedule').getList()
        .then(function(schedule) {
          callback(schedule);
        });
    };

    var getTeamStandings = function(teamAbbr, summarize, callback) {
      Restangular.one('season', getActiveSeason())
        .one('teams', teamAbbr)
        .all('standings').getList({
          week: (summarize) ? 17 : getActiveWeek()
        }).then(function(standings) {
          callback(standings);
        }, function(error) {
        });
    };

    // Games

    var getGames = function(callback) {
       Restangular.one('season', getActiveSeason()).all('games').getList({
        week: getActiveWeek()
      }).then(function(games) {
        callback(games);
      }, function() {
      });
    };

    var getAllGames = function(callback) {
      Restangular.one('season', getActiveSeason()).all('games').getList()
        .then(function(games) {
          callback(games);
        });
    };

    var deleteGame = function(game, callback, onError) {
      // TODO: Convert this to DELETE from POST
      // Restangular.one('games', game.gameKey).remove().then(function() {
      //   callback();
      // }, function(error) {
      //   onError(error);
      // });
      Restangular.one('season', getActiveSeason())
        .one('week', getActiveWeek())
        .one('deletegame', game.gameKey).post()
        .then(function() {
          callback();
        }, function(error) {
          onError(error);
        });
    };

    var addGame = function(game, callback, onError) {
      Restangular.all('games').post({
        season: getActiveSeason(),
        week: getActiveWeek(),
        date: game.date,
        homeTeamAbbr: game.homeTeam,
        awayTeamAbbr: game.awayTeam
      }).then(function(game) {
        callback(game);
      }, function(error) {
        onError(error);
      });
    };

    var updateGame = function(game, callback, onError) {
      Restangular.all('games').post({
        season: getActiveSeason(),
        week: getActiveWeek(),
        gameKey: game.gameKey,
        awayTeamScore: game.awayTeamScore,
        awayTeamSpread: game.awayTeamSpread,
        homeTeamScore: game.homeTeamScore,
        homeTeamSpread: game.homeTeamSpread,
        ended: game.ended
      }).then(function(game) {
        callback(game);
      }, function(error) {
        onError(error);
      });
    };

    // Admin

    var addTeams = function(callback) {
      $http.get('scripts/teams.json').success(function(data) {
        _.each(data.teams, function(team) {
          Restangular.all('teams').post(team).then(function() {
            callback(data.teams);
          });
        });
      });
    };

    var addTournament = function(name, season, callback) {
      Restangular.all('tournaments').post({
        name:name,
        season:parseInt(season)
      }).then(function(tournament) {
        callback(tournament);
      }, function() {

      });
    };

    var getTournaments = function(callback) {
      Restangular.one('season', getActiveSeason()).all('tournaments').getList()
      .then(function(tournaments) {
        callback(tournaments);
      });
    };

    // Public API

    return {
      getActiveRoute: getActiveRoute,
      setActiveRoute: setActiveRoute,
      
      getRouteParam: getRouteParam,
      setRouteParam: setRouteParam,

      getActiveSeason: getActiveSeason,
      setActiveSeason: setActiveSeason,

      getActiveWeek: getActiveWeek,
      getCurrentWeek: getCurrentWeek,

      registerBot: registerBot,
      getBots: getBots,

      register: register,
      login: login,
      logout: logout,
      passwordForgot: passwordForgot,
      passwordReset: passwordReset,
      isAuthenticated: isAuthenticated,
      authenticatedUser: authenticatedUser,
      adminUser: adminUser,

      getUsers: getUsers,
      getUserPicks: getUserPicks,
      getAllUserPicks: getAllUserPicks,
      addUserPick: addUserPick,
      getUserStats: getUserStats,

      getTeams: getTeams,
      getTeamSchedule: getTeamSchedule,
      getTeamStandings: getTeamStandings,
      addTeams: addTeams,

      getGames: getGames,
      getAllGames: getAllGames,
      deleteGame: deleteGame,
      addGame: addGame,
      updateGame: updateGame,

      addTournament: addTournament,
      getTournaments: getTournaments

    };
  });
