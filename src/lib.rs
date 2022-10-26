mod utils;
use std::str;

use wasm_bindgen::prelude::*;

// When the `wee_alloc` feature is enabled, use `wee_alloc` as the global
// allocator.
#[cfg(feature = "wee_alloc")]
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

#[wasm_bindgen]
extern {
    fn alert(s: usize);
}

pub fn generate_parenthesis(n: i32) -> Vec<String> {
    if n < 1 {
        return vec![];
    }
    let mut result = Vec::new();
    dfs(n, 0, 0, &mut result, String::new());
    result
}

fn dfs(n: i32, left: i32, right: i32, result: &mut Vec<String>, mut path: String) {
    if left == n && right == n {
        result.push(path);
        return;
    }
    if left < n {
        let mut new_path = path.clone();
        // 向序列加入左括号
        new_path.push('(');
        dfs(n, left + 1, right, result, new_path);
    }
    if right < left {
        // 向序列加入剩余的右括号
        path.push(')');
        dfs(n, left, right + 1, result, path);
    }
}

#[wasm_bindgen]
pub fn greet(n: i32) -> usize {
    let s = generate_parenthesis(n);
    let slen = s.len();

    alert(slen);
    return slen;
}
