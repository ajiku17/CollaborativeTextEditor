fd = 1

editorTextArea = document.getElementById("collaborative-text-editor")

var editor = CodeMirror.fromTextArea(editorTextArea, {
    lineNumbers: false
});

// set callbacks
editor.on("beforeChange", onBeforeChange);
editor.on("change", onChange);
editor.on("cursorActivity", onCursorActivity);

function onBeforeChange(editor, change) {
    console.log("[on Before Change]\neditor:", editor, "\nchange:", change)
}

function onChange(editor, change) {
    console.log("[on Change]\neditor:", editor, "\nchange:", change)
    switch (change.origin) {
        case "+input":
            DocumentInsertAt(fd, change.text.join(), editor.getDoc().indexFromPos({
                line: change.from.line,
                ch: change.from.ch,
                sticky: null
            }))

            break
        case "+delete":
            DocumentDeleteAt(fd, editor.getDoc().indexFromPos({
                line: change.from.line,
                ch: change.from.ch,
                sticky: null
            }))
            break;
    }
}

function onCursorActivity(editor) {
    console.log("[on Cursor Activity]\neditor:", editor)
    DocumentChangeCursor(fd, editor.getDoc().indexFromPos(editor.getCursor()))
}
