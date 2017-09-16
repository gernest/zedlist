function sendJSON(method, url, body) {
    let payload;
    if (body !== null) {
        payload = JSON.stringify(body);
    }
    return $.ajax({
        type: method,
        url: url,
        data: payload,
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
    })
}

module.exports.sendJSON=sendJSON;