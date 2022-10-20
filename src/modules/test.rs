macro_rules! infix {
    // done converting
    (@cvt () $postfix:tt) => { infix!(@pfx () $postfix) };
    //                                |    |  ^ postfix expression
    //                                |    ^ operand stack
    //                                ^ postfix interpreter

    // infix to postfix conversion using the rules at the bottom of this page: http://csis.pace.edu/~wolf/CS122/infix-postfix.htm

    // at end of input, flush the operators to postfix
    (@cvt ($ophead:tt $($optail:tt)*) ($($postfix:tt)*)) => { infix!(@cvt ($($optail)*) ($($postfix)* $ophead)) };

    // 2. push an operator onto the stack if it's empty or has a left-paren on top
    (@cvt (                 ) $postfix:tt + $($tail:tt)*) => { infix!(@cvt (+               ) $postfix $($tail)*) };
    (@cvt (                 ) $postfix:tt - $($tail:tt)*) => { infix!(@cvt (-               ) $postfix $($tail)*) };
    (@cvt (                 ) $postfix:tt * $($tail:tt)*) => { infix!(@cvt (*               ) $postfix $($tail)*) };
    (@cvt (                 ) $postfix:tt / $($tail:tt)*) => { infix!(@cvt (/               ) $postfix $($tail)*) };
    (@cvt (LP $($optail:tt)*) $postfix:tt + $($tail:tt)*) => { infix!(@cvt (+ LP $($optail)*) $postfix $($tail)*) };
    (@cvt (LP $($optail:tt)*) $postfix:tt - $($tail:tt)*) => { infix!(@cvt (- LP $($optail)*) $postfix $($tail)*) };
    (@cvt (LP $($optail:tt)*) $postfix:tt * $($tail:tt)*) => { infix!(@cvt (* LP $($optail)*) $postfix $($tail)*) };
    (@cvt (LP $($optail:tt)*) $postfix:tt / $($tail:tt)*) => { infix!(@cvt (/ LP $($optail)*) $postfix $($tail)*) };

    // 3. push a left-paren onto the stack
    (@cvt ($($operator:tt)*) $postfix:tt ($($inner:tt)*) $($tail:tt)*) => { infix!(@cvt (LP $($operator)*) $postfix $($inner)* RP $($tail)*) };

    // 4. see right-paren, pop operators to postfix until left-paren
    (@cvt (LP         $($optail:tt)*) $postfix:tt       RP $($tail:tt)*) => { infix!(@cvt ($($optail)*) $postfix               $($tail)*   ) };
    (@cvt ($ophead:tt $($optail:tt)*) ($($postfix:tt)*) RP $($tail:tt)*) => { infix!(@cvt ($($optail)*) ($($postfix)* $ophead) RP $($tail)*) };

    // 5. if an operator w/ lower precedence is on top, just push
    (@cvt (+ $($optail:tt)*) $postfix:tt * $($tail:tt)*) => { infix!(@cvt (* + $($optail)*) $postfix $($tail)*) };
    (@cvt (- $($optail:tt)*) $postfix:tt * $($tail:tt)*) => { infix!(@cvt (* - $($optail)*) $postfix $($tail)*) };
    (@cvt (+ $($optail:tt)*) $postfix:tt / $($tail:tt)*) => { infix!(@cvt (/ + $($optail)*) $postfix $($tail)*) };
    (@cvt (- $($optail:tt)*) $postfix:tt / $($tail:tt)*) => { infix!(@cvt (/ - $($optail)*) $postfix $($tail)*) };

    // 6. if an operator w/ equal precedence is on top, pop and push
    (@cvt (+ $($optail:tt)*) ($($postfix:tt)*) + $($tail:tt)*) => { infix!(@cvt (+ $($optail)*) ($($postfix)* +) $($tail)*) };
    (@cvt (- $($optail:tt)*) ($($postfix:tt)*) - $($tail:tt)*) => { infix!(@cvt (- $($optail)*) ($($postfix)* -) $($tail)*) };
    (@cvt (+ $($optail:tt)*) ($($postfix:tt)*) - $($tail:tt)*) => { infix!(@cvt (- $($optail)*) ($($postfix)* +) $($tail)*) };
    (@cvt (- $($optail:tt)*) ($($postfix:tt)*) + $($tail:tt)*) => { infix!(@cvt (+ $($optail)*) ($($postfix)* -) $($tail)*) };
    (@cvt (* $($optail:tt)*) ($($postfix:tt)*) * $($tail:tt)*) => { infix!(@cvt (* $($optail)*) ($($postfix)* *) $($tail)*) };
    (@cvt (/ $($optail:tt)*) ($($postfix:tt)*) / $($tail:tt)*) => { infix!(@cvt (/ $($optail)*) ($($postfix)* /) $($tail)*) };
    (@cvt (* $($optail:tt)*) ($($postfix:tt)*) / $($tail:tt)*) => { infix!(@cvt (/ $($optail)*) ($($postfix)* *) $($tail)*) };
    (@cvt (/ $($optail:tt)*) ($($postfix:tt)*) * $($tail:tt)*) => { infix!(@cvt (* $($optail)*) ($($postfix)* /) $($tail)*) };

    // 7. if an operator w/ higher precedence is on top, pop it to postfix
    (@cvt (* $($optail:tt)*) ($($postfix:tt)*) + $($tail:tt)*) => { infix!(@cvt ($($optail)*) ($($postfix)* *) + $($tail)*) };
    (@cvt (* $($optail:tt)*) ($($postfix:tt)*) - $($tail:tt)*) => { infix!(@cvt ($($optail)*) ($($postfix)* *) - $($tail)*) };
    (@cvt (/ $($optail:tt)*) ($($postfix:tt)*) + $($tail:tt)*) => { infix!(@cvt ($($optail)*) ($($postfix)* /) + $($tail)*) };
    (@cvt (/ $($optail:tt)*) ($($postfix:tt)*) - $($tail:tt)*) => { infix!(@cvt ($($optail)*) ($($postfix)* /) - $($tail)*) };

    // 1. operands go to the postfix output
    (@cvt $operators:tt ($($postfix:tt)*) $head:tt $($tail:tt)*) => { infix!(@cvt $operators ($($postfix)* ($head)) $($tail)*) };

    // postfix interpreter
    (@pfx ($result:expr                     ) (                     )) => { $result };
    (@pfx (($a:expr) ($b:expr) $($stack:tt)*) (+        $($tail:tt)*)) => { infix!(@pfx ((($b + $a)) $($stack)*) ($($tail)*)) };
    (@pfx (($a:expr) ($b:expr) $($stack:tt)*) (-        $($tail:tt)*)) => { infix!(@pfx ((($b - $a)) $($stack)*) ($($tail)*)) };
    (@pfx (($a:expr) ($b:expr) $($stack:tt)*) (*        $($tail:tt)*)) => { infix!(@pfx ((($b * $a)) $($stack)*) ($($tail)*)) };
    (@pfx (($a:expr) ($b:expr) $($stack:tt)*) (/        $($tail:tt)*)) => { infix!(@pfx ((($b / $a)) $($stack)*) ($($tail)*)) };
    (@pfx ($($stack:tt)*                    ) ($head:tt $($tail:tt)*)) => { infix!(@pfx ($head       $($stack)*) ($($tail)*)) };

    ($($t:tt)*) => { infix!(@cvt () () $($t)*) }
    //                      |    |  |  ^ infix expression
    //                      |    |  ^ postfix expression
    //                      |    ^ operator stack
    //                      ^ convert infix to postfix
}

fn main() {
    println!("{}", infix!(1 + 2 * 3));
    println!("{}", infix!(1 * 2 + 3));
    println!("{}", infix!(((1 + 2) * 3) * 3));
    println!("{}", infix!(( 1 + 2  * 3) * 3));
    println!("{}", infix!(1 - 2 - 1));
}