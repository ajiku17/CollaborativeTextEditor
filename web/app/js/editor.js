docId = -1

editorTextArea = document.getElementById("collaborative-text-editor")

var editor = CodeMirror.fromTextArea(editorTextArea, {
    lineNumbers: true,
    theme: "base16-dark"
});

// set callbacks
editor.on("change", onChange);
editor.on("cursorActivity", onCursorActivity);

function disableButtons() {
    console.log("buttons disabled")
}

function enableButtons() {
    console.log("buttons enabled")
}

function disableEditor() {
    editor.options.readOnly = 'nocursor'
}

function enableEditor() {
    editor.options.readOnly = false
}

function enableLoading() {
    document.querySelectorAll('div:not(.loader)').forEach(function(currentValue, currentIndex, listObj) {
        currentValue.style.display = 'none';
    })
    document.getElementsByClassName('loader').item(0).style.display = 'block'
}

function disableLoading() {
    document.querySelectorAll('div:not(.loader)').forEach(function(currentValue, currentIndex, listObj) {
        currentValue.style.display = 'block';
    })
    document.getElementsByClassName('loader').item(0).style.display = 'none'
}

function connectionChanged(connected) {
    console.log(connected)
    if (connected) {
        console.log(document.getElementsByClassName('connection').item(0))
        document.getElementsByClassName('connection').item(0).style.backgroundColor = '#00bb00'
    } else {
        document.getElementsByClassName('connection').item(0).style.backgroundColor = '#aa1111'
    }
}

function documentLoaded(documentId) {
    docId = documentId
    pushState(docId)
    console.log("Document Loaded!")
    enableEditor()
}

function onChange(editor, change) {
    let offset;
    switch (change.origin) {
        case "+input":
            offset = editor.getDoc().indexFromPos({
                line: change.from.line,
                ch: change.from.ch,
                sticky: null
            });

            for (var i = 0; i < change.text.length; i++) {
                if (i > 0) {
                    DocumentInsertAt(docId, "\n", offset)
                    offset++
                }
                for (var j = 0; j < change.text[i].length; j++) {
                    DocumentInsertAt(docId, change.text[i][j],  offset)
                    offset++
                }
            }

            break
        case "+delete":
            offset = editor.getDoc().indexFromPos({
                line: change.from.line,
                ch: change.from.ch,
                sticky: null
            });
            for (var i = 0; i < change.removed.length; i++) {
                if (i > 0) {
                    DocumentDeleteAt(docId, offset)
                }
                for (var j = 0; j < change.removed[i].length; j++) {
                    DocumentDeleteAt(docId, offset)
                }
            }

            break;
    }
}

function onCursorActivity(editor) {
    DocumentChangeCursor(docId, editor.getDoc().indexFromPos(editor.getCursor()))
}

function onDocChange(changeName, changeObj) {
    // console.log("received change with name ", changeName, " and change obj ", changeObj)
    switch (changeObj.changeName) {
        case "insert":
            // console.log("calling insert")
            editor.getDoc().replaceRange(changeObj.value,
                editor.getDoc().posFromIndex(changeObj.index),
                editor.getDoc().posFromIndex(changeObj.index),
                "ignore")

            break

        case "delete":
            // console.log("calling delete")
            editor.getDoc().replaceRange("",
                editor.getDoc().posFromIndex(changeObj.index),
                editor.getDoc().posFromIndex(changeObj.index + 1))

            break

        case "peer_cursor":
            break
    }
}

function onPeerConnect(peerId, cursorPos) {
    console.log("peer ", peerId, " connected with cursor position ", cursorPos)
    addPeer(peerId)
}

function addPeer(peerId) {
    var ul = document.getElementById('connected-peers');
    var li = document.createElement("li");
    li.setAttribute('id', peerId)
    var t = document.createTextNode(peerId + " Connected");
    li.appendChild(t);
    ul.appendChild(li);
}

function onPeerDisconnect(peerId) {
    console.log("peer ", peerId, " disconnected")
    removePeer(peerId)
}

function removePeer(peerId) {
    var li = document.getElementById(peerId);
    if (li != null && li != undefined) {
        li.textContent = peerId + " disconnected"
        li.value = peerId + "Disconnected";
    }

}

function pushState(docId) {
    var path = '?doc=' + docId
    console.log('pushing state ' + path);
    history.pushState('', '', path);
}

function popState() {
    history.pushState('', '', "/");
}

function openNewDoc() {
    editor.setValue("")
    if (docId !== -1) {
        DocumentClose(docId)
        editor.setValue("")
    }

    DocumentNew(onDocChange, onPeerConnect, onPeerDisconnect, documentLoaded)
}

function saveLocal() {
    let serializedDoc = DocumentSerialize(docId)
    downloadFile(serializedDoc, "your-document.txt")
}

function openDocById(id) {
    if (docId !== -1) {
        DocumentClose(docId)
        editor.setValue("")
    }

    DocumentOpen(id, initCallback, onDocChange, onPeerConnect, onPeerDisconnect, documentLoaded)
}

function closeDoc() {
    DocumentClose(docId)
    editor.setValue("")
    popState()
    disableEditor()
    inputElement.value = ""
}

function disconnect() {
    connectionChanged(false)
    DocumentDisconnect(docId)
}

function reconnect() {
    connectionChanged(true)
    console.log("REconnect")
    DocumentReconnect(docId)
}

function downloadFile(data, filename) {
    const file = new Blob([data], {type: "multipart/byteranges;charset=utf-8"});

    const a = document.createElement("a")
    const url = URL.createObjectURL(file)
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()

    setTimeout(function() {
        document.body.removeChild(a)
        window.URL.revokeObjectURL(url)
    }, 0);
}

function initCallback(initialText) {
    editor.setValue(initialText)
}

const inputElement = document.getElementById("inputElement")

inputElement.onchange = (e) => {
    const file = inputElement.files[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = (e) => {
        const textContent = e.target.result
        if (docId !== -1) {
            DocumentClose(docId)
            editor.setValue("")
        }
        docId = DocumentDeserialize(Uint8Array.from(textContent.split(",").map(function (item) {
            return parseInt(item, 10)
        })), initCallback, onDocChange, onPeerConnect, onPeerDisconnect, documentLoaded)
    }
    reader.onerror = (e ) => {
        const error = e.target.error
        console.error(`Error occurred while reading ${file.name}`, error)
    }
    reader.readAsText(file)
}

function parseReq () {
    console.log("parsing request")
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString)

    console.log(urlParams)
    if (urlParams.has("doc")) {
        console.log("opening doc by id", urlParams.get("doc"))
        openDocById(urlParams.get("doc"))
    }
}

function initJS() {
    // initially the editor is disabled
    console.log("init js")

    enableLoading()
    disableButtons()
    disableEditor()

    setTimeout(function () {
        parseReq()
        disableLoading()
    }, 3000)

    console.log("disabling editor")
    enableButtons()
}

// run()
initJS()

// setTimeout(parseReq, 3000)