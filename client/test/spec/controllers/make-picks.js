'use strict';

describe('Controller: MakePicksCtrl', function () {

  // load the controller's module
  beforeEach(module('clientApp'));

  var MakePicksCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    MakePicksCtrl = $controller('MakePicksCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
