'use strict';

/**
 * @ngdoc directive
 * @name clientApp.directive:weekSelector
 * @description
 * # weekSelector
 */
angular.module('clientApp')
  .directive('weekSelector', function (dataService) {
    return {
      templateUrl: 'views/week-selector.html',
      restrict: 'AE',
      replace: true,
      link: function postLink(scope) {

        scope.pool = {
          season: dataService.getActiveSeason(),
          week: dataService.getActiveWeek(),
          setWeek: function(week) {
            dataService.setRouteParam({week:week});
          }
        };
      }
    };
  });
