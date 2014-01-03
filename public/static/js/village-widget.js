/* troll farm widget file */


var WEBSOCKET_HOST = "ws://127.0.0.1:5000";
//var WEBSOCKET_HOST = "ws://troll-farm.herokuapp.com";

var DOMAIN = "http://127.0.0.1:5000";
//var DOMAIN = "https://troll-farm.herokuapp.com";
var NEW_CONNECTION_ENDPOINT = "/connect";

var TROLL_VILLAGE_WIDGET_ID = "troll-village-widget";
var TROLL_VILLAGE_WIDGET_CONTAINER_ID = "troll-village-widget-container";

var TROLL_SCRIPTS = 	["/static/js/village-module.js",];
var TROLL_STYLESHEETS = ["/static/css/widget.css",];


/* wrap in anonymous function as to not interfere with existing function and variable names */
var trollVillageWidgetModule = function() {
	/* replace _jQuery_src with your desired jQuery file source.  Otherwise defaults to 1.10.2 hosted by google */
	var _jQuery_src = "/static/lib/jQuery-2.0.3.min.js";
	var jQuery;

	function isEI8or9() {
		var rv = -1;
		if (navigator.appName == 'Microsoft Internet Explorer') {
			var ua = navigator.userAgent;
			var re  = new RegExp("MSIE ([0-9]{1,}[\.0-9]{0,})");
			if (re.exec(ua) != null)
				rv = parseFloat( RegExp.$1 );
		}
		if ((rv > 0) && (rv < 10)) {
			return true;
		} else {
			return false;
		}
	}

	var withScript = function(src, callback) {
		if (!src) {
			callback();
		}
		var script_tag = document.createElement('script');

		script_tag.setAttribute("type","text/javascript");
		script_tag.setAttribute("src", src);

		if (script_tag.readyState) {
			script_tag.onreadystatechange = function () { // For old versions of IE
				if (this.readyState == 'complete' || this.readyState == 'loaded') {
					callback();
				}
			};
		} else {
			script_tag.onload = callback;
		}
		// Try to find the head, otherwise default to the documentElement
		(document.getElementsByTagName("head")[0] || document.documentElement).appendChild(script_tag);
	};

	var withJQ = function(callback) {
		var jQuery_src = (_jQuery_src || "//ajax.googleapis.com/ajax/libs/jquery/2.0.3/jquery.min.js");
		if (window.jQuery === undefined) {
			withScript(jQuery_src, function() {
				jQuery = window.jQuery.noConflict(true);
				callback();
			});
		} else {
			jQuery = window.jQuery;
			callback();
		}
	};

	var withStyleSheet = function(src, callback) {
		if (document.createStyleSheet) {
			document.createStyleSheet(src);
		} else {
			var file = document.createElement("link")
			file.setAttribute("rel", "stylesheet")
			file.setAttribute("type", "text/css")
			file.setAttribute("href", src)

			if (typeof file !== "undefined")
				document.getElementsByTagName("head")[0].appendChild(file)
		}
		callback();
	};
	var loadDependencies = function(callback) {
		var numDependencies = TROLL_STYLESHEETS.length + TROLL_SCRIPTS.length + 1; // +1 for jquery
        var numLoaded = 0;
        function dependencyLoaded() {
                numLoaded++;
                if (numLoaded === numDependencies) {
                        callback();
                }
        }
        withJQ(dependencyLoaded);
		for(var i in TROLL_STYLESHEETS) {
			withStyleSheet(DOMAIN + TROLL_STYLESHEETS[i], dependencyLoaded);
		}
		for(var i in TROLL_SCRIPTS) {
			withScript(DOMAIN + TROLL_SCRIPTS[i], dependencyLoaded);
		}
	};
	loadDependencies(main);


	/* helper for 'getting' and 'posting' (cough... jsonp hack to get around cross origin issue...) */
	function jsonp(data_url, data, onSuccess, onError){
		jQuery.ajax({
			url: DOMAIN + data_url,
			dataType: 'jsonp',
			data: data
		}).done(function(returnedData){ 
			if(returnedData.error){ 
				if (onError){ onError(returnedData); }
				else{ console.log('ERROR: ' + returnedData.error); }
			}
			else{ onSuccess(returnedData); }
		}).fail(function(jqXHR, textStatus, errorThrown) {
			console.log('.fail');
			if (onError) { onError(jqXHR); } 
			else { console.log('ERROR IN JSONP REQUEST'); }
		});
	}
	/* check if user is using mobile browser -- return true if so */
	function isMobile() {
		var check = false;
		(function(a){if(/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(a)||/1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(a.substr(0,4)))check = true})(navigator.userAgent||navigator.vendor||window.opera);
		if(window.innerWidth <= 800 && window.innerHeight <= 600) { check = true; }
		return check;
	}
	/* check if canvas is supported -- technique used by modernizr */
	function isCanvasSupported(){
		var elem = document.createElement('canvas');
		return !!(elem.getContext && elem.getContext('2d'));
	}

	function fill_widget_content(content) {
		jQuery('#' + TROLL_VILLAGE_WIDGET_ID).html(content);
	};

	function main(){

		/* do nothing if mobile */
		if (isMobile()) { return null; }
		/* do nothing if canvas not supported */
		if (!isCanvasSupported()) { console.log("Canvas not supported"); return null; }

		jQuery(document).ready(function($) {

			var widget = document.createElement('div');
			widget.id = TROLL_VILLAGE_WIDGET_ID;
			var widgetContainer = document.getElementById(TROLL_VILLAGE_WIDGET_CONTAINER_ID)
			widgetContainer.appendChild(widget);

			fill_widget_content("<h3>hi trolls...</h3>");

			this.trollVillageModule = new TrollVillageModule(widget);
		});
	}
	return {
		jsonp: jsonp,
		fill_widget_content: fill_widget_content,
	};	
}();




