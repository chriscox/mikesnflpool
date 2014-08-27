'use strict';

describe('Controller: LeaderboardCtrl', function () {

  // load the controller's module
  beforeEach(module('clientApp'));

  var LeaderboardCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    LeaderboardCtrl = $controller('LeaderboardCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
