'use strict';

describe('Filter: substring', function () {

  // load the filter's module
  beforeEach(module('clientApp'));

  // initialize a new instance of the filter before each test
  var substring;
  beforeEach(inject(function ($filter) {
    substring = $filter('substring');
  }));

  it('should return the input prefixed with "substring filter:"', function () {
    var text = 'angularjs';
    expect(substring(text)).toBe('substring filter: ' + text);
  });

});
