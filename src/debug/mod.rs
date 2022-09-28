#[macro_export]
#[cfg(debug_assertions)]
macro_rules! debugln {
    () => {
        std::println!("\n")
    };
    ($($arg:tt)*) => {{
        std::println!(std::format_args_nl!($($arg)*));
    }};
}

#[macro_export]
#[allow_internal_unstable(print_internals)]
#[cfg(debug_assertions)]
macro_rules! debug {
    ($($arg:tt)*) => {{
        std::io::_print(std::format_args!($($arg)*));
    }};
}

#[macro_export]
#[cfg(not(debug_assertions))]
macro_rules! debugln {
    () => {
    };
    ($($arg:tt)*) => {{
    }};
}

#[macro_export]
#[allow_internal_unstable(print_internals)]
#[cfg(not(debug_assertions))]
macro_rules! debug {
    ($($arg:tt)*) => {{
    }};
}
