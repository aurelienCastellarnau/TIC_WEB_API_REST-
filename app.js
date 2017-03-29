(function(){
    'use strict';

    angular.module('clientRest', [
        'base64',
        'ngCookies',
        'clientRestDirectives',
        'clientRestControllers'
    ])
        .factory('httpRequestInterceptor',
                 ['$rootScope', '$base64', function($rootScope, $base64)
                  {
                      return {
                          request: function($config) {
                              if ($rootScope.basic)
                              {
                                  $config.headers['Authorization'] = "Basic " + $rootScope.basic;
                              }
                              return $config;
                          }
                      };
                  }])
        .config(['$httpProvider', function($httpProvider){
            $httpProvider.interceptors.push('httpRequestInterceptor');
        }]);
})();
