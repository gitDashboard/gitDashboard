var reposService = angular.module('reposService', []);

reposService.factory('Repo', ['$q','$http',function ($q,$http) {
	function list(parentFolderId){
		var respDef = $q.defer();		
		$http.get("api/v1/repos/"+parentFolderId+"/list").success(function (data){
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

	function createRepo(folderId,name,description){
		var req={
			'folderId':folderId,
			'name':name,
			'description':description
		}
		var respDef = $q.defer();
		$http.put("api/v1/admin/repo/create?folderId="+folderId,req).success(function (data){
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

	function moveRepo(repo,folderId){
		var respDef = $q.defer();
		var req={
			folderId:parseInt(folderId),
			destName:repo.name
		};
		console.log("moive");
		console.log(req);
		$http.post("api/v1/admin/repo/"+repo.id+"/move",req).success(function (data){
			data.repo = repo;
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status,"repo":repo} );
		});
		return respDef.promise;
	}

	function lockRepo(repo,lock){
		var respDef = $q.defer();
		var url = "api/v1/admin/repo/"+repo.id+"/";
		if (lock){
			url+="lock";
		}else{
			url+="unlock";
		}console.log(url);
		$http.post(url).success(function (data){
			data.repo = repo;
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status,"repo":repo} );
		});
		return respDef.promise;
	}

	function graph(repoId){
		var respDef = $q.defer();
		$http.get("api/v1/repo/"+repoId+"/graph").success(function (data){
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
		'createRepo':createRepo,
		'initRepo':initRepo,
		'permissions':permissions,
		'updatePermissions':updatePermissions,
		'updateDescription':updateDescription,
		'moveRepo':moveRepo,
		'lock':lockRepo,
		'graph':graph
	};
}]);