use std::collections::HashSet;
use std::env;
use std::process::Command;

fn main() {
    let args: Vec<String> = env::args().collect();
    let doc_id = if args.len() > 1 {
        &args[1]
    } else {
        "1ZXB6pnQEvym9pkDDoSxl49mxiiMAkzx0tS50OuwRt3Q"
    };

    let url = format!("https://drive.google.com/file/d/{}/view", doc_id);
    println!("[*] Target Document: {}", url);

    let output = Command::new("curl")
        .args(&["-s", "-L", "-A", "Mozilla/5.0", &url])
        .output()
        .expect("Failed to execute curl");

    let html = String::from_utf8_lossy(&output.stdout);

    if let Some(start) = html.find("<title>") {
        if let Some(end) = html[start + 7..].find(" - Google Drive</title>") {
            println!("[+] Title: {}", &html[start + 7..start + 7 + end]);
        }
    }

    let mut emails = HashSet::new();
    for word in html.split(|c: char| c == '"' || c == ',' || c == ';' || c.is_whitespace()) {
        if word.ends_with("@gmail.com") && word.contains('@') {
            let clean = word.trim_matches(|c: char| !c.is_alphanumeric() && c != '@' && c != '.' && c != '_' && c != '-');
            if clean.contains('@') {
                emails.insert(clean.to_string());
            }
        }
    }

    for email in &emails {
        let username = email.split('@').next().unwrap_or(email);
        println!("[+] Recovered Owner Email: {}", email);
        println!("[+] Extracted Username:   {}", username);

        let yt_url = format!("https://www.youtube.com/@{}", username);
        println!("[*] Checking linked profile: {}", yt_url);

        let yt_out = Command::new("curl")
            .args(&["-s", "-L", "-A", "Mozilla/5.0", &yt_url])
            .output();

        if let Ok(out) = yt_out {
            let yt_html = String::from_utf8_lossy(&out.stdout);
            let marker = "\"mainEntity\":{\"type\":\"Person\",\"name\":\"";
            if let Some(pos) = yt_html.find(marker) {
                let rest = &yt_html[pos + marker.len()..];
                if let Some(end) = rest.find('"') {
                    println!("[!] Discovered Linked Real Name: {}", &rest[..end]);
                }
            }

            let ch_marker = "\"channelId\":\"";
            if let Some(pos) = yt_html.find(ch_marker) {
                let rest = &yt_html[pos + ch_marker.len()..];
                if let Some(end) = rest.find('"') {
                    println!("[!] Linked Channel / Gaia ID:   {}", &rest[..end]);
                }
            }
        }
    }
}
