PACK_SELECTOR = document.getElementById("pack-select");
RESULT_LIST = document.getElementById("result-list");
EXPR_QUERY = document.getElementById("expr-query");


clearResults = function() {
    list = RESULT_LIST;
    while(list.lastElementChild.className != 'collection-header') {
        list.lastElementChild.remove();
    }
}

count = 0;

fetchPacks = function() {
    fetch("../pack")
        .then(response => response.json())
        .then(list => {
            // Alpha sort by title.
            list.sort(function(a,b){
                if(a.title > b.title) {
                    return -1;
                }
                if(a.title < b.title) {
                    return 1;
                }
                return 0;
            });
            selector = PACK_SELECTOR;
            while(selector.firstElementChild != selector.lastElementChild) {
                selector.lastElementChild.remove();
            }
            for (var i = 0; i < list.length; i++) {
                item = document.createElement("option");
                item.value = list[i].id;
                item.append(list[i].title);
                selector.lastElementChild.after(item);
            }
            // Have to re-init the selector.
            M.FormSelect.init(selector);
        });
}

packValue = function() {
    // FormSelect.getInstance(el).getSelectedValues doesn't seem to worw.
    // Single select was always returning the input before last, Wat?!
    // This is simpler and works...
    return PACK_SELECTOR.value;
}

createResultRow = function(result, expr, compErr, runErr, rolls) {
    next = document.createElement('li');
    header = document.createElement('div');
    header.className = "collapsible-header";
    if(compErr || runErr) {
        if(compErr) {
            header.append("Compiler Error: " + compErr)
        } else if(runErr) {
            header.append("Runtime Error: " + runErr)
        }
        console.log("setting badge");
        badge = document.createElement('span');
        badge.classList.add("badge");
        badge.classList.add("red");
        badge.classList.add("white-text");
        badge.append("Error");
        header.append(badge);
    } else {
        header.append(result);
    }
    next.append(header);
    txt = document.createElement('span');
    txt.append(expr);
    liBody = document.createElement('li');
    liBody.className = "collapsible-body";
    liBody.append(txt);
    next.append(liBody);

    return next;
}

getInputField = function(item) {
    while(!item.classList.contains("input-field")) {
        item = item.parentElement;
    }
    return item;
}

removeError = function() {
    getInputField(EXPR_QUERY).classList.remove("transition-error");
    getInputField(PACK_SELECTOR).classList.remove("transition-error");
}

setPackError = function() {
    getInputField(PACK_SELECTOR).classList.add("form-error");
    setTimeout(function(){getInputField(PACK_SELECTOR).classList.add("transition-error")}, 20);
    setTimeout(function(){getInputField(PACK_SELECTOR).classList.remove("form-error")}, 50);
}

setExprError = function() {
    getInputField(EXPR_QUERY).classList.add("form-error");
    setTimeout(function(){getInputField(EXPR_QUERY).classList.add("transition-error")}, 20);
    setTimeout(function(){getInputField(EXPR_QUERY).classList.remove("form-error")}, 50);
}

validateSubmit = function() {
    removeError();

    pack = packValue();
    err = false;
    if(pack === '') {
        err = true;
        setPackError();
    }
    rawExpr = EXPR_QUERY.value;
    if(rawExpr === '') {
        err = true;
        setExprError();
    }
    if(err){
        M.toast({html:"You must specify both a pack and an expression.", duration: 3000});
    }
    return {
        "pack": packValue(),
        "expression": (rawExpr[0] === '{')? rawExpr : ("{" + rawExpr + "}"),
        "err": err
    }
}

submitExpression = function(e) {
    e.preventDefault();

    req = validateSubmit();
    if(req.err) {
        return;
    }

    fetch(
        "../eval",
        {
            method: "POST",
            body: JSON.stringify(req)
        }
    ).then(response => response.json())
        .then(function(body) {
            console.log(body);
            RESULT_LIST.firstElementChild.after(createResultRow(
                body.result,
                rawExpr,
                body["compile-error"],
                body["runtime-error"],
                []
            ));        
        })

}
