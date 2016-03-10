	var currTooltip=null;
	var currTooltipText=null;
	var currGravatar=null;
	var currAuthor = null;
	var currDate = null;

	function createTooltip(stage, x,y,message,emailHash,author,commitDate){
		var length = Math.max(message.length,author.length)*9+55;
		var rect = new createjs.Shape();
		rect.graphics.setStrokeStyle(1).beginStroke("Black").beginFill("LightBlue").drawRect(x+10, y+10, length,60);
		stage.addChild(rect);

		gravatarUrl = "http://www.gravatar.com/avatar/"+emailHash+"?s=50";
		var gravatar = new createjs.Bitmap(gravatarUrl);
 		gravatar.x=x+15;
 		gravatar.y=y+15;
 		stage.addChild(gravatar);

		var author = new createjs.Text(author, "13px Mono", "black");
		author.x = x+68;
		author.y = y+40;
		author.textBaseline = "alphabetic";
		stage.addChild(author);

		var date = new createjs.Text(commitDate, "13px Mono", "black");
		date.x = x+68;
		date.y = y+55;
		date.textBaseline = "alphabetic";
		stage.addChild(date);
		
		var text = new createjs.Text(message, "13px Mono", "black");
		text.x = x+68;
		text.y = y+25;
		text.textBaseline = "alphabetic";
		stage.addChild(text); 

		stage.update();
		//to remove
		currTooltip=rect;
		currTooltipText=text;
		currGravatar=gravatar;
		currAuthor = author;
		currDate= date;
	}

	function createCommit(stage,id,x,y,message,emailHash,author,commitDate){
		var circle = new createjs.Shape();
		circle.graphics.setStrokeStyle(1).beginStroke("Black").beginFill("DeepSkyBlue").drawCircle(0, 0, 6);
		circle.x = x;
		circle.y = y;
		circle.addEventListener('mouseover',function (event){
			circle.graphics.clear().setStrokeStyle(1).beginStroke("Black").beginFill("Red").drawCircle(0, 0, 6);
			stage.update();
			createTooltip(stage,x,y,message,emailHash,author,commitDate);
		});

		circle.addEventListener('click',function (event){
			stage.scope.selCommit(id);
		});

		circle.addEventListener('mouseout',function (event){
			circle.graphics.clear().setStrokeStyle(1).beginStroke("Black").beginFill("DeepSkyBlue").drawCircle(0, 0, 6);
			stage.removeChild(currTooltip);
			stage.removeChild(currTooltipText);
			stage.removeChild(currGravatar);
			stage.removeChild(currAuthor);
			stage.removeChild(currDate);
			stage.update();
		});

		stage.addChild(circle);
		return {
			position:{
				x:x,
				y:y
			},
			message:message
		}
	}

	function connectCommit(stage,cmtSrc, cmtDst){
		var line = new createjs.Shape();
		line.graphics.setStrokeStyle(2);
		if (cmtSrc.position.y == cmtDst.position.y){
			var width = cmtDst.position.x-cmtSrc.position.x;
			line.graphics.beginStroke("Black");
			line.graphics.moveTo(cmtSrc.position.x,cmtSrc.position.y);
			line.graphics.lineTo(cmtDst.position.x,cmtDst.position.y);
		}else{
			var deltaY = cmtDst.position.y-cmtSrc.position.y;
			var deltaX = cmtDst.position.x-cmtSrc.position.x;				
			if (cmtSrc.position.y < cmtDst.position.y){
				line.graphics.beginStroke("Blue");
				line.graphics.moveTo(cmtSrc.position.x,cmtSrc.position.y);
				line.graphics.lineTo(cmtSrc.position.x,cmtDst.position.y-15);
				line.graphics.lineTo(cmtSrc.position.x+deltaX-10,cmtDst.position.y-15);
				line.graphics.lineTo(cmtDst.position.x,cmtDst.position.y);
			}
			if (cmtSrc.position.y > cmtDst.position.y){
				line.graphics.beginStroke("Red");
				line.graphics.moveTo(cmtSrc.position.x,cmtSrc.position.y);
				line.graphics.lineTo(cmtSrc.position.x+deltaX/2,cmtSrc.position.y);
				line.graphics.lineTo(cmtSrc.position.x+deltaX/2,cmtSrc.position.y+deltaY/2);
				line.graphics.lineTo(cmtDst.position.x,cmtDst.position.y);
			}
		}
		stage.addChildAt(line,0);
	}

	function drawDate(stage,y,x1,x2,color,date){
		var rect = new createjs.Shape();
		rect.graphics.setStrokeStyle(1).beginStroke("White").beginFill(color).drawRect(x1, y, x2-x1,20);
		stage.addChild(rect);
		var text = new createjs.Text(date, "12px Mono", "white");
		text.x = x1+((x2-x1)/2)-7;
		text.y = y+14;
		text.textBaseline = "alphabetic";
		stage.addChild(text);
	}

	function graphEnd(){
		$('#graphTab').scrollLeft($('#graphCanvas').attr('width')+100);
	}
	function graphBegin(){
		$('#graphTab').scrollLeft(0);
	}
	function graphForward(){
		var to = $('#graphTab').scrollLeft()+500;
		if (to< $('#graphCanvas').attr('width')+100 ){
			$('#graphTab').scrollLeft(to);
		}
	}
	function graphBackward(){
		var to = $('#graphTab').scrollLeft()-500;
		if (to<0) {
			to=0;
		}
		$('#graphTab').scrollLeft(to);
	}
