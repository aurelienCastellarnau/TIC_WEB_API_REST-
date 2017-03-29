(function(){
    'use strict';

    angular.module('clientRestControllers', [])
        .controller('userController', [
        '$scope', 
        '$http',
        '$cookies',
        '$base64',
        '$timeout',
        function($scope, $http, $cookies, $base64, $timeout){
            console.log("userController active.");
            var self = this;
            $scope.user = {};
            $scope.usersFromGet = [];
            $scope.search = "";
            $scope.count = 0;

            $scope.$watch('user', function(actual, old){
                if ($scope.autoUpdate)
                    self.editUser();
            })

            $scope.$watch('warning', function(actual, old){
                $timeout(function(){
                    $scope.warning = "";
                }, 10000);
            })

            this.getUsers = function(){
                if ($scope.id)
                {
                    $http.get("http://localhost:3000/user/" + $scope.id)
                        .then(function(r){
                        console.log(r);
                        $scope.userFromGet = r.data;
                    }, 
                              function(err){
                        console.log(err);
                        $scope.warning = manageError(err);
                    });

                }
                else
                    $http.get("http://localhost:3000/users")
                        .then(function(r){
                        console.log(r);
                        $scope.usersFromGet = r.data;
                    },
                              function(err){
                        console.log(err);
                        $scope.warning =  manageError(err);
                    });
            }

            this.sendUser = function(){
                $http.post("http://localhost:3000/user", $scope.user, {
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    xsrfCookieName: 'Auth',
                    withCredentials : true
                })
                    .then(function(r){
                    console.log(r);
                    $scope.usersFromGet = r.data;
                },
                          function(err){
                    console.log(err);
                    $scope.warning =  manageError(err);
                });
            }

            this.editUser = function(){
                if ($scope.user.uid <= 0 || typeof($scope.user.uid) == "undefined")
                    return;
                $http.put("http://localhost:3000/user/" + $scope.user.uid, $scope.user, {
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    xsrfCookieName: 'Auth',
                    withCredentials : true
                })
                    .then(function(r){
                    self.getUsers();
                    console.log(r);
                },
                          function(err){
                    console.log(err);
                    $scope.warning =  manageError(err);
                });
            }

            this.deleteUser = function(){
                if ($scope.id <= 0 || typeof($scope.id) == "undefined")
                    return;
                $http.delete("http://localhost:3000/user/" + $scope.id, null, {
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    xsrfCookieName: 'Auth',
                    withCredentials : true
                })
                    .then(function(response){
                    $scope.id = 0;
                    self.getUsers();
                    console.log(response);
                },
                          function(err){
                    console.log(err);
                    $scope.warning =  manageError(err);
                });
            }

            this.search = function() {
                var url = "http://localhost:3000/search/";
                console.log($scope.search);

                if ($scope.search != ""){
                    if ($scope.count > 0)
                        $scope.search = $scope.search + "/" + $scope.count;
                    console.log(url);
                    $http.get("http://localhost:3000/search/" + $scope.search)
                        .then(function(r){
                        $scope.usersFromGet = r.data;
                        console.log(r);
                    }, function (err){
                        console.log(err);
                        $scope.warning =  manageError(err);
                    });

                }
            }
        }])

        .controller('authenticationController', [
        '$rootScope',
        '$scope', 
        '$http',
        '$cookies', 
        '$timeout',
        '$base64',
        function($rootScope, $scope, $http, $cookies, $timeout, $base64){
            console.log("authenticationController active.");
            $scope.authUser = {}

            this.authenticate = function(){
                $rootScope.basic = $base64.encode($scope.authUser.email + ":" + $scope.authUser.password);
                $http.post("http://localhost:3000/auth", null, {
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    withCredentials : true
                })
                    .then(function(r){
                    if (r.status == 200){
                        $scope.authUser.email = "Thank You";
                        $scope.authUser.password = r.data.firstname + " " + r.data.lastname
                    }
                }
                          ,function(err){
                    console.log(err);
                    $scope.warning =  manageError(err);
                    $scope.authUser.password = "";
                });
            };

            this.HTTPBasicAuth = function(){ 
                $http.post("http://localhost:3000/auth", null, {
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    withCredentials : true
                })
                    .then(function(r){
                    console.log(r);
                    var data = r.data
                    $rootScope.basic = $base64.encode(data[0] + ":" + data[1]);
                }
                          ,function(err){
                    console.log(err);
                    $scope.warning =  manageError(err);
                });

            }

            this.logout = function(){ 
                $rootScope.basic = "";
                $http.post("http://localhost:3000/logout", null, {
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    withCredentials : true
                })
                    .then(function(r){
                    console.log(r);
                    var data = r.data

                    }
                          ,function(err){
                    console.log(err);
                    $scope.warning =  manageError(err);
                });
            }
        }]);

    function manageError(err){
        if (err.status == 400) {
            return "Code" + err.data.code + " Message" + err.data.message;
        }
        return "Status:" + err.status + " Message: " + err.data;
    }
})();
