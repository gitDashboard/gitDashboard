var userService = angular.module('userService', []);

userService.factory('User', ['$q','$http',function ($q,$http) {
	function search(query){
		var respDef = $q.defer();
		$http.get("api/v1/user/search?query="+query).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	function ldapSearch(username){
		var respDef = $q.defer();
		$http.get("api/v1/admin/user/ldapSearch?username="+username).success(function (data){
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
	function save(user){
		var respDef = $q.defer();
		$http.post("api/v1/admin/user/save",user).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	function deleteUser(user){
		var respDef = $q.defer();
		$http.delete("api/v1/admin/user/"+user.id).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	return {
		'search':search,
		'ldapSearch':ldapSearch,
		'list':list,
		'save':save,
		'deleteUser':deleteUser
	}
}]);