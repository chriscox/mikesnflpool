'use strict';

/**
 * @ngdoc directive
 * @name clientApp.directive:navbar
 * @description
 * # navbar
 */
angular.module('clientApp')
  .directive('navbar', function (dataService, $location) {
    return {
      templateUrl: 'views/navbar.html',
      restrict: 'AE',
      replace: true,
      link: function postLink(scope) {

        // scope.navItems = [
        //   {title:'Rules', url:'/rules'},
        //   {title:'Make Picks', url:'/make-picks'},
        //   {title:'Player Picks', url:'/player-picks'},
        //   {title:'Leaderboard', url:'leaderboard'},
        //   {title:'Team Stats', url:'team-stats'},
        //   {title:'Player Stats', url:'player-stats'},
        //   {title:'Log In', url:'login'},
        // ];

        // initial state
        scope.user = dataService.authenticatedUser();

        // Manually set route to keep the location.search value.
        // Otherwise controllers load twice for search '?week=X'
        scope.routeTo = function(route) {
          $location.path(route);
        };

        // watch route to determine active nav
        scope.$watch(function() {
          return dataService.getActiveRoute();
        }, function(newValue, oldValue) { 
          // toggle navbar mobile selector
          if (angular.element('.navbar-collapse').hasClass('in')) {
            angular.element('.navbar-toggle').click();
          }
          scope.route = newValue;

          // Check if admin
          scope.admin = dataService.adminUser();
        });

        scope.logout = function() {
          dataService.logout();
        };

        // listen for login/logout authenticated broadcast
        scope.$on('isAuthenticated', function(event, authenticated) {
          if (authenticated) {
            scope.user = dataService.authenticatedUser();
          } else {
            scope.user = null;
          }
        });

        scope.$on('isAdmin', function(event, admin) {
          scope.admin = dataService.adminUser();
        });
      }
    };
  });
