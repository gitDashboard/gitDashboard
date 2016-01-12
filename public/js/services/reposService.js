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
			branch = "refs/heads/master";
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

	function fileContent(repoId,fileRef){
		var respDef = $q.defer();
		$http.get("api/v1/repo/"+repoId+"/file/"+fileRef).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;	
	}

	function createFolder(path){
		var req={
			'path':path
		}
		var respDef = $q.defer();
		$http.put("api/v1/admin/repo/mkdir",req).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;	
	}

	function createRepo(path,description){
		var req={
			'path':path,
			'description':description
		}
		var respDef = $q.defer();
		$http.put("api/v1/admin/repo/create",req).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;	
	}
	
	function initRepo(path){
		var respDef = $q.defer();
		$http.put("api/v1/admin/repo/init?path="+path).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;	
	}

	function updateDescription(repoId,description){
		var respDef = $q.defer();
		$http.post("api/v1/admin/repo/"+repoId+"/description?description="+description).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;	
	}

	function permissions(repoId){
		var respDef = $q.defer();
		$http.get("api/v1/admin/repo/"+repoId+"/permissions").success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}

	function commit(repoId,commitId){
		var respDef = $q.defer();
		$http.get("api/v1/repo/"+repoId+"/commit/"+commitId).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}

	function updatePermissions(repoId,permissions){
		var respDef = $q.defer();
		var req={
			permissions:permissions
		};
		$http.post("api/v1/admin/repo/"+repoId+"/permissions",req).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}

	return {
		"list":list,
		"commits":commits,
		"commit":commit,
		"info":info,
		"files":files,
		"fileContent":fileContent,
		'createFolder':createFolder,
		'createRepo':createRepo,
		'initRepo':initRepo,
		'permissions':permissions,
		'updatePermissions':updatePermissions,
		'updateDescription':updateDescription

	};
}]);