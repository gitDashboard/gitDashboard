var gitDashboard = angular.module("gitDashboard",['angular-jwt','LocalStorageModule','authService','reposService']);

gitDashboard.config(function($interpolateProvider) {
	$interpolateProvider.startSymbol('{[{');
	$interpolateProvider.endSymbol('}]}');
});


gitDashboard.config(function ($httpProvider, jwtInterceptorProvider) {
	function localStorageService(localStorageService) {
		return localStorageService.get('jwt_token');
	}

	jwtInterceptorProvider.tokenGetter = ['localStorageService', localStorageService];
	$httpProvider.interceptors.push('jwtInterceptor');
});

gitDashboard.controller('MainCtrl', ['$scope','localStorageService', function ($scope,localStorageService) {
	$scope.isLogged=function(){
		return localStorageService.get('jwt_token') != undefined;
	}
}]);

gitDashboard.controller('LoginController', ['$scope','Auth','localStorageService', function ($scope,Auth,localStorageService) {
	$scope.types = ["internal","ldap"];
	$scope.login = function(){
		Auth.login($scope.username,$scope.password,$scope.type).then(function(data){
			if (!data.success){
				alert("Login failed");
			}else{
				localStorageService.set('jwt_token',data.jwt_token);
			}
		},function(error){
			alert(error);
		});
	};

	$scope.logout = function () {
		localStorageService.remove('jwt_token');
	};

}]);

gitDashboard.controller('ReposController',function($scope){
	$scope.test="ciao";
});