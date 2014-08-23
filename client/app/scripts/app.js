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
    // 'ui.router',
    'ng-token-auth',
  ])
  .config(function ($routeProvider, $authProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
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
      .otherwise({
        redirectTo: '/'
      });

    // // ng-token-auth config.
    // if ($location.host() === 'localhost') {
      $authProvider.configure({
        apiUrl: 'http://localhost:8080/api',
      });
    // }

  }).run(function(Restangular, $location) {
    // Restangular config
    if ($location.host() === 'localhost') {
      Restangular.setBaseUrl('http://localhost:8080/api');
    } else {
      Restangular.setBaseUrl('/api');
    }
  });
