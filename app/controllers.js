(function(){
    'use strict';

    angular.module('clientRestControllers', [])
        .controller('addUserController', [
        '$scope', 
        '$http', 
        function($scope, $http){
            console.log("addUserController active.");

            this.getUsers = function(){
                console.log($scope.user);
                if ($scope.id)
                    $http.get("http://localhost:3000/user/" + $scope.id, null , {
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    })
                        .success(function(response){
                        console.log(response);
                    })
                        .error(function(response){
                        console.log(response);
                    });
                else
                    $http.get("http://localhost:3000/user")
                        .then(function(response){
                        console.log(response);
                    });
            }

            this.sendUser = function(){
                console.log($scope.user);
                $http.post("http://localhost:3000/user", $scope.user, {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                })
                    .then(function(response){
                    console.log(response);
                });
            }
        }])
        .controller('authenticationController', [
        '$scope', 
        '$http', 
        function($scope, $http){
            console.log("authenticationController active.");
            $scope.authUser = {}
            
            this.authenticate = function(){
                console.log($scope.authUser);
                $http.post("http://localhost:3000/token", $scope.authUser, {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                })
                .then(function(response){
                    console.log(response);
                });
            }
        }]);
})();