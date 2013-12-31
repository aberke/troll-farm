
/* -------------------------- VILLAGE ---------------------- */

var TrollConnection = function() {
	var ws;


	function connect() {
		var ws = new WebSocket("ws://" + host + NEW_CONNECTION_ENDPOINT);
		return ws;
	}
	this.init = function(callback) {
		console.log("TrollConnection init")
		ws = connect();

		ws.onmessage = function(event) {
			//console.log(event);
			//console.log("onmessage: " + event.data);
		};
		ws.onopen = function(event) {
			if (callback) callback();
		}
		console.log(ws)

	}
	function send(msgBody) {
		ws.send(JSON.stringify(msgBody));
	}
	this.sendTest = function() {
		var msg = {"message-type": "test",};
		send(msg);
	}
	this.sendTrollsRequest = function() {
		var msg = {"message-type": "trolls"};
		send(msg);
	}
}

var TrollVillageModule = function(widgetDiv) {

	this.localTroll;

	this.trollConnection;

	this.widgetDiv = widgetDiv;
	this.canvas;
	this.context;

	/* board is a 10x10 grid with 40x40 px squares  */
	this.board = {"width": 10,
				  "height": 10,
				  "cellSize": 40,
				 }

	this.createCanvas = function() {
		this.canvas = document.createElement('canvas');
		this.canvas.id = "trollVillageModule-canvas";
		this.canvas.width = this.board.width*this.board.cellSize;
		this.canvas.height = this.board.height*this.board.cellSize;
		this.widgetDiv.appendChild(this.canvas);
		this.context = this.canvas.getContext("2d");
	}
	this.drawBoard = function() {
		if (! this.canvas ) this.createCanvas();

		for (var x=0; x<=this.board.width; x+=1) {
			this.context.moveTo(x*this.board.cellSize, 0);
			this.context.lineTo(x*this.board.cellSize, this.board.height*this.board.cellSize);
		}
		for (var y=0; y<=this.board.height; y+=1) {
			this.context.moveTo(0, y*this.board.cellSize);
			this.context.lineTo(this.board.width*this.board.cellSize, y*this.board.cellSize);
		}	
		this.context.strokeStyle = "black";
		this.context.stroke();
	}

	/* Define an object to hold all our images for the game so images are only ever created once. */
	this.imageRepository = new function() {
        // Define images
        this.troll = new Image();
        this.otherTroll = new Image();

        // Ensure all images have loaded before starting the game
        var numImages = 2;
        var numLoaded = 0;
        function imageLoaded() {
            numLoaded++;
            if (numLoaded === numImages) {
                    //this.init();
            }
        }
        this.troll.onload = function() {
                imageLoaded();
        }
        this.otherTroll.onload = function() {
                imageLoaded();
        }

        // Set images src
        this.troll.src = DOMAIN + "/static/img/troll.gif";
        this.otherTroll.src = "/static/img/other-troll.gif";
	}
	this.init = function() {
		this.drawBoard();

		var trollConnection = new TrollConnection();
		trollConnection.init(function() {
			trollConnection.sendTest();
			trollConnection.sendTrollsRequest();

		});
		this.trollConnection = trollConnection;


		this.localTroll = new Troll();

		setKeyListeners(this.localTroll.move);

		this.localTroll.setImage(this.imageRepository.troll);
		this.localTroll.init(0,0, this.context, this.board);

	}

}

var Troll = function() {
	//var _x,_y,_width,_height;
	this.x;
	this.y;
	this.padding = 5;
	this.x_px;
	this.y_px;
	this.width = 20;
	this.height = 20;
	this.board;

	this.img;
	this.context;
	var self = this;

	this.move = function(direction) {
		if (direction == "left") {
			if (self.x > 0)  self.x--;
		} else if (direction == "right") {
			if (self.x < (self.board.width-1)) self.x++;
		} else if (direction == "up") {
			if (self.y > 0) self.y--;
		} else if (direction == "down") {
			if (self.y < (self.board.height-1)) self.y++;
		} else {
			console.log("direction: " + direction);
		}
		self.draw();
	}

	this.draw = function() {
		if (this.x_px)
			this.context.clearRect(self.x_px, self.y_px, self.width, self.height);

		this.x_px = this.x*this.board.cellSize + this.padding;
		this.y_px = this.y*this.board.cellSize + this.padding;
		
		this.context.drawImage(this.img, this.x_px, this.y_px, this.width, this.height);
	}
	this.setImage = function(img) {
		this.img = img;
	}
	this.update = function(x, y) {
		this.x = x;
		this.y = y;
		this.draw();
	}
	this.init = function(x, y, context, board) {
		this.x = x;
		this.y = y;
		this.context = context;
		this.board = board;
		this.draw();
	}

}

/* --------------------------------------------------------------- */

KEY_CODES = {
  32: 'space',
  37: 'left',
  38: 'up',
  39: 'right',
  40: 'down',
}
function setKeyListeners(onkeydownCall) {

	document.onkeydown = function(e) {
		// Firefox and opera use charCode instead of keyCode to return which key was pressed.
		var keyCode = (e.keyCode) ? e.keyCode : e.charCode;
		if (KEY_CODES[keyCode]) {
			e.preventDefault();
			onkeydownCall(KEY_CODES[keyCode]);
		}
	}
	document.onkeyup = function(e) {
	  var keyCode = (e.keyCode) ? e.keyCode : e.charCode;
	  if (KEY_CODES[keyCode]) {
	    e.preventDefault();
	  }
	}


}


