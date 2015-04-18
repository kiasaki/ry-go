# Extension language core

## Buffers

`current-buffer`

`other-buffer`

? `set-buffer buffer` sets current buffer

? `switch-to-buffer` swap current buffer, plus displays that buffer in current window

`buffer-name buffer`

`buffer-filename buffer`

`buffer-read-only? buffer`

`buffer-size buffer` size of buffer in characters

`point buffer` index of cursor in buffer as int in characters

`point-min` minimum index your cursor can get to in buffer

`point-max` maximum index your cursor can get to in buffer

`save-excursion expr ..` saves current point and mark and resets them before returning

`goto-char i` moves point to specified position

`push-mark` save current point in mark buffer

`get-buffer name` finds a buffer by name

`get-buffer-create name filename` find or creates a buffer by name

## editor

`redraw-editor`

`editor-root-window`

`editor-first-window` Return to top left most leaf window of the editor

`editor-char-height` Height in char of editor

`editor-char-width` Width in char of editor

`set-editor-selected-window window`

`editor-selected-window`

`split-window window side`

`delete-window`

`delete-other-windows`

`delete-windows-on buffer`

## windows

```
RD = Right divider
LS = Left-schematic
LF = Left-frings
LF = Left-margin
        ____________________________________________
       |______________ Header Line ______________|RD| ^
     ^ |LS|LF|LM|                       |RM|RF|RS|  | |
     | |  |  |  |                       |  |  |  |  | |
Window |  |  |  |       Text Area       |  |  |  |  | Window
Body | |  |  |  |     (Window Body)     |  |  |  |  | Total
Height |  |  |  |                       |  |  |  |  | Height
     | |  |  |  |<- Window Body Width ->|  |  |  |  | |
     v |__|__|__|_______________________|__|__|__|  | |
       |_______________ Mode Line _______________|__| |
       |_____________ Bottom Divider _______________| v
        <---------- Window Total Width ------------>

active window = visible
selected window = window with cursor focus
```

`make-window`

`quit-window`

`window-parent` Returns parent window or nil

`window-top-child`

`window-left-child`

`window-leaf? window` returns non-nil if window is displaying a buffer, in other words, isn't a vertical or horizontal combination

`window-combined? window horizontal` Returns non-nil if and only is the window is part of a vertical combination, if _horizontal_ is non-nil it return non-nil if and only is the window is part of a horizontal combination

`window-child window` Return left most or top most child

`window-in-direction direction window wrap` Finds a window in the given direction. Direction must be `above`, `below`, `left` or `right` and wrap is a boolean indicating if window-in-direction return nil at edges or wraps

`window-full-height?` & `window-full-width?`

`window-total-height` & `window-total-width`

`window-body-height` & `window-body-width`

`window-buffer window`

`set-window-buffer window buffer`

`get-buffer-window buffer` Gets first active window displaying the given buffer

`get-buffer-window-list`

`window-point`

`set-window-point`

`recenter` centers point in the middle of scrollable region

`scroll-up down`

`scroll-down window`

`window-edges`

`window-inside-edges`

## status bar

`message message ...` printf to footer status line

`warning message ...` printf in red to footer status line

`interactive type message` prints message to footer status line and return user input

`current-message`

`clear-message`

`with-temp-message mess &body`

## hooks?

`after-save-hooks`

`pre-save-hooks`

`mini-buffer-setup-hook`

`mini-buffer-exit-hook`

`kill-buffer-hook`

`startup-hook`

## standard error

`wrong-type-error` Displays warning with message "Wrong type of argument"

`file-read-error`

`file-write-error`

`file-locked-error`
