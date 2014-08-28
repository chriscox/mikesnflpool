'use strict';

describe('Controller: TeamStatsCtrl', function () {

  // load the controller's module
  beforeEach(module('clientApp'));

  var TeamStatsCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    TeamStatsCtrl = $controller('TeamStatsCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
