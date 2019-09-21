#!/usr/bin/guile  --no-auto-compile
!#
(define (>= x y)
    (not (< x y))
)
(define (square x)
     (* x x)
)
(define (sum_max x y z)
(cond 
((and (>= x z ) (>= y z)) (+ (square x) (square y)))
((and (>= x y ) (>= z y)) (+ (square x) (square z)))
((and (>= y x ) (>= z x)) (+ (square z) (square y))) 
)
)
(display (sum_max 1 2 3))
(newline)
