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
    'restangular',
    'ngAnimate',
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'ngTouch'
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl'
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

  }).run(function(Restangular, $location) {
    if ($location.host() === 'localhost') {
      // Set to local play server
      Restangular.setBaseUrl('http://localhost:8080/api');
    } else {
      // For dist
      Restangular.setBaseUrl('/api');
    }
    // Restangular.setRequestSuffix('.json');
    // Restangular.setRestangularFields({
    //   id: "_id",
    //   route: "restangularRoute",
    //   selfLink: "self.href"
    // });
  });
