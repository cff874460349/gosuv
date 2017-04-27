/* index.js */
var W = {};

var testPrograms = [{
  program: {
    name: "gggg",
    command: "",
    dir: "",
    autoStart: true,
  },
  status: "running",
}];

var vm = new Vue({
  el: "#app",
  data: {
    isConnectionAlive: true,
    log: {
      content: '',
      follow: true,
      line_count: 0,
    },
    programs: [],
	config: {},
  },
  methods: {
    addNewProgram: function() {
      console.log("Add")
      var form = $("#formNewProgram");
      form.submit(function(e) {
        console.log("HellO")
        e.preventDefault();
        console.log(e);
        $("#newTest").modal('hide')
        return false;
      });
      // console.log(this.program.name);
    },
    updateBreadcrumb: function() {
      var pathname = decodeURI(location.pathname || "/");
      var parts = pathname.split('/');
      this.breadcrumb = [];
      if (pathname == "/") {
        return this.breadcrumb;
      }
      var i = 2;
      for (; i <= parts.length; i += 1) {
        var name = parts[i - 1];
        var path = parts.slice(0, i).join('/');
        this.breadcrumb.push({
          name: name + (i == parts.length ? ' /' : ''),
          path: path
        })
      }
      return this.breadcrumb;
    },
    refresh: function() {
      // ws.send("Hello")
      $.ajax({
        url: "/api/programs",
        success: function(data) {
          vm.programs = data;
          Vue.nextTick(function() {
            $('[data-toggle="tooltip"]').tooltip()
          })
        }
      });
    },
    reload: function() {
      $.ajax({
        url: "/api/reload",
        method: "POST",
        success: function(data) {
          if (data.status == 0) {
            alert("reload success");
          } else {
            alert(data.value);
          }
        }
      });
    },
    test: function() {
      console.log("test");
    },
    cmdStart: function(name) {
      console.log(name);
      $.ajax({
        url: "/api/programs/" + name + "/start",
        method: 'post',
        success: function(data) {
          console.log(data);
        }
      })
    },
    cmdStop: function(name) {
      $.ajax({
        url: "/api/programs/" + name + "/stop",
        method: 'post',
        success: function(data) {
          console.log(data);
        }
      })
    },
    cmdTail: function(name) {
      var that = this;
      if (W.wsLog) {
        W.wsLog.close()
      }
      W.wsLog = newWebsocket("/ws/logs/" + name, {
        onopen: function(evt) {
          that.log.content = "";
          that.log.line_count = 0;
        },
        onmessage: function(evt) {
          // strip ansi color
          // console.log("DT:", evt.data)
          that.log.content += evt.data.replace(/\033\[[0-9;]*m/g, "");
          that.log.line_count = $.trim(that.log.content).split(/\r\n|\r|\n/).length;
          if (that.log.follow) {
            var pre = $(".realtime-log")[0];
            setTimeout(function() {
              pre.scrollTop = pre.scrollHeight - pre.clientHeight;
            }, 1);
          }
        }
      });
      this.log.follow = true;
      $("#modalTailf").modal({
        show: true,
        keyboard: true,
        // keyboard: false,
        // backdrop: 'static',
      })
    },
	cmdEdit: function(command) {
      //console.log(command);
	  console.log(JSON.parse(command));
	  this.config = JSON.parse(command);
	  if (this.config.kyee.kyee == '1')
	  {
		  this.config.kyee.kyee = true;
		  $("#fields_kyee").show();
	  }
	  else
	  {
		  this.config.kyee.kyee = false;
		  $("#fields_kyee").hide();
	  }
	  if (this.config.ladrip.ladrip == '1')
	  {
	      this.config.ladrip.ladrip = true;
		  $("#fields_lianxin").show();
	  }
	  else
	  {
		this.config.ladrip.ladrip = false;
		$("#fields_lianxin").hide();
	  }
	  if (this.config.ewell.ewell == '1')
	  {
	      this.config.ewell.ewell = true;
		  $("#fields_ewell").show();
	  }
	  else
	  {
		this.config.ewell.ewell = false;
		$("#fields_ewell").hide();
	  }
	  if (this.config.kyee.data_hb == '1')
	      this.config.kyee.data_hb = true;
	  else
		this.config.kyee.data_hb = false;
	  if (this.config.kyee.data_off == '1')
	      this.config.kyee.data_off = true;
	  else
		this.config.kyee.data_off = false;
	  if (this.config.kyee.data_out == '1')
	      this.config.kyee.data_out = true;
	  else
		this.config.kyee.data_out = false;
	  if (this.config.ladrip.data_weight == '1')
	      this.config.ladrip.data_weight = true;
	  else
		this.config.ladrip.data_weight = false;
	  if (this.config.ladrip.data_e2 == '1')
	      this.config.ladrip.data_e2 = true;
	  else
		this.config.ladrip.data_e2 = false;
	  if (this.config.ewell.data_dat == '1')
	      this.config.ewell.data_dat = true;
	  else
		this.config.ewell.data_dat = false;
	  if (this.config.ewell.data_sts == '1')
	      this.config.ewell.data_sts = true;
	  else
		this.config.ewell.data_sts = false;
	  
      $("#newTest").modal({
        show: true,
        keyboard: true,
        // keyboard: false,
        // backdrop: 'static',
      })
    },
    cmdDelete: function(name) {
      if (!confirm("Confirm delete \"" + name + "\"")) {
        return
      }
      $.ajax({
        url: "/api/programs/" + name,
        method: 'delete',
        success: function(data) {
          console.log(data);
        }
      })
    },
    canStop: function(status) {
      switch (status) {
        case "running":
        case "retry wait":
          return true;
      }
    },
  },
  
})

Vue.filter('fromNow', function(value) {
  return moment(value).fromNow();
})

Vue.filter('formatBytes', function(value) {
  var bytes = parseFloat(value);
  if (bytes < 0) return "-";
  else if (bytes < 1024) return bytes + " B";
  else if (bytes < 1048576) return (bytes / 1024).toFixed(0) + " KB";
  else if (bytes < 1073741824) return (bytes / 1048576).toFixed(1) + " MB";
  else return (bytes / 1073741824).toFixed(1) + " GB";
})

Vue.filter('colorStatus', function(value) {
  var makeColorText = function(text, color) {
    return "<span class='status' style='background-color:" + color + "'>" + text + "</span>";
  }
  switch (value) {
    case "stopping":
      return makeColorText(value, "#996633");
    case "running":
      return makeColorText(value, "green");
    case "fatal":
      return makeColorText(value, "red");
    default:
      return makeColorText(value, "gray");
  }
  return value;
})

Vue.directive('disable', function(value) {
  this.el.disabled = !!value
})

$(function() {
  vm.refresh();

  $("#formNewProgram").submit(function(e) {
    var url = "/api/programs",
      data = $(this).serialize();
    $.ajax({
      type: "POST",
      url: url,
      data: data,
      success: function(data) {
        if (data.status === 0) {
          $("#newTest").modal('hide');
        } else {
          window.alert(data.error);
        }
      },
      error: function(err) {
        console.log(err.responseText);
      }
    })
    e.preventDefault()
  });


  console.log("HEE")

  function newEventWatcher() {
    W.events = newWebsocket("/ws/events", {
      onopen: function(evt) {
        vm.isConnectionAlive = true;
      },
      onmessage: function(evt) {
        console.log("response:" + evt.data);
        vm.refresh();
      },
      onclose: function(evt) {
        W.events = null;
        vm.isConnectionAlive = false;
        console.log("Reconnect after 3s")
        setTimeout(newEventWatcher, 3000)
      }
    });
  };

  newEventWatcher();

  // cancel follow log if people want to see the original data
  $(".realtime-log").bind('mousewheel', function(evt) {
    if (evt.originalEvent.wheelDelta >= 0) {
      vm.log.follow = false;
    }
  })
  $('#modalTailf').on('hidden.bs.modal', function() {
    // do something…
    console.log("Hiddeen")
    if (W.wsLog) {
      console.log("wsLog closed")
      W.wsLog.close()
    }
  })
});
