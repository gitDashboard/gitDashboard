gitDashboard.controller('RepoController',['$scope','$routeParams','Repo','$location',function($scope,$routeParams,Repo,$location){
	var repoId = parseInt($routeParams.repoId);
	$scope.page = 1;
	$scope.count = 10;
	Repo.info(repoId).then(function(data){
		if (data.success){
			console.log(data);
			$scope.repo = data.info;
		}else{
			console.log(data.error);
			alert(data.error.message);
		}
	},function(error){
		console.log(error);
		alert(error);
	});

	$scope.returnToFolder = function(){
		$location.path("").search({path:$scope.repo.folderPath});
	}

	$scope.decPage=function(){
		if ($scope.page>1){
			$scope.page--;
		}
		$scope.getCommits();
	}
	$scope.incPage=function(){
		if ($scope.commits.length>0){
			$scope.page++;
			$scope.getCommits();
		}
	}
	$scope.getCommits=function(){
		Repo.commits(repoId,$scope.currRef,($scope.page-1)*parseInt($scope.count),parseInt($scope.count),$scope.ascending).then(function(data){
			console.log(data);
			if (data.success){
				$scope.commits = data.commits;
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	}
	$scope.getFiles=function(parent){
		if (parent==null){
			$scope.currPath=[];
			parentId=null;
		}else{
			parentId = parent.id;
		}
		$scope.currPath.push(parent);
		console.log($scope.currPath);
		Repo.files(repoId,$scope.currRef,parentId).then(function(data){
			console.log(data);
			if (data.success){
				$scope.parentTreeId = data.parentTreeId;
				$scope.files = data.files;
			}
		},function(error){
			console.log(error);
			if (error.status==401){
				$location.path("login");
			}
		});
	}
	$scope.inFolder=function(){
		return $scope.currPath.length>1;
	}
	$scope.upFilesDir=function(){
		if ($scope.currPath.length>0){
			$scope.currPath.pop();//current
			parent = $scope.currPath.pop();
			$scope.getFiles(parent);
			$scope.fileContent=null;
		}
	}
	$scope.getPath=function(){
		var strPath = "";
		for (var i=0;i<$scope.currPath.length;i++){
			if ($scope.currPath[i]!=null){
				strPath+="/"+$scope.currPath[i].name;
			}
		}
		if ($scope.file!=null){
			strPath+="/"+$scope.file.name;
		}
		return strPath;
	}

	$scope.codemirrorLoaded = function(_editor){
		console.log(_editor);
	}

	$scope.openFile=function(file){
		if (file.isDir){
			$scope.getFiles(file);
			$scope.file=null;
			$scope.showFile=false;
		}else{
			$scope.showFile=true;
			Repo.fileContent(repoId,file.id).then(function(data){
				console.log(data);
				if (data.success){
					$scope.file={
						content:atob(data.content),
						name:file.name
					};
					var codeMirrorInstance = $('.CodeMirror')[0].CodeMirror;
					codeMirrorInstance.setOption("value",$scope.file.content);
					if (file.name.indexOf("xml")>-1){
						codeMirrorInstance.setOption("mode","text/xml");
					}
					if (file.name.indexOf("java")>-1){
						codeMirrorInstance.setOption("mode","text/x-java");
					}
				}else{
					alert(data.error.message)
				}
			},function(error){
				console.log(error);
				if (error.status==401){
					$location.path("login");
				}
			});
		}
		console.log(file);
	}


	$scope.getPermissions=function(){
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
	$scope.addPermission=function(){
		$scope.permissions.push({
			userName:"",
			groupName:"",
			type:"",
			ref:""
		});
	}
	$scope.removePermission=function(pos){
		$scope.permissions.splice(pos,1);
	}
	$scope.selUser=function(permission){
		$scope.currPerm=permission;
		$('#searchUserPopup').modal('show');		
	}

	$scope.selGroup=function(permission){
		$scope.currPerm=permission;
		$('#searchGroupPopup').modal('show');		
	}

	$scope.updatePermissions=function(){
		Repo.updatePermissions(repoId,$scope.permissions);
	}

	$scope.setCurrView=function(view){
		$scope.currView = view;
	}

	$scope.showFile=false;
	$scope.currRef="refs/heads/master";	
	$scope.currPath=[];
	$scope.getCommits();
	$scope.getFiles(null);
}]);	