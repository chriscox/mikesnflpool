'use strict';

/**
 * @ngdoc function
 * @name clientApp.controller:PasswordResetCtrl
 * @description
 * # PasswordResetCtrl
 * Controller of the clientApp
 */
angular.module('clientApp')
  .controller('PasswordResetCtrl', function ($scope, dataService) {

    $scope.user = {
      'fields': [
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

    $scope.submit = function() {
      var user = $scope.user.fields;
      var password = user[0];
      var passwordConfirm = user[1];

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
        dataService.passwordReset (
          password.value,
        function() {
          $scope.hasSuccess = true;
        }, function(error) {
          console.log(error)
          if (error.status === 400) {
            password.hasError = true;
            $scope.errorText = error.data;
          }
          $scope.hasSuccess = false;
        });
      }
    };

  });
