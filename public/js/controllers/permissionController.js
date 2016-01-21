gitDashboard.controller('PermissionController',['$scope','$location','Repo','Folder','$routeParams',function($scope,$location,Repo,Folder,$routeParams){
	var repoId =null;
	if ($routeParams.repoId!=undefined){
		repoId = parseInt($routeParams.repoId);
		Repo.info(repoId).then(function(data){
			$scope.info=data.info;
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	}
	var folderId =null;
	if ($routeParams.folderId!=undefined){
		folderId = parseInt($routeParams.folderId);
		Folder.getFolder(folderId).then(function(data){
			$scope.info=data.folder
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	}
	$scope.selUser=function(permission){
		$scope.currPerm=permission;
		$('#searchUserPopup').modal('show');		
	}

	$scope.getPermissions=function(){
		if (folderId!=null){
			Folder.permissions(folderId).then(function(data){
				if (data.success){
					$scope.permissions = data.permissions;
				}else{
					alert(data.error.message);
				}
			},function(error){
				console.log(error);
				if (error.status==401){
					$location.path("login");
				}
			});
		}
		if (repoId!=null){
			Repo.permissions(repoId).then(function(data){
				if (data.success){
					$scope.permissions = data.permissions;
				}else{
					alert(data.error.message);
				}
			},function(error){
				console.log(error);
				if (error.status==401){
					$location.path("login");
				}
			});
		}
	}
	$scope.addPermission=function(){
		$scope.permissions.push({
			users:[],
			types:[],
			ref:""
		});
	}
	$scope.isGranted=function(permission,type){
		return permission.types.indexOf(type)>-1;
	}

	$scope.grant = function(permission,type){
		var currTypePos = permission.types.indexOf(type);
		if (currTypePos==-1){
			permission.types.push(type);
		}else{
			permission.types.splice(currTypePos,1);
		}
	}
	$scope.removeUser = function(permission,userIndex){
		permission.users.splice(userIndex,1);
	}

	$scope.removePermission=function(pos){
		$scope.permissions.splice(pos,1);
	}

	$scope.updatePermissions=function(){
		if (folderId!=null){
			Folder.updatePermissions(folderId,$scope.permissions);
		}
		if (repoId!=null){
			Repo.updatePermissions(repoId,$scope.permissions);
		}
	}
}]);