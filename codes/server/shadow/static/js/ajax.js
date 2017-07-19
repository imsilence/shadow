function ajax(url, method, args, callback, dataType) {
    jQuery.ajax({
        type: method,
        url: url,
        data: data,
        success: function(response, status, xhr) {
            if(validate_response(response, status, xhr)) {
                callback(response, status, xhr);
            }
        },
        dataType: dataType
    });
}

function ajax_get(url, args, callback, dataType) {
    jQuery.post(url, args, function(response, status, xhr) {
        if(validate_response(response, status, xhr)) {
            callback(response, status, xhr);
        }
    }, dataType);
}

function ajax_post(url, args, callback, dataType) {
    jQuery.post(url, args, function(response, status, xhr) {
        if(validate_response(response, status, xhr)) {
            callback(response, status, xhr);
        }
    }, dataType);
}

function validate_response(response, status, xhr) {
    return true;
}
