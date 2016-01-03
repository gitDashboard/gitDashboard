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
	function info(repoId){
		var respDef = $q.defer();
		$http.get("api/v1/repo/"+repoId+"/info").success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}

	function commits(repoId,branch,start,count){
		if (branch==undefined){
			branch = "master";
		}
		var req ={
			"repoId":repoId,
			"branch":branch,
			"start":start,
			"count":count
		};
		console.log(req);
		var respDef = $q.defer();
		$http.post("api/v1/repo/"+repoId+"/commits",req).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	};

	function files(repoId,refName,parent){
		if (refName==undefined){
			refName = "master";
		}
		if (parent==undefined){
			parent=""
		}
		var req ={
			"repoId":repoId,
			"refName":refName,
			"parent":parent,
		};
		console.log(req);
		var respDef = $q.defer();
		$http.post("api/v1/repo/"+repoId+"/files",req).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	};

	return {
		"list":list,
		"commits":commits,
		"info":info,
		"files":files
	};
}]);