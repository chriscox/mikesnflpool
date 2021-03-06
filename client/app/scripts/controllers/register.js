/*global _:false */
'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:RegisterCtrl
 * @description
 * # RegisterCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('RegisterCtrl', function ($scope, dataService) {
    if (dataService.adminUser()) {
      $scope.closed = false;
    } else {
      $scope.closed = true;
    }
    
    $scope.user = {
      'fields': [
        {
          'id':'firstName',
          'type':'text',
          'name':'First Name',
          'hasError':false,
          'value':''
        },
        {
          'id':'lastName',
          'type':'text',
          'name':'Last Name',
          'hasError':false,
          'value':''
        },
        {
          'id':'email',
          'type':'email',
          'name':'Email',
          'hasError':false,
          'value':''
        },
        {
          'id':'password',
          'type':'password',
          'name':'Password',
          'hasError':false,
          'value':''
        },
        {
          'id':'passwordConfirm',
          'type':'password',
          'name':'Confirm Password',
          'hasError':false,
          'value':''
        }
      ]
    };

    $scope.hasError = function() {
      return (_.findWhere($scope.user.fields, { hasError:true }) != null);
    };

    $scope.clearErrors = function() {
      _.each($scope.user.fields, function(field) {
        field.hasError = false;
      });
    };

    $scope.getTournaments = function() {
      dataService.getTournaments(function(tournaments) {
        // TODO: Update this to handle more than a single single tournament. 
        $scope.tournamentKey = tournaments[0].tournamentKey;
      });
    };


    $scope.register = function() {
      var user = $scope.user.fields;
      var firstName = user[0];
      var lastName = user[1];
      var email = user[2];
      var password = user[3];
      var passwordConfirm = user[4];

      // reset errors
      $scope.clearErrors();

      // ensure all fields non-empty
      _.each(user.fields, function(field) {
        if (field.value === '') {
          field.hasError = true;
          $scope.errorText = 'All fields are required.';
        }
      });

      // test password confirm
      if (password.value !== passwordConfirm.value) {
        password.hasError = true;
        passwordConfirm.hasError = true;
        $scope.errorText = 'Password confirmation must match.';
      }

      // submit
      if (!$scope.hasError()) {
        dataService.register(
          $scope.tournamentKey,
          firstName.value, 
          lastName.value, 
          email.value, 
          password.value,
        function() {
          $scope.hasSuccess = true;
        }, function(error) {
          if (error.status === 400) {
            email.hasError = true;
            $scope.errorText = error.data;
          }
          $scope.hasSuccess = false;
        });
      }
    };

    $scope.getTournaments();
  });
