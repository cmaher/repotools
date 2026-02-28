use std::io;

fn main() {
    println!("hello");
}

fn helper() {
    // does stuff
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_main() {}
}
