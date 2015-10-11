'use strict';

describe('Controller: AdminTournamentPlayersCtrl', function () {

  // load the controller's module
  beforeEach(module('clientApp'));

  var AdminTournamentPlayersCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    AdminTournamentPlayersCtrl = $controller('AdminTournamentPlayersCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
