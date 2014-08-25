'use strict';

/**
 * @ngdoc overview
 * @name clientApp
 * @description
 * # clientApp
 *
 * Main module of the application.
 */
angular
  .module('clientApp', [
    'ngAnimate',
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'ngTouch',
    'ui.select2',
    'ui.bootstrap.datetimepicker',
    'ui.bootstrap',
    'restangular',
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .when('/make-picks', {
        templateUrl: 'views/make-picks.html',
        controller: 'MakePicksCtrl',
        access: 'user'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl'
      })
      .when('/register', {
        templateUrl: 'views/register.html',
        controller: 'RegisterCtrl'
      })
      .when('/login', {
        templateUrl: 'views/login.html',
        controller: 'LoginCtrl'
      })
      // admin
      .when('/admin/games', {
        templateUrl: 'views/admin-games.html',
        controller: 'AdminGamesCtrl',
        access: 'admin'
      })
      .when('/admin/teams', {
        templateUrl: 'views/admin-teams.html',
        controller: 'AdminTeamsCtrl',
        access: 'admin'
      })
      .when('/admin/tournaments', {
        templateUrl: 'views/admin-tournaments.html',
        controller: 'AdminTournamentsCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });

  }).run(function(Restangular, $location, $rootScope, dataService) {
    // Restangular config
    if ($location.host() === 'localhost') {
      Restangular.setBaseUrl('http://localhost:8080/api');
    } else {
      Restangular.setBaseUrl('/api');
    }
      // register listener to watch route changes
    $rootScope.$on( "$routeChangeStart", function(event, next, current) {
      var isAuthenticated = dataService.isAuthenticated();
      if (next.access == 'user' && !isAuthenticated) {
        // redirect
        $location.path("/login");
      } else if (next.access == 'admin' && !isAuthenticated) {
        // redirect
        $location.path("/login");
      }     
    });
  });
