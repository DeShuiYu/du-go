use std::time;
use walkdir::{DirEntry, WalkDir};
async fn calc_dir_disksize(root: &DirEntry) -> anyhow::Result<u64> {
    let total_file_size = WalkDir::new(root.path())
        .into_iter()
        .filter_map(|e| e.ok())
        .filter(|e| e.file_type().is_file())
        .map(|e| e.metadata().unwrap().len())
        .sum();
    Ok(total_file_size)
}

#[tokio::main]
async fn main() {
    let st = time::Instant::now();
    let root = "/Users/dsy/SSD/";
    let mut stream = WalkDir::new(root)
        .min_depth(1)
        .max_depth(1)
        .into_iter()
        .map(|entry| async move {
            if let Ok(entry) = entry {
                let totoal_fize_size = calc_dir_disksize(&entry).await.unwrap_or(0);
                println!(
                    "{},{}",
                    entry.path().display(),
                    humansize::format_size(totoal_fize_size, humansize::DECIMAL)
                );
            }
        });

    while let Some(res) = stream.next() {
        res.await;
    }
    println!("total time: {:?}", st.elapsed());
}
