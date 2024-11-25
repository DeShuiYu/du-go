use std::env::args;
use std::{fs, time};
use walkdir::WalkDir;
use rayon::prelude::*;
use humansize;

fn get_dir_size(path: &str) -> Result<u64, std::io::Error> {
    let total = WalkDir::new(path)
        .into_iter()
        .par_bridge()  // 转换为并行迭代器
        .filter_map(|e| e.ok())
        .filter(|e| e.path().is_file())
        .map(|e| fs::metadata(e.path()).map(|m| m.len()).unwrap_or(0))
        .sum();
    Ok(total)
}


fn main() {
    let st = time::Instant::now();
    args().skip(1).for_each(|arg| {
        WalkDir::new(arg).min_depth(1).max_depth(1).into_iter().par_bridge().for_each(|entry| {
            let entry = entry.unwrap();
            if let Ok(size) = get_dir_size(entry.path().to_str().unwrap()) {
                println!("{},{}", entry.path().display(), humansize::format_size(size,humansize::BINARY))
            }
        })
    });
    println!("Total time: {:.2?}", st.elapsed());
}

