routie({
  '/': showChoices,
  '/likes-dislikes': visLikesAndDislikes,
  '/individual-work': visIndividualWork
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



function visIndividualWork() {
  hideChoices();
  $("#visualizationTitle").html("By Individual Work");
  $("#visualizationData").html("<select id='visualizationSelect' placeholder='Find a work...'></select>");
  $("#visualizationSelect").selectize();
  $.ajax({
    url: "/visualize/work-list",
    success: function(data) {
      
    },
    dataType: "json"
  })
}
