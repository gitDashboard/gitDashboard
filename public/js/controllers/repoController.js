gitDashboard.controller('RepoController',['$scope','$routeParams','Repo','$location',function($scope,$routeParams,Repo,$location){
	var repoId = parseInt($routeParams.repoId);
	$scope.page = 1;
	$scope.count = 10;
	$scope.info = function(){
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
	}

	if ($routeParams.fileId!=null && $routeParams.fileName!=null){
		$scope.openFileContent($routeParams.fileId,$routeParams.fileName);
	}

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

	$scope.openFileContent = function(fileId,name){
		$scope.setCurrView("public/fragment/repo/file.html");
		Repo.fileContent(repoId,fileId).then(function(data){
			console.log(data);
			if (data.success){
				$scope.file={
					content:atob(data.content),
					name:name
				};
				var fileEditor = CodeMirror.fromTextArea(document.getElementById("fileContentTxtArea"),{
					lineNumbers: true,
					matchBrackets: true,
					readOnly:true,
					extraKeys: {
						"F11": function(cm) {
							cm.setOption("fullScreen", !cm.getOption("fullScreen"));
						},
						"Esc": function(cm) {
							if (cm.getOption("fullScreen")) cm.setOption("fullScreen", false);
						}
					}
				});
				fileEditor.setSize("100%",950);
				fileEditor.setOption("value",$scope.file.content);
				if ($scope.file.name.indexOf("xml")>-1){
					fileEditor.setOption("mode","text/xml");
				}
				if ($scope.file.name.indexOf("java")>-1){
					fileEditor.setOption("mode","text/x-java");
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

	$scope.openFile=function(file){
		if (file.isDir){
			$scope.getFiles(file);
			$scope.file=null;
			$scope.setCurrView(null);
		}else{
			$scope.openFileContent(file.id,file.name);
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

	$scope.updateDescription=function(){
		Repo.updateDescription(repoId,$scope.repo.description).then(function(data){
			if (!data.success){
				alert(data.error.message);
			}
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

	$scope.selGroup=function(permission){
		$scope.currPerm=permission;
		$('#searchGroupPopup').modal('show');		
	}

	$scope.updatePermissions=function(){
		Repo.updatePermissions(repoId,$scope.permissions);
	}

	$scope.setCurrView=function(view){
		$scope.viewHistory.push($scope.currView);
		$scope.currView = view;
	}

	$scope.backView = function(){
		prevView = $scope.viewHistory.pop();
		$scope.currView = prevView;
	}

	$scope.getCommit=function(commitId){
		Repo.commit(repoId,commitId).then(function(data){
			console.log(data);
			if (data.success){
				$scope.currCommit = data.commit;
				$scope.currCommit.files = data.files
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

	$scope.selCommit=function(commit){
		$scope.getCommit(commit.id)
		$scope.setCurrView('public/fragment/repo/commit.html');
	}

	$scope.lockRepo=function(repo){
		Repo.lock(repo,!repo.locked).then(function(data){
			console.log(data);
			if (data.success){
				$scope.info();				
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

	$scope.renameRepo=function(repo){
		Repo.moveRepo(repo,"").then(function(data){
			console.log(data);
			if (data.success){
				$scope.info();		
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
	$scope.info();
	$scope.viewHistory=[];
	$scope.showFile=false;
	$scope.currRef="refs/heads/master";	
	$scope.getCommits();
	$scope.getFiles(null);
}]);	