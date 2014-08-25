'use strict';

describe('Controller: AdminTournamentsCtrl', function () {

  // load the controller's module
  beforeEach(module('clientApp'));

  var AdminTournamentsCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    AdminTournamentsCtrl = $controller('AdminTournamentsCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
