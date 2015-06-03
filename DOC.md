# API

## Buffers

### Native (Go)

`make-buffer name filename` Creates a new buffer instance

## Editor

### Native (Go)

`start-editor` Start the editor and inits the frontend

`stop-editor` Stops the editor and closes the frontend

`editor-keypresses-chan` Return the channel with incoming keypresses

`clear-editor` Fills the screen with the default background color

`set-cell` Takes x, y, character, fg, bg

### Runtime (syp)

`quit` Stops the editor and exits
