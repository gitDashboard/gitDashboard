var userService = angular.module('userService', []);

userService.factory('User', ['$q','$http',function ($q,$http) {
	function search(query){
		var respDef = $q.defer();
		$http.get("api/v1/admin/user/search?query="+query).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	function list(){
		var respDef = $q.defer();
		$http.get("api/v1/admin/user/list").success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	return {
		'search':search,
		'list':list
	}
}]);