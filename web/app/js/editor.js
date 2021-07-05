fd = -1

editorTextArea = document.getElementById("collaborative-text-editor")

var editor = CodeMirror.fromTextArea(editorTextArea, {
    lineNumbers: true
});

// set callbacks
editor.on("change", onChange);
editor.on("cursorActivity", onCursorActivity);

function disableEditor() {
    editor.options.readOnly = 'nocursor'
}

function enableEditor() {
    editor.options.readOnly = false
}

function documentLoaded() {
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
                    DocumentInsertAt(fd, "\n", offset)
                    offset++
                }
                for (var j = 0; j < change.text[i].length; j++) {
                    DocumentInsertAt(fd, change.text[i][j],  offset)
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
                    DocumentDeleteAt(fd, offset)
                }
                for (var j = 0; j < change.removed[i].length; j++) {
                    DocumentDeleteAt(fd, offset)
                }
            }

            break;
    }
}

function onCursorActivity(editor) {
    DocumentChangeCursor(fd, editor.getDoc().indexFromPos(editor.getCursor()))
}

function openNewDoc() {
    editor.setValue("")
    fd = DocumentNew()
    documentLoaded()
}

function saveLocal() {
    let serializedDoc = DocumentSerialize(fd)
    downloadFile(serializedDoc, "your-document.txt")
}

function closeDoc() {
    DocumentClose(fd)
    editor.setValue("")
    disableEditor()
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
        fd = DocumentDeserialize(Uint8Array.from(textContent.split(",").map(function (item) {
            return parseInt(item, 10)
        })), initCallback)
        documentLoaded()
    }
    reader.onerror = (e) => {
        const error = e.target.error
        console.error(`Error occurred while reading ${file.name}`, error)
    }
    reader.readAsText(file)
}

// initially the editor is disabled
disableEditor()