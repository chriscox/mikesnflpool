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


    var _currentSeason = 2014;
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
      return 1;
      // var now = new Date();
      // if (now <= new Date("9/9/2013")) {
      //   return 1;
      // } else if (now <= new Date("9/16/2013")) {
      //   return 2;
      // } else if (now <= new Date("9/23/2013")) {
      //   return 3;
      // } else if (now <= new Date("9/30/2013")) {
      //   return 4;
      // } else if (now <= new Date("10/7/2013")) {
      //   return 5;
      // } else if (now <= new Date("10/14/2013")) {
      //   return 6;
      // } else if (now <= new Date("10/21/2013")) {
      //   return 7;
      // } else if (now <= new Date("10/28/2013")) {
      //   return 8;
      // } else if (now <= new Date("11/4/2013")) {
      //   return 9;
      // } else if (now <= new Date("11/11/2013")) {
      //   return 10;
      // } else if (now <= new Date("11/18/2013")) {
      //   return 11;
      // } else if (now <= new Date("11/25/2013")) {
      //   return 12;
      // } else if (now <= new Date("12/2/2013")) {
      //   return 13;
      // } else if (now <= new Date("12/9/2013")) {
      //   return 14;
      // } else if (now <= new Date("12/16/2013")) {
      //   return 15;
      // } else if (now <= new Date("12/23/2013")) {
      //   return 16;
      // } else {
      //   return 17;
      // }
    }

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

    // Register

    var register = function(firstName, lastName, email, password, remember, callback, onError) {
      Restangular.all('register').post({
        firstName:firstName,
        lastName:lastName,
        email:email,
        password:password
      }).then(function(response) {
        if (remember) {
          setUserCookie(response.user);
        }
        callback();
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

    var clearCookie = function() {
      $cookieStore.remove('user');
    };

    var getToken = function() {
      var token = $cookieStore.get('user.sessionToken');
      return {'sessionToken':token};
    };

    var login = function(email, password, callback, onError) {
      Restangular.all('login').post({
        email:email,
        password:password
      }).then(function(response) {
        setUserCookie(response.user);
        $rootScope.$broadcast('isAuthenticated', true);
        callback();
      }, function(error) {
        onError(error);
      });
    };

    var logout = function() {
      clearCookie();
      $rootScope.$broadcast('isAuthenticated', false);
      $location.path('/login');
    };

    var isAuthenticated = function() {
      return (getUserCookie() !== null && getUserCookie() !== undefined);
    };

    var authenticatedUser = function() {
      if (isAuthenticated) {
        return getUserCookie();
      } else {
        return null;
      }
    };

    // Users

    var getUsers = function(callback) {
      if (_users) {
        return callback(_users);
      } else {
        Restangular.all('users').getList({
        }).then(function(users) {
          _users = users;
          return callback(_users);
        }, function(error) {
          parseError(error);
        });
      }
    };

    // UserPicks

    var getUserPicks = function(userId, callback, onError) {
      Restangular.one('season', _activeSeason)
        .one('user', userId)
        .all('userpicks').getList()
        .then(function(userPicks) {
          callback(userPicks);
        }, function(error) {
          onError(error);
        });
    };

    var getAllUserPicks = function(callback) {
      Restangular.one('season', _activeSeason)
        .all('userpicks').getList({week: getActiveWeek()})
        .then(function(userPicks) {
          callback(userPicks);
        });
    };

    var addUserPick = function(game, user, callback, onError) {
      Restangular.all('userpicks').post({
        user: (user === null) ? authenticatedUser() : user,
        game:game
      }).then(function(pick) {
        callback(pick);
      }, function(error) {
        onError(error);
      });
    };

    var getUserStats = function(callback) {
      Restangular.one('season', getActiveSeason()).all('userstats').getList({
      }).then(function(userStats) {
        callback(userStats);
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
        });
    };

    var addTeam = function(team, callback) {
      Restangular.all('teams').post({
        name: team.name,
        abbr: team.abbr
      }).then(function(team) {
        callback(team);
      });
    };

    // Games

    var getGames = function(callback) {
      Restangular.one('season', getActiveSeason()).all('games').getList({
        week: getActiveWeek()
      }).then(function(games) {
        callback(games);
      });
    };

    var getAllGames = function(callback) {
      Restangular.one('season', getActiveSeason()).all('games').getList()
        .then(function(games) {
          callback(games);
        });
    };

    var deleteGame = function(game, callback, onError) {
      Restangular.one('games', game.id).remove().then(function() {
        callback();
      }, function(error) {
        onError(error);
      });
    };

    var addGame = function(game, callback) {
      Restangular.all('games').post({
        Season: getActiveSeason(),
        Week: getActiveWeek(),
        date: game.date,
        // homeTeam: game.homeTeam,
        // awayTeam: game.awayTeam,
      }).then(function(game) {
        callback(game);
      }, function(error) {
        console.log(error);
      });
    };

    var updateGame = function(game, callback, onError) {
      Restangular.all('games').post({
        season: getActiveSeason(),
        week: getActiveWeek(),
        game: game
      }).then(function(game) {
        callback(game);
      }, function(error) {
        onError(error);
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

      register: register,
      login: login,
      logout: logout,
      isAuthenticated: isAuthenticated,
      authenticatedUser: authenticatedUser,

      getUsers: getUsers,
      getUserPicks: getUserPicks,
      getAllUserPicks: getAllUserPicks,
      addUserPick: addUserPick,
      getUserStats: getUserStats,

      getTeams: getTeams,
      getTeamSchedule: getTeamSchedule,
      getTeamStandings: getTeamStandings,
      addTeam: addTeam,

      getGames: getGames,
      getAllGames: getAllGames,
      deleteGame: deleteGame,
      addGame: addGame,
      updateGame: updateGame,

    };
  });
