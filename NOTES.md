# Extension language core

## Buffers

`current-buffer`

`other-buffer`

`set-buffer buffer` sets current buffer

`switch-to-buffer` swap current buffer, plus displays that buffer in current window

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

`get-buffer-create` find or creates a buffer by name

## status bar

`message message ...` printf to footer status line

`warning message ...` printf in red to footer status line

`interactive type message` prints message to footer status line and return user input
