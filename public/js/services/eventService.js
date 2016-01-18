var eventService = angular.module('eventService', [])

eventService.factory('Event', ['$http','$q',function ($http,$q) {
	function search(searchParams){
		var defResponse = $q.defer();
		$http.post("api/v1/admin/event/search ",searchParams).success(function(data){
			defResponse.resolve(data);
		}).error(function(data,status){
			defResponse.reject("status:"+status+" data:"+data);
		});
		return defResponse.promise;
	};
	return {
		'search':search
	};
}]);