var folderService = angular.module('folderService', [])

folderService.factory('Folder', ['$http','$q',function ($http,$q) {
	function list(parentId){
		if (parentId==null){
			parentId=0;
		}
		var defResponse = $q.defer();
		$http.get("api/v1/folders/"+parentId+"/list").success(function(data){
			defResponse.resolve(data);
		}).error(function(data,status){
			defResponse.reject("status:"+status+" data:"+data);
		});
		return defResponse.promise;
	};
	function createFolder(parentId,name,description){
		var req={
			'name':name,
			'description':description
		}
		var respDef = $q.defer();
		$http.put("api/v1/admin/folder/"+parentId+"/mkdir",req).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;	
	};
	function getFolder(folderId){
		var respDef = $q.defer();
		$http.get("api/v1/folder/"+folderId).success(function (data){
			respDef.resolve(data);
		}).error(function (data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;	
	}
	function permissions(folderId){
		var respDef = $q.defer();
		$http.get("api/v1/admin/folder/"+folderId+"/permissions").success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	function updatePermissions(folderId,permissions){
		var respDef = $q.defer();
		var req={
			permissions:permissions
		};
		$http.post("api/v1/admin/folder/"+folderId+"/permissions",req).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	function setAdmins(folderId,admins){
		var respDef = $q.defer();
		var req={
			admins:admins
		};
		$http.post("api/v1/admin/folder/"+folderId+"/admins",req).success(function (data){
			respDef.resolve(data);
		}).error(function(data,status){
			respDef.reject({"error":data,"status":status} );
		});
		return respDef.promise;
	}
	return {
		'list':list,
		'createFolder':createFolder,
		'getFolder':getFolder,
		'permissions':permissions,
		'updatePermissions':updatePermissions,
		'setAdmins':setAdmins
	};
}]);