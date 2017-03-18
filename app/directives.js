(function(){
    'use strict';

    angular.module('clientRestDirectives', [])
        .directive('dNavbar', function(){
            return {
                restrict: 'E',
                replace: true,
                templateUrl: "pub/view/navbar.html"
            }
        })
        .directive('connect', function(){
            return {
                restrict: 'E',
                replace: true,
                templateUrl: "pub/view/connect.html"
            }
        })
        .directive('dUserForm', function(){
            return {
                restrict: 'E',
                replace: true,
                templateUrl: "pub/view/user_form.html"
            }
        });
})();