'use strict';

describe('Controller: AdminUserpicksCtrl', function () {

  // load the controller's module
  beforeEach(module('clientApp'));

  var AdminUserpicksCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    AdminUserpicksCtrl = $controller('AdminUserpicksCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
