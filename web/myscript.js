function validateForm() {
    for(var i =0;i<16;i++){
        var element = "" + i
        var x = document.forms["boggle"][element].value;
        if( x.length == 0 ){
            alert("Cannot have empty cells");
            return false;
        }
        var regex = /^[a-zA-Z()]*$/
        if(!x.match(regex)){
            alert("Only aphabets allowed");
            return false;
        }
        if (x.length > 1) {
            if (x.toUpperCase() != "QU"){
                alert("Boggle borad can have only one letter per block or 'Qu'");
                return false;
            }
        }
    }
    
  } 