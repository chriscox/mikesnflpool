'use strict';

describe('Controller: AdminGamesCtrl', function () {

  // load the controller's module
  beforeEach(module('clientApp'));

  var AdminGamesCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    AdminGamesCtrl = $controller('AdminGamesCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
