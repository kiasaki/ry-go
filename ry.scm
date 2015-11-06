(use (srfi 1)
     ncurses)

(define editor-height 0)
(define editor-width 0)

(define minibuffer-text "")

(define (set-minibuffer-message message)
  (set! minibuffer-text message))

(define (log-err-about msg)
  (display (string-append "Error: " msg) (current-error-port))
  (newline (current-error-port))
  (exit 1))

(define (log-err msg)
  (display msg (current-error-port))
  (newline (current-error-port))
  msg)

(define (move-cursor-up lines pos)
  (if (> (cdr pos) 0)
    (cons (car pos) (- (cdr pos) 1))
    pos))

(define (move-cursor-down lines pos)
  (if (< (cdr pos) (- (length lines) 1))
    (cons (car pos) (+ (cdr pos) 1))
    pos))

(define (move-cursor-left lines pos)
  (if (> (car pos) 0)
    (cons (- (car pos) 1) (cdr pos))
    pos))

(define (move-cursor-right lines pos)
  (if (< (cdr pos) (length lines))
    (if (< (car pos) (- (string-length (list-ref lines (cdr pos))) 1))
      (cons (+ (car pos) 1) (cdr pos))
      pos)
    pos))

(define (ensure-valid-position lines pos)
  (let* ([x (car pos)]
         [y (cdr pos)]
         [height (- (length lines) 1)]
         [ny (if (< y height) y height)]
         [width (- (string-length (list-ref lines ny)) 1)]
         [nx (if (< x width) x width)])
    (set-minibuffer-message (string-append
                              " x: " (number->string x)
                              " y: " (number->string y)
                              " nx: " (number->string y)
                              " ny: " (number->string y)
                              " width: " (number->string width)
                              " height: " (number->string height)
                              ))
    (cons nx ny)))

