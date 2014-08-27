'use strict';

describe('Controller: PlayerPicksCtrl', function () {

  // load the controller's module
  beforeEach(module('clientApp'));

  var PlayerPicksCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    PlayerPicksCtrl = $controller('PlayerPicksCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
