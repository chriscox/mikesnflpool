'use strict';

/**
 * @ngdoc filter
 * @name clientApp.filter:range
 * @function
 * @description
 * # range
 * Filter in the clientApp.
 */
angular.module('clientApp')
  .filter('range', function () {
    return function (input, total) {
      total = parseInt(total, 10);
      for (var i=0; i<total; i++) {
        input.push(i);
      }
      return input;
    };
  });
