var groupService = angular.module('groupService', []);

groupService.factory('Group', ['$q','$http',function ($q,$http) {
	function search(query){
		var respDef = $q.defer();
		$http.get("api/v1/admin/group/search?query="+query).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	function list(){
		var respDef = $q.defer();
		$http.get("api/v1/admin/group/list").success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	function save(group){
		var respDef = $q.defer();
		$http.post("api/v1/admin/group/save",group).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	function deleteGroup(group){
		var respDef = $q.defer();
		$http.delete("api/v1/admin/group/"+group.id).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	return {
		'search':search,
		'list':list,
		'save':save,
		'deleteGroup':deleteGroup
	}
}]);