var reposService = angular.module('reposService', []);

reposService.factory('Repo', ['$q','$http',function ($q,$http) {
	function list(path){
		var respDef = $q.defer();
		var req ={
			"subPath":path
		};
		console.log(req);
		$http.post("api/v1/repos/list",req).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject("Error: "+data+" Status:"+status);
		});
		return respDef.promise;
	};
	function commits(repoId,branch,start,count,ascending){
		if (branch==undefined){
			branch = "master";
		}
		var req ={
			"repoId":repoId,
			"branch":branch,
			"start":start,
			"count":count,
			"ascending":ascending
		};
		console.log(req);
		var respDef = $q.defer();
		$http.post("api/v1/repo/"+repoId+"/commits",req).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject("Error: "+data+" Status:"+status);
		});
		return respDef.promise;
	};

	return {
		"list":list,
		"commits":commits
	};
}]);