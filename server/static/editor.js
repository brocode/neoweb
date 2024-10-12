function initKeyListener() {
    // Listen for global keypresses on the entire document
    document.addEventListener("keydown", function (event) {


        // F keys should stay in the browser for now
        if (event.key.startsWith('F') && event.key.length > 1 && !isNaN(event.key.slice(1))) {
            return;
        }

        // Modifier keys
        const modifierKeys = ['Shift', 'Control', 'Alt', 'Meta'];
        if (modifierKeys.includes(event.key)) {
            return;
        }

        // Prevent copy and paste (Ctrl+C, Ctrl+V, Cmd+C, Cmd+V)
        if ((event.ctrlKey || event.metaKey) && (event.key === 'c' || event.key === 'v')) {
            return;
        }
        event.preventDefault()

        // Send the keypress information to the server
        fetch('/keypress', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                key: event.key,
                shiftKey: event.shiftKey,
                ctrlKey: event.ctrlKey,
                altKey: event.altKey,
                metaKey: event.metaKey
            })
        })
            .then(response => response.text())
            .then(data => console.debug('Keypress response:', data))
            .catch((error) => console.error('Error:', error));
    });
}

function initEventSource() {
    const source = new EventSource('/events');

    source.addEventListener('render', function (event) {
        console.debug(event.data);

        document.getElementById('editor').innerHTML = event.data;
    });

    source.onerror = function (error) {
        console.error('Error receiving SSE:', error);
    };
}


function initPaste() {
    document.addEventListener('paste', (event) => {
        // Prevent the default paste behavior
        event.preventDefault();

        // Get the pasted content as plain text
        const pastedData = (event.clipboardData || window.clipboardData).getData('text');

        // Send the plain text to your endpoint
        sendPastedText(pastedData);
    });
}

function sendPastedText(text) {
    fetch("/paste", {
        method: "POST",
        headers: {
            "Content-Type": "text/plain"
        },
        body: text
    })
        .then(response => response.text())
        .then(data => console.debug('Paste response:', data))
        .catch((error) => console.error('Error:', error));
}
