var authService = angular.module('authService', [])

authService.factory('Auth', ['$http','$q',function ($http,$q) {
	function login(username,password,type){
		var defResponse = $q.defer();
		var user = {
			"username":username,
			"password":password,
			"type":type
		};
		$http.post("api/v1/auth/login",user).success(function(data){
			defResponse.resolve(data);
		}).error(function(data,status){
			defResponse.reject("status:"+status+" data:"+data);
		});
		return defResponse.promise;
	};
	return {
		'login':login
	};
}]);