(define (split-elt l elt)
  (let loop ((head '())
             (tail l)
             (i 0))
   (if (eq? tail '())
     (values l '())
     (if (= elt i)
       (values (reverse head) tail)
       (loop (cons (car tail) head)
             (cdr tail)
             (+ i 1))))))

(define (insert-string lines pos str)
  (call-with-values
    (lambda () (split-elt lines (cdr pos)))
    (lambda (head rest)
      (call-with-values
        (lambda () (split-elt (string->list (car rest)) (car pos)))
        (lambda (lhead lrest) (append head (cons (list->string (append lhead (string->list str) lrest)) (cdr rest))))))))

(define (insert-char lines pos new-char)
  (call-with-values
    (lambda () (split-elt lines (cdr pos)))
    (lambda (head rest)
      (call-with-values
        (lambda () (split-elt (string->list (car rest)) (car pos)))
        (lambda (lhead lrest) (append head (cons (list->string (append lhead (cons new-char lrest))) (cdr rest))))))))

(define (change-char lines pos new-char)
  (call-with-values
    (lambda () (split-elt lines (cdr pos)))
    (lambda (head rest)
      (call-with-values
        (lambda () (split-elt (string->list (car rest)) (car pos)))
        (lambda (lhead lrest) (append head (cons (list->string (append lhead (cons new-char (cdr lrest)))) (cdr rest))))))))

(define (delete-char lines pos)
  (if (and (< (cdr pos) (length lines)) (>= (cdr pos) 0))
    (if (and (< (car pos) (string-length (list-ref lines (cdr pos)))) (>= (car pos) 0))
      (call-with-values
        (lambda () (split-elt lines (cdr pos)))
        (lambda (head rest)
          (call-with-values
            (lambda () (split-elt (string->list (car rest)) (car pos)))
            (lambda (lhead lrest) (append head (cons (list->string (append lhead (cdr lrest))) (cdr rest)))))))
      lines)
    lines))

(define (delete-char-left lines pos running mode)
  (let ([nlines (delete-char lines (cons (- (car pos) 1) (cdr pos)))])
    (values nlines (move-cursor-left nlines pos) running mode)))

(define (delete-line lines line)
  (if (< line (length lines))
    (call-with-values
      (lambda () (split-elt lines line))
      (lambda (head rest) (append head (cdr rest))))
    lines))

(define (define-binding alist)
  (lambda (lines pos running mode)
    (let ([f (assv (getch) alist)])
     (if f
       ((cdr f) lines pos running mode)
       (values lines pos running mode)))))

(define (input-string end-marker)
  (move (- editor-height 1) 0)
  (let loop ([l (list)])
   (let ([c (getch)])
    (if (char=? c end-marker)
      (list->string (reverse l))
      (let* ([updated-list (cons c l)]
             [current-input (list->string (reverse updated-list))])
        (mvaddstr (- editor-height 1) 0 current-input)
        (loop updated-list))))))

(define normal-mode
  (define-binding
    (list
      (cons #\q (lambda (lines pos running mode) (values lines pos #f mode)))
      (cons #\i (lambda (lines pos running mode) (values lines pos running insert-mode)))
      (cons #\h (lambda (lines pos running mode) (values lines (move-cursor-left lines pos) running mode)))
      (cons #\j (lambda (lines pos running mode) (values lines (move-cursor-down lines pos) running mode)))
      (cons #\k (lambda (lines pos running mode) (values lines (move-cursor-up lines pos) running mode)))
      (cons #\l (lambda (lines pos running mode) (values lines (move-cursor-right lines pos) running mode)))
      (cons #\d
       (define-binding
         (list
           (cons #\d (lambda (lines pos running mode) (values (delete-line lines (cdr pos)) (ensure-valid-position lines pos) running mode)))
           (cons #\h delete-char-left)
           (cons #\j (lambda (lines pos running mode) (values (delete-char-left lines pos) (move-cursor-left lines pos) running mode)))
           (cons #\k (lambda (lines pos running mode) (values (delete-char lines pos) pos running mode)))
           (cons #\l (lambda (lines pos running mode) (values (delete-char lines pos) pos running mode))))))
      (cons #\: (lambda (lines pos running mode) (eval (input-string (integer->char 10)))))
      (cons #\x (lambda (lines pos running mode) (values (delete-char lines pos) pos running mode)))
      (cons #\r (lambda (lines pos running mode) (values (change-char lines pos (getch)) pos running mode))))))

(define (insert-mode lines pos running mode)
  (let ([c (getch)])
   (cond [(char=? c (integer->char 27)) (values lines pos running normal-mode)]
         [else (values (insert-char lines pos c) (move-cursor-right lines pos) running mode)])))

(define (display-lines lines) ;impure
  (wclear (stdscr))
  (let loop ([l lines]
             [y 0])
    (if (not (null? l))
      (begin
        (mvaddstr y 0 (car l))
        (loop (cdr l) (+ y 1)))))
  (wrefresh (stdscr)))

(define (display-status-bar lines pos)
  (let ([pos-text (string-append (number->string (car pos)) "," (number->string (cdr pos)))])
    (attron A_REVERSE)
    (mvaddstr (- editor-height 2) 0 (make-string editor-width #\-))
    (mvaddstr (- editor-height 2) (- editor-width (string-length pos-text) 1) pos-text)
    (attroff A_REVERSE)))

(define (display-minibuffer text)
  (mvaddstr (- editor-height 1) 0 text))

(define (main) ;impure
  (initscr)
  (cbreak)
  (noecho)
  (keypad (stdscr) #t)
  (let loop ([lines (list "hello" "world" "bluh")]
             [pos (cons 0 0)]
             [running #t]
             [mode normal-mode])
    (if running
      (begin
        (let-values ([[my mx] (getmaxyx (stdscr))])
          (set! editor-width mx)
          (set! editor-height my))
        (display-lines lines)
        (display-status-bar lines pos)
        (display-minibuffer minibuffer-text)
        (move (cdr pos) (car pos))
        (call-with-values (lambda () (mode lines pos running mode)) loop))))
  (endwin))

(main)
