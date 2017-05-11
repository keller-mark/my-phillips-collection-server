var workList = [];
var $select;

routie({
  '/': showChoices,
  '/likes-dislikes': visLikesAndDislikes,
  '/individual-work': visIndividualWork,
  '/individual-work/?:id': visIndividualWork,
});

$(document).ready(function() {
  routie('/');

});


function hideChoices() {
  $("#visualizationChoices").hide();
}

function showChoices() {
  $("#visualizationChoices").show();
  $("#visualizationTitle").html("");
  $("#visualizationData").html("");
}


function visLikesAndDislikes() {

  console.log("Showing likes and dislikes");
  hideChoices();


  $("#visualizationTitle").html("By Likes and Dislikes");
  
  $.ajax({
    url: "/visualize/likes",
    success: function(data) {
      console.log(data);
      var numWorks = Object.keys(data).length;
      var tableDataSet = new Array(numWorks);
      var i = 0;
      $.each(data, function(key, value) {
	console.log(key);
	console.log(value);
	console.log(i);
	tableDataSet[i] = new Array(5);
	tableDataSet[i][0] = "" + value["Likes"];
	tableDataSet[i][1] = "" + value["Dislikes"];
	tableDataSet[i][2] = "" + (value["Likes"] - value["Dislikes"]);
	tableDataSet[i][3] = value["TheWork"]["Title"];
	tableDataSet[i][4] = value["TheWork"]["Maker"];
	tableDataSet[i][5] = value["TheWork"]["Year"];

	i++;
      });

      $("#visualizationData").html("<table id='visualizationTable' class='table table-striped' width='100%'></table>");
      $("#visualizationTable").DataTable({
	data: tableDataSet,
	columns: [
	  { title: "Likes"},
	  { title: "Dislikes"},
	  { title: "Net" },
	  { title: "Title" },
	  { title: "Maker"},
	  { title: "Year" }
	],
	order: [[2, "desc"]]
      });  
    },
    dataType: "json"
  });
}


function visIndividualWork(id) {
  hideChoices();
  $("#visualizationTitle").html("By Individual Work");
  if(workList.length == 0) {
    (function() {
      $.ajax({
	url: "/visualize/work-list",
	success: function(data) {
	  workList = data;
	  visIndividualWork(id);
	},
	dataType: "json"
      });
    })();
    return;
  } else {

  $("#visualizationData").html("<select id='visualizationSelect' placeholder='Find a work...'></select>");
  $select = $("#visualizationSelect").selectize({
      valueField: 'ID',
      labelField: 'Title',
      searchField: 'Title',
      options: workList,
      onChange: function(value) {
	if(!value.length) return;
	routie("/individual-work/" + value);
      }
  });
  var selectize = $select[0].selectize;

  if(typeof id !== 'undefined') {
    console.log("id is " + id);
    selectize.setValue(id);
    loadWork(id);
  }
  }
}

function loadWork(id) {
    $.ajax({
      url: "/work/" + id,
      success: function(data) {
	
	if(!$("#image" + id).length) {
	  $("#visualizationData").append("<img id='image" + id + "' class='individual-work-img col-xs-6' src='http://www.phillipscollection.org/willo/w/size3/" + data["PhillipsID"] + "w.jpg'>");
	  
	  $("#visualizationData").append("<table id='statTable" + id + "' class='col-xs-6'></table>");
	  loadWorkStats(id);
	  $("#visualizationData").append("<table id='infoTable" + id + "' class='table table-striped'></table>");
	  $.each(data, function(key, value) {
	    if(key != "ID" && key != "PhillipsID" && key != "CreatedAt" && key != "UpdatedAt" && key != "DeletedAt") {
	      key = key.replace(/([A-Z])/g, ' $1').trim()
	      $("#infoTable" + id).append("<tr><th>" + key + "</th><td>" + value + "</td></tr>");
	    }
	  });

	}
      }
    });
}

var ageStrings = {
  "1" : "Child (0-12)",
  "2" : "Teen (13-17)",
  "3" : "Young Adult (18-30)",
  "4" : "Adult (30-60)",
  "5" : "Senior (60+)"
};
var genderStrings = {
  "1" : "Male",
  "2" : "Female",
  "3" : "Other"
};

function loadWorkStats(id) {
  $.ajax({
    url: '/visualize/work/' + id,
    success: function(data) {
      console.log(data);
      $.each(data, function(key, value) {
  	newKey = key.replace(/([A-Z])/g, ' $1').trim();
	if(key == "Likes" || key == "Dislikes" || key == "Net") {
	  $("#statTable" + id).append("<tr><th>" + key + "</th><td>" + value + "</td></tr>");
	} else {
	  $("#statTable" + id).append("<tr><th>" + newKey + "</th><td><table id='" + key + id + "' class='innertable'></table></td></tr>");
	  $.each(value, function(group, num) {
	    if(num != 0) {
	      if(key == "LikesByAge" || key == "DislikesByAge" || key == "NetByAge") {
		$("#" + key + id).append("<tr><th>" + ageStrings[group] + "</th><td>" + num + "</td></tr>");
	      } else if(key == "LikesByGender" || key == "DislikesByGender" || key == "NetByGender") {
		$("#" + key + id).append("<tr><th>" + genderStrings[group] + "</th><td>" + num + "</td></tr>");
	      } else {
		$("#" + key + id).append("<tr><th>" + group + "</th><td>" + num + "</td></tr>");
	      }
	    }
	  });
	}
      });

    },
    dataType: "json"
  });
}
