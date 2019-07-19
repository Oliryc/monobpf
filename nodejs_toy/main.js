    var fs = require('fs');
    var util = require('util');
     
    var start = process.hrtime();
     
    for (var i = 0; i < 30; ++i) {
      (function (id) {
        fs.readdir('.', function () {
          var end = process.hrtime(start);
          console.log(util.format('readdir %d finished in %ds', id, end[0] + end[1] / 1e9));
        });
      })(i);
    }
